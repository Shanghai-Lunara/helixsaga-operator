package serviceloadbalancer

import corev1 "k8s.io/api/core/v1"

func Annotation(svc corev1.ServiceType, isWhiteListOn bool) map[string]string {
	annotations := make(map[string]string, 0)
	switch svc {
	case corev1.ServiceTypeLoadBalancer:
		for k, v := range Get().Annotations {
			annotations[k] = v
		}
		switch isWhiteListOn {
		case true:
			for k, v := range Get().WhiteListOn {
				annotations[k] = v
			}
		case false:
			for k, v := range Get().WhiteListOff {
				annotations[k] = v
			}
		}
	}
	return annotations
}
