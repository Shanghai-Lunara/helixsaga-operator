package helixsaga

import (
	"context"
	k8sCoreV1 "github.com/nevercase/k8s-controller-custom-resource/core/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
)

// ListPodByLabels
func ListPodByLabels(ki kubernetes.Interface, namespace, controllerName, specName string) (*corev1.PodList, error) {
	timeout := int64(10)
	opts := metav1.ListOptions{
		LabelSelector: GetLabelSelector(controllerName, specName),
		//FieldSelector:  fields.OneTermEqualSelector("status.phase", string(corev1.PodRunning)).String(),
		TimeoutSeconds: &timeout,
	}
	//fmt.Println(fields.OneTermEqualSelector("status.phase", string(corev1.PodRunning)).String())
	return ki.CoreV1().Pods(namespace).List(context.Background(), opts)
}

// GetLabelSelector returns the LabelSelector of the metav1.ListOptions
func GetLabelSelector(controllerName string, specName string) string {
	ls := labels.NewSelector()
	req1, err := labels.NewRequirement(k8sCoreV1.LabelApp, selection.Equals, []string{OperatorKindName})
	if err != nil {
		klog.Fatal(err)
	}
	ls = ls.Add(*req1)
	if controllerName != "" {
		req2, err := labels.NewRequirement(k8sCoreV1.LabelController, selection.Equals, []string{controllerName})
		if err != nil {
			klog.Fatal(err)
		}
		ls = ls.Add(*req2)
	}
	if specName != "" {
		req3, err := labels.NewRequirement(k8sCoreV1.LabelName, selection.Equals, []string{specName})
		if err != nil {
			klog.Fatal(err)
		}
		ls = ls.Add(*req3)
	}
	return ls.String()
}

// ExposePodInformationByEnvs appends envs to expose more information about the current pod
func ExposePodInformationByEnvs(envs []corev1.EnvVar) []corev1.EnvVar {
	envs = append(envs, corev1.EnvVar{
		Name: NodeName,
		ValueFrom: &corev1.EnvVarSource{
			FieldRef: &corev1.ObjectFieldSelector{
				FieldPath: "spec.nodeName",
			},
		},
	})
	envs = append(envs, corev1.EnvVar{
		Name: HostIP,
		ValueFrom: &corev1.EnvVarSource{
			FieldRef: &corev1.ObjectFieldSelector{
				FieldPath: "status.hostIP",
			},
		},
	})
	envs = append(envs, corev1.EnvVar{
		Name: PodName,
		ValueFrom: &corev1.EnvVarSource{
			FieldRef: &corev1.ObjectFieldSelector{
				FieldPath: "metadata.name",
			},
		},
	})
	envs = append(envs, corev1.EnvVar{
		Name: PodNamespace,
		ValueFrom: &corev1.EnvVarSource{
			FieldRef: &corev1.ObjectFieldSelector{
				FieldPath: "metadata.namespace",
			},
		},
	})
	envs = append(envs, corev1.EnvVar{
		Name: PodIP,
		ValueFrom: &corev1.EnvVarSource{
			FieldRef: &corev1.ObjectFieldSelector{
				FieldPath: "status.podIP",
			},
		},
	})
	envs = append(envs, corev1.EnvVar{
		Name: PodServiceAccount,
		ValueFrom: &corev1.EnvVarSource{
			FieldRef: &corev1.ObjectFieldSelector{
				FieldPath: "spec.serviceAccountName",
			},
		},
	})
	return envs
}
