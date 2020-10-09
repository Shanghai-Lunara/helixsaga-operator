package helixsaga

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	harbor "github.com/nevercase/harbor-api"
	"k8s.io/apimachinery/pkg/watch"
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
	return &Watcher{
		harborHub:   hi,
		name:        wo.Name(),
		opt:         wo,
		watchResult: wi,
		ctx:         subCtx,
		cancel:      cancel,
	}, nil
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
				// todo reWatch
				tick := time.NewTicker(time.Millisecond * 200)
				for {
					select {
					case <-tick.C:
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
			harborWatch := msg.Object.(*harbor.Option)
			// todo handle the message which was received from the watch channel
			_ = harborWatch
		}
	}
}

func (w *Watcher) Close() {
	w.once.Do(func() {
		w.cancel()
	})
}

type WatchOption struct {
	Namespace    string
	OperatorName string
	SpecName     string
	Image        string
	ImageID      string
	ImageInfo    *ImageInfo
}

type ImageInfo struct {
	Domain     string
	Project    string
	Repository string
	Tag        string
}

func (wo *WatchOption) Convert() {
	wo.ImageInfo = ConvertImageToObject(wo.Image)
}

func (wo *WatchOption) Name() string {
	return fmt.Sprintf("%s-%s-%s-%s", wo.Namespace, wo.OperatorName, wo.SpecName, wo.Image)
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
	hb, err := hi.Get(wo.ImageInfo.Domain)
	if err != nil {
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
