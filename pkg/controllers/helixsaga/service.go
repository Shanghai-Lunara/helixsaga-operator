package helixsaga

import (
	helixSagav1 "github.com/Shanghai-Lunara/helixsaga-operator/pkg/apis/helixsaga/v1"
	"github.com/Shanghai-Lunara/helixsaga-operator/pkg/serviceloadbalancer"
	k8scorev1 "github.com/nevercase/k8s-controller-custom-resource/core/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewService(hs *helixSagav1.HelixSaga, spec *helixSagav1.HelixSagaAppSpec) *corev1.Service {
	labels := map[string]string{
		k8scorev1.LabelApp:        OperatorKindName,
		k8scorev1.LabelController: hs.Name,
		k8scorev1.LabelName:       spec.Name,
	}
	annotations := make(map[string]string, 0)
	switch spec.ServiceType {
	case corev1.ServiceTypeLoadBalancer:
		annotations = serviceloadbalancer.Get().Annotations
		if spec.ServiceWhiteList == true {
			for k, v := range serviceloadbalancer.Get().WhiteList {
				annotations[k] = v
			}
		}
	}
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      k8scorev1.GetServiceName(spec.Name),
			Namespace: hs.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(hs, helixSagav1.SchemeGroupVersion.WithKind(OperatorKindName)),
			},
			Labels:      labels,
			Annotations: annotations,
		},
		Spec: corev1.ServiceSpec{
			Type:     k8scorev1.GetServiceType(spec.ServiceType),
			Ports:    spec.ServicePorts,
			Selector: labels,
		},
	}
}
