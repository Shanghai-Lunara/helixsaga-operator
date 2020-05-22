package helixsaga

import (
	"time"

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
	newResource := new.(*helixSagaV1.HelixSaga)
	oldResource := old.(*helixSagaV1.HelixSaga)
	if newResource.ResourceVersion == oldResource.ResourceVersion {
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
	for _, v := range hs.Spec.Services {
		klog.Info("v:", v)
		if err := NewStatefulSetAndService(ks, clientSet, hs, v.Spec); err != nil {
			klog.V(2).Info(err)
			return err
		}
	}
	recorder.Event(hs, coreV1.EventTypeNormal, SuccessSynced, MessageResourceSynced)
	return nil
}