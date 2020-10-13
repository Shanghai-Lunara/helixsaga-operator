package helixsaga

import (
	"fmt"
	"k8s.io/apimachinery/pkg/types"
	"sync"

	k8sCoreV1 "github.com/nevercase/k8s-controller-custom-resource/core/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog"
)

// ListPodByLabels
func ListPodByLabels(ki kubernetes.Interface, namespace, controllerName, specName string) (*corev1.PodList, error) {
	timeout := int64(10)
	opts := metav1.ListOptions{
		LabelSelector:  GetLabelSelector(controllerName, specName),
		FieldSelector:  fields.OneTermEqualSelector("status.phase", string(corev1.PodRunning)).String(),
		TimeoutSeconds: &timeout,
	}
	fmt.Println(fields.OneTermEqualSelector("status.phase", string(corev1.PodRunning)).String())
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

func PatchPod(ki kubernetes.Interface, namespace, controllerName, specName string) error {
	timeout := int64(10)
	opts := metav1.ListOptions{
		LabelSelector:  GetLabelSelector(controllerName, specName),
		TimeoutSeconds: &timeout,
	}
	pl, err := ki.CoreV1().Pods(namespace).List(opts)
	if err != nil {
		klog.V(2).Info(err)
		return err
	}
	var wg sync.WaitGroup
	wg.Add(len(pl.Items))
	for _, v := range pl.Items {
		go func(v *corev1.Pod) {
			klog.Infof("patch pod name:%s", v.Name)
			defer wg.Done()
			data, err := GePodImagePatch(v.Namespace, v.Name, v.Spec.Containers[0].Image)
			if err != nil {
				klog.V(2).Info(err)
				return
			}
			_, err = ki.CoreV1().Pods(v.Name).Patch(v.Name, types.MergePatchType, data)
			if err != nil {
				klog.V(2).Info(err)
				return
			}
		}(&v)
	}
	wg.Wait()
	return nil
}

func GePodImagePatch(ns, specName, image string) ([]byte, error) {
	patch := map[string]interface{}{
		"metadata": map[string]interface{}{
			"name":      k8sCoreV1.GetStatefulSetName(specName),
			"namespace": ns,
		},
		"spec": map[string]interface{}{
			"containers": []map[string]interface{}{
				{
					"name":  k8sCoreV1.GetContainerName(specName),
					"image": image,
				},
			},
		},
	}
	return json.Marshal(patch)
}
