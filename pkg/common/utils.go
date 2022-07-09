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

package common

import (
	"crypto/sha256"
	"encoding/hex"
	"github.ibm.com/Steve-Glover/5GOperator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

func CMHash(cm *corev1.ConfigMap) string {
	hash := sha256.Sum256([]byte(cm.String()))
	return hex.EncodeToString(hash[:])
}

func GetPodSettings(transition *v1alpha1.Transition) *v1alpha1.PodSettingsSpec {
	requests_memory := transition.Spec.Config.Properties["requests_memory"]
	requests_cpu := transition.Spec.Config.Properties["requests_cpu"]
	limits_memory := transition.Spec.Config.Properties["limits_memory"]
	limits_cpu := transition.Spec.Config.Properties["limits_cpu"]
	return &v1alpha1.PodSettingsSpec{Resources: &v1alpha1.ResourcesSpec{
		Requests: v1alpha1.RequestsSpec{
			Memory: &requests_memory,
			CPU:    &requests_cpu,
		},
		Limits: v1alpha1.LimitsSpec{
			Memory: &limits_memory,
			CPU:    &limits_cpu,
		},
	}}
}
