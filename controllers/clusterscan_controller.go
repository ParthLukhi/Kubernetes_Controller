package controllers

import (
	"context"
	"fmt"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	scanv1 "github.com/ParthLukhi/cluster-scan-controller/api/v1"
)

// ClusterScanReconciler reconciles a ClusterScan object
type ClusterScanReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=scans.mydomain.com,resources=clusterscans,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=scans.mydomain.com,resources=clusterscans/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=batch,resources=jobs;cronjobs,verbs=get;list;watch;create;update;patch;delete

func (r *ClusterScanReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Fetch the ClusterScan instance
	scan := &scanv1.ClusterScan{}
	err := r.Get(ctx, req.NamespacedName, scan)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Check if the scan is one-off or recurring
	if scan.Spec.OneOff {
		// Handle one-off job
		job := &batchv1.Job{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("%s-job", scan.Name),
				Namespace: req.Namespace,
			},
			Spec: batchv1.JobSpec{
				Template: corev1.PodTemplateSpec{
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  "scan",
								Image: "busybox",
								Command: []string{
									"sh",
									"-c",
									"echo Scan completed",
								},
							},
						},
						RestartPolicy: corev1.RestartPolicyOnFailure,
					},
				},
			},
		}

		// Set the owner reference
		if err := controllerutil.SetControllerReference(scan, job, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}

		// Create or update the Job
		foundJob := &batchv1.Job{}
		err = r.Get(ctx, types.NamespacedName{Name: job.Name, Namespace: job.Namespace}, foundJob)
		if err != nil && client.IgnoreNotFound(err) != nil {
			return ctrl.Result{}, err
		}

		if err := r.Create(ctx, job); err != nil {
			return ctrl.Result{}, err
		}

		scan.Status.JobName = job.Name
		scan.Status.LastRunTime = &metav1.Time{Time: time.Now()}
		if err := r.Status().Update(ctx, scan); err != nil {
			return ctrl.Result{}, err
		}

		return ctrl.Result{}, nil
	}

	// Handle recurring CronJob
	cronJob := &batchv1beta1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-cronjob", scan.Name),
			Namespace: req.Namespace,
		},
		Spec: batchv1beta1.CronJobSpec{
			Schedule: scan.Spec.Schedule,
			JobTemplate: batchv1beta1.JobTemplateSpec{
				Spec: batchv1.JobSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "scan",
									Image: "busybox",
									Command: []string{
										"sh",
										"-c",
										"echo Recurring scan completed",
									},
								},
							},
							RestartPolicy: corev1.RestartPolicyOnFailure,
						},
					},
				},
			},
		},
	}

	// Set the owner reference
	if err := controllerutil.SetControllerReference(scan, cronJob, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	// Create or update the CronJob
	foundCronJob := &batchv1beta1.CronJob{}
	err = r.Get(ctx, types.NamespacedName{Name: cronJob.Name, Namespace: cronJob.Namespace}, foundCronJob)
	if err != nil && client.IgnoreNotFound(err) != nil {
		return ctrl.Result{}, err
	}

	if err := r.Create(ctx, cronJob); err != nil {
		return ctrl.Result{}, err
	}

	scan.Status.JobName = cronJob.Name
	if err := r.Status().Update(ctx, scan); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *ClusterScanReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&scanv1.ClusterScan{}).
		Owns(&batchv1.Job{}).
		Owns(&batchv1beta1.CronJob{}).
		Complete(r)
}
