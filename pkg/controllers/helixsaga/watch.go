package helixsaga

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	helixsagav1 "github.com/Shanghai-Lunara/helixsaga-operator/pkg/apis/helixsaga/v1"
	helixSagaClientSet "github.com/Shanghai-Lunara/helixsaga-operator/pkg/generated/helixsaga/clientset/versioned"
	harbor "github.com/nevercase/harbor-api"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog"
)

// Watchers was a set which contains all the watchers
type Watchers struct {
	mu        sync.Mutex
	items     map[string]*Watcher
	harborHub harbor.HubInterface
}

// NewWatchers returns the pointer of the Watchers
func NewWatchers(c []harbor.Config) *Watchers {
	return &Watchers{
		items:     make(map[string]*Watcher, 0),
		harborHub: harbor.NewHub(c),
	}
}

// Subscribe adds a new Watcher into the Watchers
func (ws *Watchers) Subscribe(wo *WatchOption) error {
	ws.mu.Lock()
	defer ws.mu.Unlock()
	name := wo.Name()
	if _, ok := ws.items[name]; !ok {
		w, err := NewWatcher(context.Background(), ws.harborHub, wo)
		if err != nil {
			klog.V(2).Info(err)
			return err
		}
		ws.items[name] = w
	}
	// update WatchOption
	ws.items[name].opt = wo
	return nil
}

// UnSubscribe removes the specific Watcher from the Watchers
func (ws *Watchers) UnSubscribe(wo *WatchOption) {
	ws.mu.Lock()
	defer ws.mu.Unlock()
	name := wo.Name()
	if t, ok := ws.items[name]; ok {
		t.Close()
	}
	delete(ws.items, name)
}

type Watcher struct {
	harborHub harbor.HubInterface

	name        string
	opt         *WatchOption
	watchResult watch.Interface

	once   sync.Once
	ctx    context.Context
	cancel context.CancelFunc
}

func NewWatcher(ctx context.Context, hi harbor.HubInterface, wo *WatchOption) (*Watcher, error) {
	wi, err := WatchHarborImage(hi, wo)
	if err != nil {
		klog.V(2).Info(err)
		return nil, err
	}
	subCtx, cancel := context.WithCancel(ctx)
	w := &Watcher{
		harborHub:   hi,
		name:        wo.Name(),
		opt:         wo,
		watchResult: wi,
		ctx:         subCtx,
		cancel:      cancel,
	}
	go w.Loop()
	return w, nil
}

func (w *Watcher) Loop() {
	defer w.watchResult.Stop()
	for {
		select {
		case <-w.ctx.Done():
			return
		case msg, isClose := <-w.watchResult.ResultChan():
			klog.Info("Watcher Loop msg:", msg)
			if !isClose {
				// reWatch
				tick := time.NewTicker(time.Millisecond * 200)
				for {
					select {
					case <-w.ctx.Done():
						return
					case <-tick.C:
						klog.V(2).Infof("Watcher ResultChan reconnect name:%s", w.name)
						wi, err := WatchHarborImage(w.harborHub, w.opt)
						if err != nil {
							klog.V(2).Info(err)
							continue
						}
						w.watchResult = wi
						break
					}
				}
			}
			t := msg.Object.(harbor.Option)
			// handle the message which was received from the watch channel
			image, err := w.opt.GetPodImage()
			if err != nil {
				klog.V(2).Info(err)
				// todo handle error
				continue
			}
			hash := harbor.GetHashFromDockerImageId(image)
			if hash == "" {
				klog.V(2).Infof("image:%s hash:%s", image, hash)
				continue
			}
			if hash != t.Sha256 {
				//if err := PatchStatefulSet(w.opt.K8sClientSet, w.opt.HelixSagaClient, w.opt.HelixSaga, w.opt.SpecName, w.opt.Image); err != nil {
				//	klog.V(2).Info(err)
				//}
				//if err := PatchPod(w.opt.K8sClientSet, w.opt.HelixSaga.Namespace, w.opt.HelixSaga.Name, w.opt.SpecName); err != nil {
				//	klog.V(2).Info(err)
				//}
				var replica int32
				if replica, err = RetryPatchHelixSaga(w.opt.K8sClientSet, w.opt.HelixSagaClient, w.opt.HelixSaga.Namespace, w.opt.HelixSaga.Name, w.opt.SpecName, 0); err != nil {
					klog.V(2).Info(err)
					continue
				}
				if _, err = RetryPatchHelixSaga(w.opt.K8sClientSet, w.opt.HelixSagaClient, w.opt.HelixSaga.Namespace, w.opt.HelixSaga.Name, w.opt.SpecName, replica); err != nil {
					klog.V(2).Info(err)
					continue
				}
			}
		}
	}
}

