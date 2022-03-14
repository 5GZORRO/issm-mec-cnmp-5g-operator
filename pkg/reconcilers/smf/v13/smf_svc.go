package v13

import (
	fivegv1alpha1 "github.ibm.com/Steve-Glover/5GOperator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func SmfService(cr *fivegv1alpha1.Smf) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
			Labels: map[string]string{
				// smf name
				"app": cr.Name,
			},
		},
		Spec: corev1.ServiceSpec{
			Type: "NodePort",
			Ports: []corev1.ServicePort{
				{
					Name:       "smf-api",
					Port:       8000,
					Protocol:   "TCP",
					//NodePort:   38000,
				},
				{
					Name:       "smf-ext",
					Port:       8080,
					Protocol:   "TCP",
					//NodePort:   38080,
				},
			},
			Selector: map[string]string{
				// smf name
				"app": cr.Name,
			},
			SessionAffinity: "None",
		},
	}
}
