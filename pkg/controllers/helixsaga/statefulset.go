package helixsaga

import (
	"fmt"

	helixSagaV1 "github.com/Shanghai-Lunara/helixsaga-operator/pkg/apis/helixsaga/v1"
	k8sCoreV1 "github.com/nevercase/k8s-controller-custom-resource/core/v1"
	appsV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/json"
)

func NewStatefulSet(hs *helixSagaV1.HelixSaga, spec helixSagaV1.HelixSagaAppSpec) *appsV1.StatefulSet {
	labels := map[string]string{
		k8sCoreV1.LabelApp:        OperatorKindName,
		k8sCoreV1.LabelController: hs.Name,
		k8sCoreV1.LabelName:       spec.Name,
	}
	t := coreV1.HostPathDirectoryOrCreate
	hostPath := &coreV1.HostPathVolumeSource{
		Type: &t,
		Path: fmt.Sprintf("%s/%s/helixsaga/%s", spec.VolumePath, hs.Namespace, spec.Name),
	}
	return &appsV1.StatefulSet{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      k8sCoreV1.GetStatefulSetName(spec.Name),
			Namespace: hs.Namespace,
			OwnerReferences: []metaV1.OwnerReference{
				*metaV1.NewControllerRef(hs, helixSagaV1.SchemeGroupVersion.WithKind(OperatorKindName)),
			},
			Labels: labels,
		},
		Spec: appsV1.StatefulSetSpec{
			Replicas: spec.Replicas,
			Selector: &metaV1.LabelSelector{
				MatchLabels: labels,
			},
			Template: coreV1.PodTemplateSpec{
				ObjectMeta: metaV1.ObjectMeta{
					Labels: labels,
				},
				Spec: coreV1.PodSpec{
					Volumes: []coreV1.Volume{
						hs.Spec.ConfigMap.Volume,
						{
							Name: "task-pv-storage",
							VolumeSource: coreV1.VolumeSource{
								HostPath: hostPath,
							},
						},
					},
					Containers: []coreV1.Container{
						{
							Name:  k8sCoreV1.GetContainerName(spec.Name),
							Image: spec.Image,
							Ports: spec.ContainerPorts,
							Env:   spec.Env,
							VolumeMounts: []coreV1.VolumeMount{
								hs.Spec.ConfigMap.VolumeMount,
								{
									MountPath: "/data",
									Name:      "task-pv-storage",
								},
							},
							Command:         spec.Command,
							Args:            spec.Args,
							Resources:       spec.Resources,
							ImagePullPolicy: coreV1.PullAlways,
						},
					},
					ImagePullSecrets: spec.ImagePullSecrets,
				},
			},
		},
	}
}

func GetStatefulSetImagePatch(hs *helixSagaV1.HelixSaga, specName, image string) ([]byte, error) {
	patch := map[string]interface{}{
		"metadata": map[string]interface{}{
			"name":      k8sCoreV1.GetStatefulSetName(specName),
			"namespace": hs.Namespace,
		},
		"spec": map[string]interface{}{
			"template": map[string]interface{}{
				"spec": map[string]interface{}{
					"containers": []map[string]interface{}{
						{
							"name":  k8sCoreV1.GetContainerName(specName),
							"image": image,
						},
					},
				},
			},
		},
	}
	return json.Marshal(patch)
}
