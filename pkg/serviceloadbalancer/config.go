package serviceloadbalancer

import (
	"github.com/Shanghai-Lunara/pkg/zaplogger"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type Annotations struct {
	Annotations map[string]string `yaml:"annotations"`
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
	return annotations
}
