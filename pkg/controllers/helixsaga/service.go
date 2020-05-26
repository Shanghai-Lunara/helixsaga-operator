package helixsaga

import (
	"fmt"

	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	helixSagaV1 "github.com/Shanghai-Lunara/helixsaga-operator/pkg/apis/helixsaga/v1"
	k8sCoreV1 "github.com/nevercase/k8s-controller-custom-resource/core/v1"
)

func NewService(hs *helixSagaV1.HelixSaga, spec helixSagaV1.HelixSagaCoreSpec) *coreV1.Service {
	labels := map[string]string{
		"app":        operatorKindName,
		"controller": hs.Name,
		"role":       spec.Name,
	}
	return &coreV1.Service{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      fmt.Sprintf(k8sCoreV1.ServiceNameTemplate, spec.Name),
			Namespace: hs.Namespace,
			OwnerReferences: []metaV1.OwnerReference{
				*metaV1.NewControllerRef(hs, helixSagaV1.SchemeGroupVersion.WithKind(operatorKindName)),
			},
			Labels: labels,
		},
		Spec: coreV1.ServiceSpec{
			Ports:    spec.ServicePorts,
			Selector: labels,
		},
	}
}
