package helixsaga

import (
	"sync"

	harbor "github.com/nevercase/harbor-api"
	"k8s.io/apimachinery/pkg/watch"
)

type watchers struct {
	mu sync.Mutex

	items map[string]watch.Interface

	harborHub harbor.HubInterface
}

func (w *watchers) Subscribe(s string) (watch.Interface, error) {


}

func (w *watchers) UnSubscribe(s string) error {


	return nil
}

type one struct {
	Namespace string
	Name      string

	SpecName string
	Image    string
}
