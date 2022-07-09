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
	"github.ibm.com/Steve-Glover/5GOperator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

type OperationRunner interface {
	Start(pod *corev1.Pod, transition *v1alpha1.Transition, componentType string, containerName string) (string, error)
	Stop(pod *corev1.Pod, transition *v1alpha1.Transition, componentType string, containerName string) (string, error)
	RunOperation(pod *corev1.Pod, transition *v1alpha1.Transition, componentType string, containerName string, commands []string) (string, error)
}
