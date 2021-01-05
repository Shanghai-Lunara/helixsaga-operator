package helixsaga

const controllerAgentName = "helix-saga-controller"
const OperatorKindName = "HelixSaga"

const (
	// SuccessSynced is used as part of the Event 'reason' when a Foo is synced
	SuccessSynced = "Synced"
	// ErrResourceExists is used as part of the Event 'reason' when a Foo fails
	// to sync due to a Deployment of the same name already existing.
	ErrResourceExists = "ErrResourceExists"

	ErrResourceNotMatch = "ErrResourceNotMatch err:%s"

	// MessageResourceExists is the message used for Events when a resource
	// fails to sync due to a Deployment already existing
	MessageResourceExists = "Resource %q already exists and is not managed by Foo"
	// MessageResourceSynced is the message used for an Event fired when a Foo
	// is synced successfully
	MessageResourceSynced = "Foo synced successfully"
)

const (
	NodeName          = "NODE_NAME"
	HostIP            = "HOST_IP"
	PodName           = "POD_NAME"
	PodNamespace      = "POD_NAMESPACE"
	PodIP             = "POD_IP"
	PodServiceAccount = "POD_SERVICE_ACCOUNT"
)
