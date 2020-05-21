package helixsaga

import (
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/klog"

	helixSagaV1 "github.com/Shanghai-Lunara/helixsaga-operator/pkg/apis/helixsaga/v1"
	helixSagaClientSet "github.com/Shanghai-Lunara/helixsaga-operator/pkg/generated/helixsaga/clientset/versioned"
	k8sCoreV1 "github.com/nevercase/k8s-controller-custom-resource/core/v1"
)

func NewNginxPhpFpm(ks k8sCoreV1.KubernetesResource, client helixSagaClientSet.Interface, hs *helixSagaV1.HelixSaga, spec helixSagaV1.HelixSagaCoreSpec) error {
	ss, err := ks.StatefulSet().Get(hs.Namespace, spec.Name)
	if err != nil {
		klog.Info("statefulSet err:", err)
		if !errors.IsNotFound(err) {
			return err
		}
		klog.Info("new statefulSet")
		if ss, err = ks.StatefulSet().Create(hs.Namespace, spec.Name, NewStatefulSet(hs, spec)); err != nil {
			return err
		}
	}
	klog.Info("rds:", *spec.Replicas)
	klog.Info("statefulSet:", *ss.Spec.Replicas)
	if spec.Replicas != nil && *spec.Replicas != *ss.Spec.Replicas {
		if ss, err = ks.StatefulSet().Update(hs.Namespace, NewStatefulSet(hs, spec)); err != nil {
			klog.Info(err)
			return err
		}
	}
	if err = updateFooStatus(hs, client, ss); err != nil {
		return err
	}
	return nil
}
