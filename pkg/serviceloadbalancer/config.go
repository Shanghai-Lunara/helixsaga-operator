package serviceloadbalancer

import (
	"github.com/Shanghai-Lunara/pkg/zaplogger"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type Annotations struct {
	Annotations map[string]string `yaml:"annotations"`
	WhiteList   map[string]string `yaml:"whiteList"`
}

var annotations *Annotations

func Init(configFile string) *Annotations {
	annotations = &Annotations{}
	var data []byte
	var err error
	if data, err = ioutil.ReadFile(configFile); err != nil {
		zaplogger.Sugar().Fatal(err)
	}
	if err := yaml.Unmarshal(data, annotations); err != nil {
		zaplogger.Sugar().Fatal(err)
	}
	return annotations
}

func Get() *Annotations {
	if annotations == nil {
		annotations = &Annotations{
			Annotations: make(map[string]string, 0),
			WhiteList:   make(map[string]string, 0),
		}
	}
	return annotations
}
