package helixsaga

import (
	"fmt"

	appsV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"

	helixSagaV1 "github.com/Shanghai-Lunara/helixsaga-operator/pkg/apis/helixsaga/v1"
	helixSagaClientSet "github.com/Shanghai-Lunara/helixsaga-operator/pkg/generated/helixsaga/clientset/versioned"
	k8sCoreV1 "github.com/nevercase/k8s-controller-custom-resource/core/v1"
)

func NewNginxPhpFpm(ks k8sCoreV1.KubernetesResource, client helixSagaClientSet.Interface, hs *helixSagaV1.HelixSaga, spec helixSagaV1.NginxPhpFpmSpec) error {
	ss, err := ks.StatefulSet().Get(hs.Namespace, spec.Name)
	if err != nil {
		klog.Info("statefulSet err:", err)
		if !errors.IsNotFound(err) {
			return err
		}
		klog.Info("new statefulSet")
		if ss, err = ks.StatefulSet().Create(hs.Namespace, spec.Name, NewStatefulSet(hs, spec)); err != nil {
			return err
		}
	}
	klog.Info("rds:", *spec.Replicas)
	klog.Info("statefulSet:", *ss.Spec.Replicas)
	if spec.Replicas != nil && *spec.Replicas != *ss.Spec.Replicas {
		if ss, err = ks.StatefulSet().Update(hs.Namespace, NewStatefulSet(hs, spec)); err != nil {
			klog.Info(err)
			return err
		}
	}
	if err = updateFooStatus(hs, client, ss); err != nil {
		return err
	}
	return nil
}

func NewStatefulSet(hs *helixSagaV1.HelixSaga, spec helixSagaV1.NginxPhpFpmSpec) *appsV1.StatefulSet {
	labels := map[string]string{
		"app":        operatorKindName,
		"controller": hs.Name,
	}
	t := coreV1.HostPathDirectoryOrCreate
	hostPath := &coreV1.HostPathVolumeSource{
		Type: &t,
		Path: fmt.Sprintf("/mnt/ssd1/mysql/%s", spec.Name),
	}
	objectName := fmt.Sprintf(k8sCoreV1.StatefulSetNameTemplate, spec.Name)
	containerName := fmt.Sprintf(k8sCoreV1.ContainerNameTemplate, spec.Name)
	standard := &appsV1.StatefulSet{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      objectName,
			Namespace: hs.Namespace,
			OwnerReferences: []metaV1.OwnerReference{
				*metaV1.NewControllerRef(hs, helixSagaV1.SchemeGroupVersion.WithKind(operatorKindName)),
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
						{
							Name: "task-pv-storage",
							VolumeSource: coreV1.VolumeSource{
								HostPath: hostPath,
							},
						},
						//{
						//	Name: "test-configmap",
						//	VolumeSource: coreV1.VolumeSource{
						//		ConfigMap: &coreV1.ConfigMapVolumeSource{
						//			Items: []coreV1.KeyToPath{
						//				{
						//					Key:  "mysql.php",
						//					Path: "mysql.php",
						//				},
						//				{
						//					Key:  "redis.php",
						//					Path: "redis.php",
						//				},
						//				{
						//					Key:  "version.php",
						//					Path: "version.php",
						//				},
						//				{
						//					Key:  "gs.php",
						//					Path: "gs.php",
						//				},
						//			},
						//		},
						//	},
						//},
					},
					Containers: []coreV1.Container{
						{
							Name:  containerName,
							Image: spec.Image,
							Ports: []coreV1.ContainerPort{
								{
									ContainerPort: NginxPhpFpmDefaultPort,
								},
							},
							Env: spec.Env,
							VolumeMounts: []coreV1.VolumeMount{
								{
									MountPath: "/data",
									Name:      "task-pv-storage",
								},
								//{
								//	MountPath: "/var/www/app/conf",
								//	Name:      "test-configmap",
								//},
							},
						},
					},
					ImagePullSecrets: spec.ImagePullSecrets,
				},
			},
		},
	}
	return standard
}
