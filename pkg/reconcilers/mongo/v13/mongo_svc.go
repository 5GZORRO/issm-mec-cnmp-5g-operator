/*
Copyright 2021 IBM.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
