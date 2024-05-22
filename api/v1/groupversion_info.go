package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// GroupVersion is group version used to register these objects
var GroupVersion = schema.GroupVersion{Group: "scans.mydomain.com", Version: "v1"}

// SchemeGroupVersion is group version used to register these objects
var SchemeGroupVersion = GroupVersion

var (
	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilder = &SchemeBuilder{GroupVersion}
	// AddToScheme is used to add go types to the GroupVersionKind scheme
	AddToScheme = SchemeBuilder.AddToScheme
)

// Resource takes an unqualified resource and returns a Group qualified GroupResource
func Resource(resource string) schema.GroupResource {
	return GroupVersion.WithResource(resource).GroupResource()
}
