package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ClusterScanSpec defines the desired state of ClusterScan
type ClusterScanSpec struct {
	Schedule string `json:"schedule,omitempty"` // Cron schedule for recurring scans
	OneOff   bool   `json:"oneOff,omitempty"`   // Flag for one-off scans
}

// ClusterScanStatus defines the observed state of ClusterScan
type ClusterScanStatus struct {
	LastRunTime *metav1.Time `json:"lastRunTime,omitempty"`
	JobName     string       `json:"jobName,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// ClusterScan is the Schema for the clusterscans API
type ClusterScan struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterScanSpec   `json:"spec,omitempty"`
	Status ClusterScanStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ClusterScanList contains a list of ClusterScan
type ClusterScanList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterScan `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ClusterScan{}, &ClusterScanList{})
}
