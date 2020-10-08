package helixsaga

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/watch"
	"sync"

	harbor "github.com/nevercase/harbor-api"
)

type Watchers struct {
	mu sync.Mutex

	items map[string]*Watcher

	harborHub harbor.HubInterface
}

func NewWatchers(c []harbor.Config) *Watchers {
	return &Watchers{
		items:     make(map[string]*Watcher, 0),
		harborHub: harbor.NewHub(c),
	}
}

func (ws *Watchers) Subscribe(wo *WatchOption) {
	ws.mu.Lock()
	defer ws.mu.Unlock()
	name := wo.Name()
	if _, ok := ws.items[name]; !ok {
		ws.items[name] = NewWatcher(context.Background(), ws.harborHub, wo)
	}
}

func (ws *Watchers) UnSubscribe(wo *WatchOption) error {
	ws.mu.Lock()
	defer ws.mu.Unlock()
	name := wo.Name()
	if t, ok := ws.items[name]; ok {
		t.Close()
	}
	delete(ws.items, name)
	return nil
}

type Watcher struct {
	harborHub harbor.HubInterface

	name        string
	opt         *WatchOption
	watchResult watch.Interface

	ctx    context.Context
	cancel context.CancelFunc
}

func NewWatcher(ctx context.Context, harborHub harbor.HubInterface, wo *WatchOption) *Watcher {
	subCtx, cancel := context.WithCancel(ctx)
	return &Watcher{
		harborHub: harborHub,
		name:      wo.Name(),
		opt:       wo,
		ctx:       subCtx,
		cancel:    cancel,
	}
}

func (w *Watcher) Loop() {
	defer w.watchResult.Stop()
	for {
		select {
		case <-w.ctx.Done():
			return
		case msg, isClose := <-w.watchResult.ResultChan():
			if !isClose {
				// todo reWatch
			}
			harborWatch := msg.Object.(*harbor.Option)
			// todo handle the message which was received from the watch channel
			_ = harborWatch
		}
	}
}

func (w *Watcher) Close() {
	w.cancel()
}

type WatchOption struct {
	Namespace    string
	OperatorName string
	SpecName     string
	Image        string
	ImageID      string
}

func (wo *WatchOption) Name() string {
	return fmt.Sprintf("%s-%s-%s-%s", wo.Namespace, wo.OperatorName, wo.SpecName, wo.Image)
}

// HarborUrl returns the image's harbor domain
func (wo *WatchOption) HarborUrl() string {
	// image: harbor.domain.com/helix-saga/go-all:latest
	// image: xxxx.xxx.xxx:443/helix-saga/go-all:latest
	// image: http://xxxx.xxx.xxx:80
	// image: https://xxxx.xxx.xxx
	return ""
}

func WatchHarborImage(harbor harbor.HubInterface, wo *WatchOption) (watch.Interface, error) {
	return watch.NewEmptyWatch(), nil
}