package helixsaga

import (
	"testing"
)

func TestGetLabelSelector(t *testing.T) {
	type args struct {
		controllerName string
		specName       string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "TestGetLabelSelector_1",
			args: args{
				controllerName: "hso-develop",
				specName:       "hso-develop-game",
			},
			want: "app=HelixSaga,controller=hso-develop,name=hso-develop-game",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetLabelSelector(tt.args.controllerName, tt.args.specName); got != tt.want {
				t.Errorf("GetLabelSelector() = (%v), want (%v)", got, tt.want)
			}
		})
	}
}