func (w *Watcher) Close() {
	w.once.Do(func() {
		w.cancel()
		w.opt.Close()
	})
}

type WatchOption struct {
	Namespace    string
	OperatorName string
	SpecName     string
	Image        string
	ImageID      string
	ImageInfo    *ImageInfo

	K8sClientSet    kubernetes.Interface
	HelixSagaClient helixSagaClientSet.Interface
	StatefulSet     *appsv1.StatefulSet
	HelixSaga       *helixsagav1.HelixSaga

	ctx    context.Context
	cancel context.CancelFunc
}

type ImageInfo struct {
	Domain     string
	Project    string
	Repository string
	Tag        string
}

func NewWatchOption(ctx context.Context, ki kubernetes.Interface, helixSagaClient helixSagaClientSet.Interface, hs *helixsagav1.HelixSaga, specName, image string) *WatchOption {
	sub, cancel := context.WithCancel(ctx)
	wo := &WatchOption{
		Namespace:       hs.Namespace,
		OperatorName:    hs.Name,
		SpecName:        specName,
		Image:           image,
		ImageInfo:       ConvertImageToObject(image),
		K8sClientSet:    ki,
		HelixSagaClient: helixSagaClient,
		HelixSaga:       hs,
		ctx:             sub,
		cancel:          cancel,
	}
	return wo
}

func (wo *WatchOption) Name() string {
	return fmt.Sprintf("%s-%s-%s-%s", wo.Namespace, wo.OperatorName, wo.SpecName, wo.Image)
}

func (wo *WatchOption) GetPodImage() (string, error) {
	tick := time.NewTicker(time.Millisecond * 500)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			pl, err := ListPodByLabels(wo.K8sClientSet, wo.Namespace, wo.OperatorName, wo.SpecName)
			if err != nil {
				klog.V(2).Infof("WatchOption GetPodImage ListPodByLabels err:%v", err)
				continue
			}
			var hash string
			for _, v := range pl.Items {
				if v.Status.Phase != corev1.PodRunning {
					klog.Infof("Pod name:%s Status.Phase:%s", v.Name, v.Status.Phase)
					continue
				}
				if len(v.Status.ContainerStatuses) == 0 {
					klog.Infof("Pod name:%s ContainerStatuses was empty", v.Name)
					continue
				}
				if v.Status.ContainerStatuses[0].Image != wo.Image {
					klog.Infof("Pod name:%s image:%s was not match the WatchOption image:%s", v.Name, v.Status.ContainerStatuses[0].Image, wo.Image)
					continue
				}
				if v.Status.ContainerStatuses[0].ImageID == "" {
					klog.Infof("Pod name:%s ContainerStatuses Container:%s image:%s ImageID was empty",
						v.Name, v.Status.ContainerStatuses[0].Name, v.Status.ContainerStatuses[0].Image)
					continue
				}
				hash = v.Status.ContainerStatuses[0].ImageID
				break
			}
			if hash == "" {
				klog.Infof("WatchOption GetPodImage HelixSaga Name")
				continue
			}
			return hash, nil
		case <-wo.ctx.Done():
			return "", fmt.Errorf("WatchOption GetPodImage ctx cancel")
		}
	}
}

func (wo *WatchOption) Close() {
	wo.cancel()
}

// ConvertImageToObject returns the ImageInfo with an image string
func ConvertImageToObject(image string) *ImageInfo {
	// image: harbor.domain.com/helix-saga/go-all:latest
	// image: xxxx.xxx.xxx.xxx:8080/helix-saga/go-all:latest
	t := strings.Split(image, "/")
	t2 := strings.Split(t[2], ":")
	return &ImageInfo{
		Domain:     t[0],
		Project:    t[1],
		Repository: t2[0],
		Tag:        t2[1],
	}
}

func WatchHarborImage(hi harbor.HubInterface, wo *WatchOption) (watch.Interface, error) {
	hb, err := hi.Get(fmt.Sprintf("%s%s", harbor.HttpPrefix, wo.ImageInfo.Domain))
	if err != nil {
		klog.V(2).Info(err)
		return nil, err
	}
	if err := hb.Login(); err != nil {
		klog.V(2).Info(err)
		return nil, err
	}
	opt := harbor.Option{
		APIVersion: "v1",
		Kind:       "",
		Project:    wo.ImageInfo.Project,
		Repository: wo.ImageInfo.Repository,
		Tag:        wo.ImageInfo.Tag,
	}
	var t watch.Interface
	if t, err = hb.Watch(opt); err != nil {
		klog.V(2).Info(err)
		return nil, err
	}
	return t, nil
}
