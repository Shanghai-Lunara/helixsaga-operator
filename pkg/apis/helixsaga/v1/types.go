package v1

import (
	coreV1 "k8s.io/api/core/v1"
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

	Status HelixSagaStatus `json:"status"`
}

//HelixSagaSpec is the spec for a HelixSaga resource
type HelixSagaSpec struct {
	ConfigMap           HelixSagaConfigMap  `json:"config_map"`
	VersionSpec         HelixSagaCoreSpec   `json:"version_spec"`
	ApiSpec             HelixSagaCoreSpec   `json:"api_spec"`
	GameSpec            HelixSagaCoreSpec   `json:"game_spec"`
	PayNotifySpec       HelixSagaCoreSpec   `json:"pay_notify_spec"`
	GmtSpec             HelixSagaCoreSpec   `json:"gmt_spec"`
	FriendSpec          HelixSagaCoreSpec   `json:"friend_spec"`
	QueueSpec           HelixSagaCoreSpec   `json:"queue_spec"`
	RankSpec            HelixSagaCoreSpec   `json:"rank_spec"`
	ChatSpec            PhpWorkermanSpec    `json:"chat_spec"`
	HeartSpec           PhpWorkermanSpec    `json:"heart_spec"`
	CampaignSpec        CampaignSpec        `json:"campaign_spec"`
	GuildWarSpec        GuildWarSpec        `json:"guild_war_spec"`
	AppNotificationSpec AppNotificationSpec `json:"app_notification_spec"`
}

type HelixSagaConfigMap struct {
	Volume      coreV1.Volume      `json:"volume"`
	VolumeMount coreV1.VolumeMount `json:"volumeMount"`
}

//HelixSagaCoreSpec is the sub spec for a HelixSaga resource
type HelixSagaCoreSpec struct {
	// Name of the container specified as a DNS_LABEL.
	// Each container in a pod must have a unique name (DNS_LABEL).
	// Cannot be updated.
	Name string `json:"name" protobuf:"bytes,1,rep,name=name"`
	// Replicas is the number of desired replicas.
	// This is a pointer to distinguish between explicit zero and unspecified.
	// Defaults to 1.
	// More info: https://kubernetes.io/docs/concepts/workloads/controllers/replicationcontroller#what-is-a-replicationcontroller
	// +optional
	Replicas *int32 `json:"replicas,omitempty" protobuf:"bytes,2,rep,name=replicas"`
	// Docker image name.
	// More info: https://kubernetes.io/docs/concepts/containers/images
	// This field is optional to allow higher level config management to default or override
	// container images in workload controllers like Deployments and StatefulSets.
	// +optional
	Image string `json:"image,omitempty" protobuf:"bytes,3,opt,name=image"`
	// ImagePullSecrets is an optional list of references to secrets in the same namespace to use for pulling any of the images used by this PodSpec.
	// If specified, these secrets will be passed to individual puller implementations for them to use. For example,
	// in the case of docker, only DockerConfig type secrets are honored.
	// More info: https://kubernetes.io/docs/concepts/containers/images#specifying-imagepullsecrets-on-a-pod
	// +optional
	// +patchMergeKey=name
	// +patchStrategy=merge
	ImagePullSecrets []coreV1.LocalObjectReference `json:"imagePullSecrets,omitempty" patchStrategy:"merge" patchMergeKey:"name" protobuf:"bytes,4,rep,name=imagePullSecrets"`
	// List of environment variables to set in the container.
	// Cannot be updated.
	// +optional
	// +patchMergeKey=name
	// +patchStrategy=merge
	Env []coreV1.EnvVar `json:"env,omitempty" patchStrategy:"merge" patchMergeKey:"name" protobuf:"bytes,5,rep,name=env"`
	// Pod volumes to mount into the container's filesystem.
	// Cannot be updated.
	// +optional
	// +patchMergeKey=mountPath
	// +patchStrategy=merge
	VolumeMounts []coreV1.VolumeMount `json:"volumeMounts,omitempty" patchStrategy:"merge" patchMergeKey:"mountPath" protobuf:"bytes,6,rep,name=volumeMounts"`
}

//PhpWorkermanSpec is the sub spec for a HelixSaga resource
type PhpWorkermanSpec struct {
	RegisterSpec       HelixSagaCoreSpec `json:"register_spec"`
	GatewaySpec        HelixSagaCoreSpec `json:"gateway_spec"`
	BusinessWorkerSpec HelixSagaCoreSpec `json:"business_worker_spec"`
}

