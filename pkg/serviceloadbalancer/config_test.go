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
					"service.beta.kubernetes.io/alibaba-cloud-loadbalancer-id":                  "abc111",
					"service.beta.kubernetes.io/alicloud-loadbalancer-force-override-listeners": "false",
				},
			},
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

func TestLabels(t *testing.T) {
	tests := []struct {
		name string
		want *Annotations
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Labels(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Labels() = %v, want %v", got, tt.want)
			}
		})
	}
}
