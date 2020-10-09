package helixsaga

import (
	"reflect"
	"testing"
)

func TestConvertImageToObject(t *testing.T) {
	type args struct {
		image string
	}
	tests := []struct {
		name string
		args args
		want *ImageInfo
	}{
		{
			name: "TestConvertImageToObject_1",
			args: args{
				image: "harbor.domain.com/helix-saga/go-all:latest",
			},
			want: &ImageInfo{
				Domain:     "harbor.domain.com",
				Project:    "helix-saga",
				Repository: "go-all",
				Tag:        "latest",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertImageToObject(tt.args.image); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertImageToObject() = %v, want %v", got, tt.want)
			}
		})
	}
}
