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
	c := &Annotations{}
	var data []byte
	var err error
	if data, err = ioutil.ReadFile(configFile); err != nil {
		zaplogger.Sugar().Fatal(err)
	}
	if err := yaml.Unmarshal(data, c); err != nil {
		zaplogger.Sugar().Fatal(err)
	}
	return c
}

func Get() *Annotations {
	return annotations
}
