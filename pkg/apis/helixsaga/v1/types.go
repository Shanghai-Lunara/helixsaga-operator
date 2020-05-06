package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

//HelixSaga describes a HelixSaga resource
type HelixSaga struct {
	// TypeMeta is the metadata for the resource, like kind and apiversion
	metav1.TypeMeta `json:",inline"`
	// ObjectMeta contains the metadata for the particular object, including
	// things like...
	//  - name
	//  - namespace
	//  - self link
	//  - labels
	//  - ... etc ...
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec is the custom resource spec
	Spec HelixSagaSpec `json:"spec"`

	Status HelixSagaSpecStatus `json:"status"`
}

//HelixSagaSpec is the spec for a HelixSaga resource
type HelixSagaSpec struct {
	MasterSpec HelixSagaDeploymentSpec `json:"master_spec"`
	SlaveSpec  HelixSagaDeploymentSpec `json:"slave_spec"`
}

//HelixSagaDeploymentSpec is the sub spec for a HelixSaga resource
type HelixSagaDeploymentSpec struct {
	DeploymentName   string `json:"deploymentName"`
	Replicas         *int32 `json:"replicas"`
	Image            string `json:"image"`
	ImagePullSecrets string `json:"imagePullSecrets"`
}

// HelixSagaStatus is the status for a HelixSaga resource
type HelixSagaSpecStatus struct {
	MasterStatus HelixSagaDeploymentStatus `json:"master_status"`
	SlaveStatus  HelixSagaDeploymentStatus `json:"slave_status"`
}

//HelixSagaDeploymentStatus is the sub status for a HelixSaga resource
type HelixSagaDeploymentStatus struct {
	AvailableReplicas int32 `json:"availableReplicas"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

//HelixSagaList is a list of HelixSaga resources
type HelixSagaList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []HelixSaga `json:"items"`
}