//Campaign is the sub spec for a HelixSaga resource
type CampaignSpec struct {
	GatewaySpec HelixSagaCoreSpec `json:"gateway_spec"`
}

//GuildWarSpec is the sub spec for a HelixSaga resource
type GuildWarSpec struct {
	RegisterSpec HelixSagaCoreSpec `json:"register_spec"`
	GatewaySpec  HelixSagaCoreSpec `json:"gateway_spec"`
}

//AppNotificationSpec is the sub spec for a HelixSaga resource
type AppNotificationSpec struct {
	DispatchSpec HelixSagaCoreSpec `json:"dispatch_spec"`
	LogicSpec    HelixSagaCoreSpec `json:"logic_spec"`
}

// HelixSagaStatus is the status for a HelixSaga resource
type HelixSagaStatus struct {
	VersionStatus         HelixSagaCoreStatus   `json:"version_status"`
	ApiStatus             HelixSagaCoreStatus   `json:"api_status"`
	GameStatus            HelixSagaCoreStatus   `json:"game_status"`
	PayNotifyStatus       HelixSagaCoreStatus   `json:"pay_notify_status"`
	GmtStatus             HelixSagaCoreStatus   `json:"gmt_status"`
	FriendStatus          HelixSagaCoreStatus   `json:"friend_status"`
	QueueStatus           HelixSagaCoreStatus   `json:"queue_status"`
	RankStatus            HelixSagaCoreStatus   `json:"rank_status"`
	ChatStatus            PhpWorkermanStatus    `json:"chat_status"`
	HeartStatus           PhpWorkermanStatus    `json:"heart_status"`
	CampaignStatus        CampaignStatus        `json:"campaign_status"`
	GuildWarStatus        GuildWarStatus        `json:"guild_war_status"`
	AppNotificationStatus AppNotificationStatus `json:"app_notification_status"`
}

//HelixSagaCoreStatus is the sub status for a HelixSaga resource
type HelixSagaCoreStatus struct {
	// The generation observed by the deployment controller.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty" protobuf:"varint,1,opt,name=observedGeneration"`

	// Total number of non-terminated pods targeted by this deployment (their labels match the selector).
	// +optional
	Replicas int32 `json:"replicas,omitempty" protobuf:"varint,2,opt,name=replicas"`

	// Total number of non-terminated pods targeted by this deployment that have the desired template spec.
	// +optional
	UpdatedReplicas int32 `json:"updatedReplicas,omitempty" protobuf:"varint,3,opt,name=updatedReplicas"`

	// Total number of ready pods targeted by this deployment.
	// +optional
	ReadyReplicas int32 `json:"readyReplicas,omitempty" protobuf:"varint,4,opt,name=readyReplicas"`

	// Total number of available pods (ready for at least minReadySeconds) targeted by this deployment.
	// +optional
	AvailableReplicas int32 `json:"availableReplicas,omitempty" protobuf:"varint,5,opt,name=availableReplicas"`

	// Total number of unavailable pods targeted by this deployment. This is the total number of
	// pods that are still required for the deployment to have 100% available capacity. They may
	// either be pods that are running but not yet available or pods that still have not been created.
	// +optional
	UnavailableReplicas int32 `json:"unavailableReplicas,omitempty" protobuf:"varint,6,opt,name=unavailableReplicas"`
}

//PhpWorkermanStatus is the sub Status for a HelixSaga resource
type PhpWorkermanStatus struct {
	RegisterStatus       HelixSagaCoreStatus `json:"register_status"`
	GatewayStatus        HelixSagaCoreStatus `json:"gateway_status"`
	BusinessWorkerStatus HelixSagaCoreStatus `json:"business_worker_status"`
}

//Campaign is the sub Status for a HelixSaga resource
type CampaignStatus struct {
	GatewayStatus HelixSagaCoreStatus `json:"gateway_status"`
}

//GuildWarStatus is the sub Status for a HelixSaga resource
type GuildWarStatus struct {
	RegisterStatus HelixSagaCoreStatus `json:"register_status"`
	GatewayStatus  HelixSagaCoreStatus `json:"gateway_status"`
}

//AppNotificationStatus is the sub Status for a HelixSaga resource
type AppNotificationStatus struct {
	DispatchStatus HelixSagaCoreStatus `json:"dispatch_status"`
	LogicStatus    HelixSagaCoreStatus `json:"logic_status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

//HelixSagaList is a list of HelixSaga resources
type HelixSagaList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []HelixSaga `json:"items"`
}
