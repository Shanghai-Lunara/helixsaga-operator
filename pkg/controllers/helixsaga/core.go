package helixsaga

import (
	appsV1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/klog"

	helixSagaV1 "github.com/Shanghai-Lunara/helixsaga-operator/pkg/apis/helixsaga/v1"
	helixSagaClientSet "github.com/Shanghai-Lunara/helixsaga-operator/pkg/generated/helixsaga/clientset/versioned"
	k8sCoreV1 "github.com/nevercase/k8s-controller-custom-resource/core/v1"
)

func NewStatefulSetAndService(ks k8sCoreV1.KubernetesResource, client helixSagaClientSet.Interface, hs *helixSagaV1.HelixSaga, spec helixSagaV1.HelixSagaCoreSpec) error {
	ss, err := ks.StatefulSet().Get(hs.Namespace, spec.Name)
	if err != nil {
		klog.Info("statefulSet err:", err)
		if !errors.IsNotFound(err) {
			return err
		}
		klog.Info("new statefulSet")
		if ss, err = ks.StatefulSet().Create(hs.Namespace, NewStatefulSet(hs, spec)); err != nil {
			return err
		}
		if _, err = ks.Service().Create(hs.Namespace, NewService(hs, spec)); err != nil {
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
		if _, err = ks.Service().Update(hs.Namespace, NewService(hs, spec)); err != nil {
			return err
		}
	}
	if err = updateStatus(hs, client, ss, spec.Name); err != nil {
		return err
	}
	return nil
}

func updateStatus(foo *helixSagaV1.HelixSaga, clientSet helixSagaClientSet.Interface, ss *appsV1.StatefulSet, name string) error {
	// NEVER modify objects from the store. It's a read-only, local cache.
	// You can use DeepCopy() to make a deep copy of original object and modify this copy
	// Or create a copy manually for better performance
	fooCopy := foo.DeepCopy()
	t := make([]helixSagaV1.HelixSagaCore, 0)
	for _, v := range fooCopy.Spec.Applications {
		if v.Spec.Name == name {
			v.Status.Replicas = ss.Status.Replicas
			v.Status.AvailableReplicas = ss.Status.Replicas
		}
		t = append(t, v)
	}
	fooCopy.Spec.Applications = t

	// If the CustomResourceSubResources feature gate is not enabled,
	// we must use Update instead of UpdateStatus to update the Status block of the RedisOperator resource.
	// UpdateStatus will not allow changes to the Spec of the resource,
	// which is ideal for ensuring nothing other than resource status has been updated.
	_, err := clientSet.HelixsagaV1().HelixSagas(foo.Namespace).Update(fooCopy)
	return err
}
