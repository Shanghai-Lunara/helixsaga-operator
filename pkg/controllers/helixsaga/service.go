package helixsaga

import (
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	helixSagaV1 "github.com/Shanghai-Lunara/helixsaga-operator/pkg/apis/helixsaga/v1"
	k8sCoreV1 "github.com/nevercase/k8s-controller-custom-resource/core/v1"
)

func NewService(hs *helixSagaV1.HelixSaga, spec helixSagaV1.HelixSagaAppSpec) *coreV1.Service {
	labels := map[string]string{
		k8sCoreV1.LabelApp:        OperatorKindName,
		k8sCoreV1.LabelController: hs.Name,
		k8sCoreV1.LabelName:       spec.Name,
	}
	return &coreV1.Service{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      k8sCoreV1.GetServiceName(spec.Name),
			Namespace: hs.Namespace,
			OwnerReferences: []metaV1.OwnerReference{
				*metaV1.NewControllerRef(hs, helixSagaV1.SchemeGroupVersion.WithKind(OperatorKindName)),
			},
			Labels: labels,
		},
		Spec: coreV1.ServiceSpec{
			Type:     k8sCoreV1.GetServiceType(spec.ServiceType),
			Ports:    spec.ServicePorts,
			Selector: labels,
		},
	}
}
