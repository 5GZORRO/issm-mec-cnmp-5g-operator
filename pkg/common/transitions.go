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
