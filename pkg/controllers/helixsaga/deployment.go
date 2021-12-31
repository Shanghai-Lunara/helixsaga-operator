package helixsaga

import (
	"fmt"
	helixSagaV1 "github.com/Shanghai-Lunara/helixsaga-operator/pkg/apis/helixsaga/v1"
	k8sCoreV1 "github.com/nevercase/k8s-controller-custom-resource/core/v1"
	appsV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func compareDeployment(original *appsV1.Deployment, updateSpec *helixSagaV1.HelixSagaAppSpec) bool {
	if updateSpec.Replicas != nil && *updateSpec.Replicas != *original.Spec.Replicas {
		return true
	}
	if updateSpec.Image != original.Spec.Template.Spec.Containers[0].Image {
		return true
	}
	// compare Affinity
	// compare Tolerations
	return false
}

func NewDeployment(hs *helixSagaV1.HelixSaga, spec *helixSagaV1.HelixSagaAppSpec) *appsV1.Deployment {
	labels := map[string]string{
		k8sCoreV1.LabelApp:        OperatorKindName,
		k8sCoreV1.LabelController: hs.Name,
		k8sCoreV1.LabelName:       spec.Name,
	}
	dp := &appsV1.Deployment{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      k8sCoreV1.GetStatefulSetName(spec.Name),
			Namespace: hs.Namespace,
			OwnerReferences: []metaV1.OwnerReference{
				*metaV1.NewControllerRef(hs, helixSagaV1.SchemeGroupVersion.WithKind(OperatorKindName)),
			},
			Labels: labels,
		},
		Spec: appsV1.DeploymentSpec{
			Replicas: spec.Replicas,
			Selector: &metaV1.LabelSelector{
				MatchLabels: labels,
			},
			Template: coreV1.PodTemplateSpec{
				ObjectMeta: metaV1.ObjectMeta{
					Labels: labels,
				},
				Spec: coreV1.PodSpec{
					Containers: []coreV1.Container{
						{
							Name:            k8sCoreV1.GetContainerName(spec.Name),
							Image:           spec.Image,
							Ports:           spec.ContainerPorts,
							Env:             ExposePodInformationByEnvs(spec.Env),
							Command:         spec.Command,
							Args:            spec.Args,
							Resources:       spec.Resources,
							ImagePullPolicy: coreV1.PullAlways,
						},
					},
					ImagePullSecrets:   spec.ImagePullSecrets,
					NodeSelector:       spec.NodeSelector,
					ServiceAccountName: spec.ServiceAccountName,
					Affinity:           spec.Affinity,
					Tolerations:        spec.Tolerations,
				},
			},
		},
	}
	// configmap
	dp.Spec.Template.Spec.Volumes = []coreV1.Volume{
		hs.Spec.ConfigMap.Volume,
	}
	dp.Spec.Template.Spec.Containers[0].VolumeMounts = []coreV1.VolumeMount{
		hs.Spec.ConfigMap.VolumeMount,
	}
	if spec.VolumePath != "" {
		t := coreV1.HostPathDirectoryOrCreate
		hostPath := &coreV1.HostPathVolumeSource{
			Type: &t,
			Path: fmt.Sprintf("%s/%s/helixsaga/%s", spec.VolumePath, hs.Namespace, spec.Name),
		}
		dp.Spec.Template.Spec.Volumes = append(dp.Spec.Template.Spec.Volumes,
			coreV1.Volume{
				Name: "task-pv-storage",
				VolumeSource: coreV1.VolumeSource{
					HostPath: hostPath,
				},
			},
		)
		dp.Spec.Template.Spec.Containers[0].VolumeMounts = append(dp.Spec.Template.Spec.Containers[0].VolumeMounts,
			coreV1.VolumeMount{
				MountPath: "/data",
				Name:      "task-pv-storage",
			},
		)
	}
	return dp
}
