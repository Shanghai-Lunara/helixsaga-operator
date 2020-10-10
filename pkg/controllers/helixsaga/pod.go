package helixsaga

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog"

	k8sCoreV1 "github.com/nevercase/k8s-controller-custom-resource/core/v1"
)

// ListPodByLabels
func ListPodByLabels(ki kubernetes.Interface, namespace, controllerName, specName string) (*corev1.PodList, error) {
	timeout := int64(10)
	opts := metav1.ListOptions{
		LabelSelector:  GetLabelSelector(controllerName, specName),
		TimeoutSeconds: &timeout,
	}
	return ki.CoreV1().Pods(namespace).List(opts)
}

// GetLabelSelector returns the LabelSelector of the metav1.ListOptions
func GetLabelSelector(controllerName string, specName string) string {
	req1, err := labels.NewRequirement(k8sCoreV1.LabelApp, selection.Equals, []string{OperatorKindName})
	if err != nil {
		klog.Fatal(err)
	}
	req2, err := labels.NewRequirement(k8sCoreV1.LabelController, selection.Equals, []string{controllerName})
	if err != nil {
		klog.Fatal(err)
	}
	req3, err := labels.NewRequirement(k8sCoreV1.LabelName, selection.Equals, []string{specName})
	if err != nil {
		klog.Fatal(err)
	}
	ls := labels.NewSelector()
	ls = ls.Add(*req1, *req2, *req3)
	return ls.String()
}

