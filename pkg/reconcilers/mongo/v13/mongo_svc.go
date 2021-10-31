package v13

import (
	fivegv1alpha1 "github.ibm.com/Steve-Glover/5GOperator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func MongoService(cr *fivegv1alpha1.Mongo) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mongodb",
			Namespace: cr.Namespace,
			Labels: map[string]string{
				"app": "mongodb",
			},
		},
		Spec: corev1.ServiceSpec{
			Type: "ClusterIP",
			Ports: []corev1.ServicePort{
				{
					Name:       "mongodb",
					Port:       27017,
					Protocol:   "TCP",
					TargetPort: intstr.FromInt(27017),
				},
			},
			Selector: map[string]string{
				"app": "mongodb",
			},
			SessionAffinity: "None",
		},
	}
}
