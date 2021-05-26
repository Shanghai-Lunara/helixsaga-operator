package serviceloadbalancer

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

func TestInit(t *testing.T) {
	path, err := os.Getwd()
	if err != nil {
		t.Errorf("os.Getwd err:%v\n", err)
	}
	type args struct {
		configFile string
	}
	tests := []struct {
		name           string
		args           args
		want           *Annotations
		expectedResult bool
	}{
		{
			name: "TestInit_case1",
			args: args{
				configFile: fmt.Sprintf("%s/svc.yaml", path),
			},
			want: &Annotations{
				Annotations: map[string]string{
					"service.beta.kubernetes.io/alibaba-cloud-loadbalancer-id":                  "abc",
					"service.beta.kubernetes.io/alicloud-loadbalancer-force-override-listeners": "false",
				},
				WhiteListOn: map[string]string{
					"service.beta.kubernetes.io/alibaba-cloud-loadbalancer-acl-status": "on",
					"service.beta.kubernetes.io/alibaba-cloud-loadbalancer-acl-id":     "${YOUR_ACL_ID}",
					"service.beta.kubernetes.io/alibaba-cloud-loadbalancer-acl-type":   "white",
				},
				WhiteListOff: map[string]string{
					"service.beta.kubernetes.io/alibaba-cloud-loadbalancer-acl-status": "off",
				},
			},
			expectedResult: true,
		},
		{
			name: "TestInit_case2",
			args: args{
				configFile: fmt.Sprintf("%s/svc.yaml", path),
			},
			want: &Annotations{
				Annotations: map[string]string{
					"service.beta.kubernetes.io/alibaba-cloud-loadbalancer-id":                  "abc",
					"service.beta.kubernetes.io/alicloud-loadbalancer-force-override-listeners": "false",
				},
				WhiteListOn: map[string]string{
					"service.beta.kubernetes.io/alibaba-cloud-loadbalancer-acl-statusx": "on",
					"service.beta.kubernetes.io/alibaba-cloud-loadbalancer-acl-id":      "${YOUR_ACL_ID}",
					"service.beta.kubernetes.io/alibaba-cloud-loadbalancer-acl-type":    "white1",
				},
				WhiteListOff: map[string]string{
					"service.beta.kubernetes.io/alibaba-cloud-loadbalancer-acl-status": "off",
				},
			},
			expectedResult: false,
		},
		{
			name: "TestInit_case3",
			args: args{
				configFile: fmt.Sprintf("%s/svc.yaml", path),
			},
			want: &Annotations{
				Annotations: map[string]string{
					"service.beta.kubernetes.io/alibaba-cloud-loadbalancer-id":                  "abc111",
					"service.beta.kubernetes.io/alicloud-loadbalancer-force-override-listeners": "false",
				},
				WhiteListOn: map[string]string{
					"service.beta.kubernetes.io/alibaba-cloud-loadbalancer-acl-status": "on",
					"service.beta.kubernetes.io/alibaba-cloud-loadbalancer-acl-id":     "${YOUR_ACL_ID}",
					"service.beta.kubernetes.io/alibaba-cloud-loadbalancer-acl-type":   "white",
				},
			},
			expectedResult: false,
		},
		{
			name: "TestInit_case4",
			args: args{
				configFile: fmt.Sprintf("%s/svc.yaml", path),
			},
			want:           &Annotations{},
			expectedResult: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.expectedResult {
			case true:
				if got := Init(tt.args.configFile); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Init() = %v, want %v", got, tt.want)
				}
			case false:
				if got := Init(tt.args.configFile); reflect.DeepEqual(got, tt.want) {
					t.Errorf("Init() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestPointer(t *testing.T) {
	path, err := os.Getwd()
	if err != nil {
		t.Errorf("os.Getwd err:%v\n", err)
	}
	_ = Init(fmt.Sprintf("%s/svc.yaml", path))

	fmt.Printf("init:%#v\n", annotations)
	b1 := make(map[string]string, 0)
	b1 = Get().Annotations
	fmt.Printf("b1-1: %#v\n", b1)
	for k, v := range Get().WhiteListOn {
		b1[k] = v
	}
	fmt.Printf("b1-2: %#v\n", b1)
	fmt.Printf("Annotations: %#v\n", Get().Annotations)
	fmt.Printf("now:%#v\n", annotations)
}
