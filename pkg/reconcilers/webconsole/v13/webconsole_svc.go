package v13

import (
	fivegv1alpha1 "github.ibm.com/Steve-Glover/5GOperator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func WebconsoleService(cr *fivegv1alpha1.Webconsole) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "webconsole",
			Namespace: cr.Namespace,
			Labels: map[string]string{
				"app": "webconsole",
			},
		},
		Spec: corev1.ServiceSpec{
			Type: "NodePort",
			Ports: []corev1.ServicePort{
				{
					Name:       "webui",
					Port:       5000,
					Protocol:   "TCP",
					NodePort:   30050,
				},
			},
			Selector: map[string]string{
				"app": "webconsole",
			},
			SessionAffinity: "None",
		},
	}
}
