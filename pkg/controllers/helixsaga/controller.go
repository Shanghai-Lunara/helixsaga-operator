package helixsaga

import (
	"time"

	appsV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog"

	helixSagaV1 "github.com/Shanghai-Lunara/helixsaga-operator/pkg/apis/helixsaga/v1"
	helixSagaClientSet "github.com/Shanghai-Lunara/helixsaga-operator/pkg/generated/helixsaga/clientset/versioned"
	helixSagaScheme "github.com/Shanghai-Lunara/helixsaga-operator/pkg/generated/helixsaga/clientset/versioned/scheme"
	informersExt "github.com/Shanghai-Lunara/helixsaga-operator/pkg/generated/helixsaga/informers/externalversions"
	informers "github.com/Shanghai-Lunara/helixsaga-operator/pkg/generated/helixsaga/informers/externalversions/helixsaga/v1"
	k8sCoreV1 "github.com/nevercase/k8s-controller-custom-resource/core/v1"
)

func NewController(
	controllerName string,
	kubeclientset kubernetes.Interface,
	sampleclientset helixSagaClientSet.Interface,
	stopCh <-chan struct{}) k8sCoreV1.KubernetesControllerV1 {

	exampleInformerFactory := informersExt.NewSharedInformerFactory(sampleclientset, time.Second*30)
	fooInformer := exampleInformerFactory.Helixsaga().V1().HelixSagas()

	//roInformerFactory := informersv2.NewSharedInformerFactory(sampleclientset, time.Second*30)

	opt := k8sCoreV1.NewOption(&helixSagaV1.HelixSaga{},
		controllerName,
		operatorKindName,
		helixSagaScheme.AddToScheme(scheme.Scheme),
		sampleclientset,
		fooInformer,
		fooInformer.Informer().HasSynced,
		fooInformer.Informer().AddEventHandler,
		CompareResourceVersion,
		Get,
		Sync)
	opts := k8sCoreV1.NewOptions()
	if err := opts.Add(opt); err != nil {
		klog.Fatal(err)
	}
	op := k8sCoreV1.NewKubernetesOperator(kubeclientset, stopCh, controllerName, opts)
	kc := k8sCoreV1.NewKubernetesController(op)
	//roInformerFactory.Start(stopCh)
	exampleInformerFactory.Start(stopCh)
	return kc
}

//func NewOption() k8sCoreV1.Option {
//
//}

func CompareResourceVersion(old, new interface{}) bool {
	newDepl := new.(*helixSagaV1.HelixSaga)
	oldDepl := old.(*helixSagaV1.HelixSaga)
	if newDepl.ResourceVersion == oldDepl.ResourceVersion {
		// Periodic resync will send update events for all known Deployments.
		// Two different versions of the same Deployment will always have different RVs.
		return true
	}
	return false
}

func Get(foo interface{}, nameSpace, ownerRefName string) (obj interface{}, err error) {
	kc := foo.(informers.HelixSagaInformer)
	return kc.Lister().HelixSagas(nameSpace).Get(ownerRefName)
}

func Sync(obj interface{}, clientObj interface{}, ks k8sCoreV1.KubernetesResource, recorder record.EventRecorder) error {
	hs := obj.(*helixSagaV1.HelixSaga)
	clientSet := clientObj.(helixSagaClientSet.Interface)
	if err := NewNginxPhpFpm(ks, clientSet, hs, hs.Spec.VersionSpec); err != nil {
		return err
	}
	recorder.Event(hs, coreV1.EventTypeNormal, SuccessSynced, MessageResourceSynced)
	return nil
}

func updateFooStatus(foo *helixSagaV1.HelixSaga, clientSet helixSagaClientSet.Interface, statefulSet *appsV1.StatefulSet) error {
	// NEVER modify objects from the store. It's a read-only, local cache.
	// You can use DeepCopy() to make a deep copy of original object and modify this copy
	// Or create a copy manually for better performance
	fooCopy := foo.DeepCopy()
	fooCopy.Status.VersionStatus.AvailableReplicas = statefulSet.Status.Replicas

	// If the CustomResourceSubResources feature gate is not enabled,
	// we must use Update instead of UpdateStatus to update the Status block of the RedisOperator resource.
	// UpdateStatus will not allow changes to the Spec of the resource,
	// which is ideal for ensuring nothing other than resource status has been updated.
	_, err := clientSet.HelixsagaV1().HelixSagas(foo.Namespace).Update(fooCopy)
	return err
}
