package helixsaga

import (
	appsV1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/klog"

	helixSagaV1 "github.com/Shanghai-Lunara/helixsaga-operator/pkg/apis/helixsaga/v1"
	helixSagaClientSet "github.com/Shanghai-Lunara/helixsaga-operator/pkg/generated/helixsaga/clientset/versioned"
	k8sCoreV1 "github.com/nevercase/k8s-controller-custom-resource/core/v1"
)

func NewStatefulSetAndService(ks k8sCoreV1.KubernetesResource, client helixSagaClientSet.Interface, hs *helixSagaV1.HelixSaga, spec helixSagaV1.HelixSagaAppSpec) error {
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
		if len(spec.ServicePorts) > 0 {
			klog.Info("new service init")
			if _, err = ks.Service().Create(hs.Namespace, NewService(hs, spec)); err != nil {
				return err
			}
		}
	}
	klog.Info("rds:", *spec.Replicas)
	klog.Info("statefulSet:", *ss.Spec.Replicas)
	if spec.Replicas != nil && *spec.Replicas != *ss.Spec.Replicas || spec.Image != ss.Spec.Template.Spec.Containers[0].Image {
		if ss, err = ks.StatefulSet().Update(hs.Namespace, NewStatefulSet(hs, spec)); err != nil {
			klog.V(2).Info(err)
			return err
		}
		if len(spec.ServicePorts) > 0 {
			if _, err = ks.Service().Update(hs.Namespace, NewService(hs, spec)); err != nil {
				return err
			}
		}
	}
	if len(spec.ServicePorts) == 0 {
		if err = ks.Service().Delete(hs.Namespace, k8sCoreV1.GetServiceName(spec.Name)); err != nil {
			return err
		}
	} else {
		_, err = ks.Service().Get(hs.Namespace, k8sCoreV1.GetServiceName(spec.Name))
		if err != nil {
			klog.Info("service err:", err)
			if !errors.IsNotFound(err) {
				return err
			}
			klog.Info("new service check")
			if _, err = ks.Service().Create(hs.Namespace, NewService(hs, spec)); err != nil {
				return err
			}
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
	t := make([]helixSagaV1.HelixSagaApp, 0)
	for _, v := range fooCopy.Spec.Applications {
		if v.Spec.Name == name {
			v.Status.ObservedGeneration = ss.Status.ObservedGeneration
			v.Status.Replicas = ss.Status.Replicas
			v.Status.ReadyReplicas = ss.Status.ReadyReplicas
			v.Status.CurrentReplicas = ss.Status.CurrentReplicas
			v.Status.UpdatedReplicas = ss.Status.UpdatedReplicas
			v.Status.CurrentRevision = ss.Status.CurrentRevision
			v.Status.UpdateRevision = ss.Status.UpdateRevision
			v.Status.CollisionCount = ss.Status.CollisionCount
		}
		t = append(t, v)
	}
	fooCopy.Spec.Applications = t

	// If the CustomResourceSubResources feature gate is not enabled,
	// we must use Update instead of UpdateStatus to update the Status block of the RedisOperator resource.
	// UpdateStatus will not allow changes to the Spec of the resource,
	// which is ideal for ensuring nothing other than resource status has been updated.
	_, err := clientSet.NevercaseV1().HelixSagas(foo.Namespace).Update(fooCopy)
	return err
}

func DeleteStatefulSetAndService(ks k8sCoreV1.KubernetesResource, namespace string, name string) error {
	if err := ks.StatefulSet().Delete(namespace, name); err != nil {
		klog.V(2).Info(err)
		return err
	}
	if err := ks.Service().Delete(namespace, k8sCoreV1.GetServiceName(name)); err != nil {
		klog.V(2).Info(err)
		return err
	}
	return nil
}

func PatchStatefulSet(ks k8sCoreV1.KubernetesResource, client helixSagaClientSet.Interface, hs *helixSagaV1.HelixSaga, spec helixSagaV1.HelixSagaAppSpec) error {
	ss := GetStatefulSetImagePatch(hs, spec)
	data, err := json.Marshal(*ss)
	if err != nil {
		klog.V(2).Info(err)
	}
	ss, err = ks.StatefulSet().Patch(hs.Namespace, hs.Name, types.MergePatchType, data)
	if err != nil {
		klog.V(2).Info(err)
		return err
	}
	if err = updateStatus(hs, client, ss, spec.Name); err != nil {
		return err
	}
	return nil
}
