package helixsaga

import (
	"fmt"
	k8sCoreV1 "github.com/nevercase/k8s-controller-custom-resource/core/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
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

var fakeImage = "harbor.domain.com/fake-project/box:latest"
var fakeImageID = "docker-pullable://harbor.domain.com/fake-project/box@sha256:d69e015d92a51c351b2c621ada4b3bfe250752dbe99c7f29d2b8118f60a5ef24"

var fakeNamespace1 = "test"
var fakeControllerName1 = "hso-test"

var fakeHelixSagaAppSpecName1 = "hso-test-game"

var fakePodSpecName1 = "hso-test-game-0"
var fakePodContainerName1 = "hso-test-game-0"

var fakePodSpecName2 = "hso-test-game-1"
var fakePodContainerName2 = "hso-test-game-1"

var fakePodSpecName3 = "hso-test-gmt-1"

var fakePodList = []runtime.Object{
	&corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fakePodSpecName1,
			Namespace: fakeNamespace1,
			Labels: map[string]string{
				k8sCoreV1.LabelApp:        OperatorKindName,
				k8sCoreV1.LabelController: fakeControllerName1,
				k8sCoreV1.LabelName:       fakeHelixSagaAppSpecName1,
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  fakePodContainerName1,
					Image: fakeImage,
				},
			},
		},
		Status: corev1.PodStatus{
			ContainerStatuses: []corev1.ContainerStatus{
				{
					Name:    fakePodContainerName1,
					Image:   fakeImage,
					ImageID: fakeImageID,
				},
			},
			Phase: corev1.PodPending,
		},
	},
	&corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fakePodSpecName2,
			Namespace: fakeNamespace1,
			Labels: map[string]string{
				k8sCoreV1.LabelApp:        OperatorKindName,
				k8sCoreV1.LabelController: fakeControllerName1,
				k8sCoreV1.LabelName:       fakeHelixSagaAppSpecName1,
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  fakePodContainerName2,
					Image: fakeImage,
				},
			},
		},
		Status: corev1.PodStatus{
			ContainerStatuses: []corev1.ContainerStatus{
				{
					Name:    fakePodContainerName2,
					Image:   fakeImage,
					ImageID: fakeImageID,
				},
			},
			Phase: corev1.PodReasonUnschedulable,
		},
	},
	&corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fakePodSpecName3,
			Namespace: fakeNamespace1,
			Labels: map[string]string{
				k8sCoreV1.LabelApp:        OperatorKindName,
				k8sCoreV1.LabelController: fakeControllerName1,
				k8sCoreV1.LabelName:       "2121212121",
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  fakePodContainerName2,
					Image: fakeImage,
				},
			},
		},
	},
}

func TestListPodByLabels(t *testing.T) {
	fakeInterface := fake.NewSimpleClientset(fakePodList...)
	type args struct {
		ki             kubernetes.Interface
		namespace      string
		controllerName string
		specName       string
	}
	tests := []struct {
		name    string
		args    args
		want    *corev1.PodList
		wantErr bool
	}{
		{
			name: "TestListPodByLabels_1",
			args: args{
				ki:             fakeInterface,
				namespace:      fakeNamespace1,
				controllerName: fakeControllerName1,
				specName:       fakeHelixSagaAppSpecName1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ListPodByLabels(tt.args.ki, tt.args.namespace, tt.args.controllerName, tt.args.specName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListPodByLabels() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil || got.Items == nil || len(got.Items) == 0 {
				t.Error("No PodList")
				return
			}
			for _, v := range got.Items {
				fmt.Println("name:", v.Name, " imageId:", v.Status.ContainerStatuses[0].ImageID)
			}
			if len(got.Items) != 2 {
				t.Error("ListPodByLabels returns error")
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("ListPodByLabels() = %v, want %v", got, tt.want)
			//}
		})
	}
}
