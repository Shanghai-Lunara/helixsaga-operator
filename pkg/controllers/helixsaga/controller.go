package helixsaga

import (
	"context"
	"fmt"
	"sync"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog"

	helixsagav1 "github.com/Shanghai-Lunara/helixsaga-operator/pkg/apis/helixsaga/v1"
	helixsagaclientset "github.com/Shanghai-Lunara/helixsaga-operator/pkg/generated/helixsaga/clientset/versioned"
	helixsagascheme "github.com/Shanghai-Lunara/helixsaga-operator/pkg/generated/helixsaga/clientset/versioned/scheme"
	informersext "github.com/Shanghai-Lunara/helixsaga-operator/pkg/generated/helixsaga/informers/externalversions"
	informers "github.com/Shanghai-Lunara/helixsaga-operator/pkg/generated/helixsaga/informers/externalversions/helixsaga/v1"
	harbor "github.com/nevercase/harbor-api"
	k8scorev1 "github.com/nevercase/k8s-controller-custom-resource/core/v1"
)

func NewController(
	controllerName string,
	kubeclientset kubernetes.Interface,
	sampleclientset helixsagaclientset.Interface,
	stopCh <-chan struct{}) k8scorev1.KubernetesControllerV1 {

	controller := &controller{}

	exampleInformerFactory := informersext.NewSharedInformerFactory(sampleclientset, time.Second*30)
	fooInformer := exampleInformerFactory.Nevercase().V1().HelixSagas()
	//roInformerFactory := informersv2.NewSharedInformerFactory(sampleclientset, time.Second*30)

	opt := k8scorev1.NewOption(&helixsagav1.HelixSaga{},
		controllerName,
		OperatorKindName,
		helixsagascheme.AddToScheme(scheme.Scheme),
		sampleclientset,
		fooInformer,
		fooInformer.Informer(),
		controller.CompareResourceVersion,
		controller.Get,
		controller.Sync,
		controller.SyncStatus)
	opts := k8scorev1.NewOptions()
	if err := opts.Add(opt); err != nil {
		klog.Fatal(err)
	}
	op := k8scorev1.NewKubernetesOperator(kubeclientset, stopCh, controllerName, opts)
	kc := k8scorev1.NewKubernetesController(op)
	//roInformerFactory.Start(stopCh)
	exampleInformerFactory.Start(stopCh)
	return kc
}

func NewOption(controllerName string, cfg *rest.Config, stopCh <-chan struct{}, harborConfig []harbor.Config) k8scorev1.Option {
	c, err := helixsagaclientset.NewForConfig(cfg)
	if err != nil {
		klog.Fatalf("Error building clientSet: %s", err.Error())
	}
	controller := &controller{
		watchers: NewWatchers(harborConfig),
	}
	informerFactory := informersext.NewSharedInformerFactory(c, time.Second*30)
	fooInformer := informerFactory.Nevercase().V1().HelixSagas()
	opt := k8scorev1.NewOption(&helixsagav1.HelixSaga{},
		controllerName,
		OperatorKindName,
		helixsagascheme.AddToScheme(scheme.Scheme),
		c,
		fooInformer,
		fooInformer.Informer(),
		controller.CompareResourceVersion,
		controller.Get,
		controller.Sync,
		controller.SyncStatus)
	informerFactory.Start(stopCh)
	return opt
}

type controller struct {
	mu       sync.Mutex
	watchers *Watchers
}

func (c *controller) CompareResourceVersion(old, new interface{}) bool {
	newResource := new.(*helixsagav1.HelixSaga)
	oldResource := old.(*helixsagav1.HelixSaga)
	if newResource.ResourceVersion == oldResource.ResourceVersion {
		// Periodic resync will send update events for all known HelixSaga.
		// Two different versions of the same HelixSaga will always have different RVs.
		return true
	}
	return false
}

func (c *controller) Get(foo interface{}, nameSpace, ownerRefName string) (obj interface{}, err error) {
	kc := foo.(informers.HelixSagaInformer)
	return kc.Lister().HelixSagas(nameSpace).Get(ownerRefName)
}

func (c *controller) Sync(obj interface{}, clientObj interface{}, ks k8scorev1.KubernetesResource, recorder record.EventRecorder) error {
	hs := obj.(*helixsagav1.HelixSaga)
	clientSet := clientObj.(helixsagaclientset.Interface)
	opts := metav1.GetOptions{ResourceVersion: hs.ResourceVersion}
	old, err := clientSet.NevercaseV1().HelixSagas(hs.Namespace).Get(hs.Name, opts)
	if err != nil {
		if !errors.IsNotFound(err) {
			klog.V(2).Info(err)
			return err
		} else {
			names := make(map[string]bool, len(hs.Spec.Applications))
			images := make(map[string]int, 0)
			for _, v := range hs.Spec.Applications {
				names[v.Spec.Name] = true
				if _, ok := images[v.Spec.Name]; ok {
					images[v.Spec.Name] += 1
				} else {
					images[v.Spec.Name] = 1
				}
			}
			for _, v := range old.Spec.Applications {
				if _, ok := names[v.Spec.Name]; !ok {
					// stop watching before removing apps
					wo := &WatchOption{
						Namespace:    hs.Namespace,
						OperatorName: hs.Name,
						Image:        v.Spec.Image,
					}
					if _, ok := images[v.Spec.Image]; !ok {
						klog.Infof("HelixSaga crdName:%s image:%s has been removed", hs.Name, v.Spec.Image)
						c.watchers.UnSubscribe(wo)
					}
					klog.Info("remove app-name:", v.Spec.Name)
					if err := DeleteStatefulSetAndService(ks, hs.Namespace, v.Spec.Name); err != nil {
						klog.V(2).Info(err)
						return err
					}
				}
			}
		}
	}
	for _, v := range hs.Spec.Applications {
		// starting watching the harbor before creating apps
		wo := NewWatchOption(context.Background(), ks.ClientSet(), clientSet, hs, v.Spec.Image)
		if *v.Spec.Replicas > 0 {
			if err := c.watchers.Subscribe(wo); err != nil {
				klog.V(2).Info(err)
				return err
			}
		} else {
			klog.V(4).Infof("HelixSaga crdName:%s image:%s UnSubscribe due to replicas 0", hs.Name, v.Spec.Image)
			c.watchers.UnSubscribe(wo)
		}
		klog.Info("v:", v)
		if err := NewStatefulSetAndService(ks, clientSet, hs, v.Spec, wo); err != nil {
			klog.V(2).Info(err)
			return err
		}
	}
	recorder.Event(hs, corev1.EventTypeNormal, SuccessSynced, MessageResourceSynced)
	return nil
}

func (c *controller) SyncStatus(obj interface{}, clientObj interface{}, ks k8scorev1.KubernetesResource, recorder record.EventRecorder) error {
	clientSet := clientObj.(helixsagaclientset.Interface)
	ss := obj.(*appsv1.StatefulSet)
	var objName string
	if t, ok := ss.Labels[k8scorev1.LabelController]; ok {
		objName = t
	} else {
		return fmt.Errorf(ErrResourceNotMatch, "no controller")
	}
	hs, err := clientSet.NevercaseV1().HelixSagas(ss.Namespace).Get(objName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	var appName string
	if t, ok := ss.Labels[k8scorev1.LabelName]; ok {
		appName = t
	} else {
		return fmt.Errorf(ErrResourceNotMatch, "no appName")
	}
	if err := updateStatus(hs, clientSet, ss, appName); err != nil {
		return err
	}
	recorder.Event(hs, corev1.EventTypeNormal, SuccessSynced, MessageResourceSynced)
	return nil
}
