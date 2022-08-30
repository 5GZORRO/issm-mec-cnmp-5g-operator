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
	"fmt"
	fivegv1alpha1 "github.com/5GZORRO/issm-mec-cnmp-5g-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func UpfPod(cr *fivegv1alpha1.Upf) *corev1.Pod {
	container := upfContainer(cr)
	init_container := upfInitContainer(cr)

	var anotations = map[string]string{}
	if cr.Spec.Config.DataNetworkName != "" {
		// TODO: write this better
		anotations = map[string]string{
			// Note: important for network names to be consistent with configmap
			"k8s.v1.cni.cncf.io/networks": fmt.Sprintf("[{\"name\" : \"sbi\"}, {\"name\" : \"up\"}, {\"name\" : \"gilan\"}]"),}
	} else {
		anotations = map[string]string{
			"k8s.v1.cni.cncf.io/networks": fmt.Sprintf("[{\"name\" : \"sbi\"}, {\"name\" : \"up\"}]"),}
	}
	pod := corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
			Labels: map[string]string{
				"app": cr.Name,
			},
            Annotations: anotations,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				*container,
			},
			InitContainers: []corev1.Container{
				*init_container,
			},
			DNSPolicy:       "ClusterFirst",
			RestartPolicy:   "Always",
			SecurityContext: &corev1.PodSecurityContext{},
			Volumes: []corev1.Volume{
				{
					Name: "upfconfig",
					VolumeSource: corev1.VolumeSource{
						ConfigMap: &corev1.ConfigMapVolumeSource{
							LocalObjectReference: corev1.LocalObjectReference{Name: cr.Name},
						},
					},
				},
				{
					Name: "podinfo",
					VolumeSource: corev1.VolumeSource{
						DownwardAPI: &corev1.DownwardAPIVolumeSource{
							Items: []corev1.DownwardAPIVolumeFile{
								{
									Path: "annotations",
									FieldRef: &corev1.ObjectFieldSelector{
										FieldPath: "metadata.annotations",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	if cr.Spec.Config.Elicensing != nil && cr.Spec.Config.Elicensing.IsActive{
		licenseSidecarContainer := licensingContainer(cr, false)
		licenseInitContainer := licensingContainer(cr, true)
		pod.Spec.Containers = append(pod.Spec.Containers, *licenseSidecarContainer)
		pod.Spec.InitContainers = append(pod.Spec.InitContainers, *licenseInitContainer)
	}
	return &pod
}

func upfInitContainer(cr *fivegv1alpha1.Upf) *corev1.Container {
	return &corev1.Container{
		Name: "upf-init",
		Image: cr.Spec.Config.ImageUrlInit,
		ImagePullPolicy: corev1.PullAlways,
		Env:             []corev1.EnvVar{
			{
				Name: "CR_NAME",
				Value: cr.Name,
			},
			{
				Name: "DIR",
				Value: "/etc/podinfo",
			},
			{
				Name: "FILE",
				Value: "annotations",
			},
			{
				Name: "CR_KIND",
				Value: cr.Kind,
			},
			{
				Name: "CR_UID",
				Value: string(cr.UID),
			},
		},
		TerminationMessagePath:   "/dev/termination-log",
		TerminationMessagePolicy: "File",
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      "podinfo",
				MountPath: "/etc/podinfo",
			},
		},
	}
}
func upfContainer(cr *fivegv1alpha1.Upf) *corev1.Container {
	allowPrivilegeEscalation := true
	return &corev1.Container{
		Name:            "upf",
		Image:           cr.Spec.Config.ImageUrl,
		ImagePullPolicy: corev1.PullAlways,
		Env:             []corev1.EnvVar{},
		Command:         []string{"/free5gc/free5gc-upfd", "-c", "/free5gc/config/upf.yaml"},
		//Command: []string{"/bin/sh", "-ec", "while :; do echo '.'; sleep 5 ; done"},
		SecurityContext: &corev1.SecurityContext{
			Capabilities: &corev1.Capabilities{
				Add: []corev1.Capability{"NET_ADMIN"},
			},
			AllowPrivilegeEscalation: &allowPrivilegeEscalation,
		},
		TerminationMessagePath:   "/dev/termination-log",
		TerminationMessagePolicy: "File",
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      "upfconfig",
				MountPath: "/free5gc/config",
			},
		},
	}
}

func licensingContainer(cr *fivegv1alpha1.Upf, isInitContainer bool) *corev1.Container {

	container := &corev1.Container{
		Name:            "elma-init-container",
		Image:           cr.Spec.Config.Elicensing.Image,
		ImagePullPolicy: corev1.PullAlways,
		Env: []corev1.EnvVar{
			{
				Name: "INSTANCE_ID",
				ValueFrom: &corev1.EnvVarSource{
					FieldRef: &corev1.ObjectFieldSelector{
						FieldPath: "metadata.uid",
					},
				},
			},
			{
				Name: "DESCRIPTOR_ID",
				Value: cr.Spec.Config.Elicensing.DescriptorId,
			},
			{
				Name:  "ELMA_IP",
				Value: cr.Spec.Config.Elicensing.ElmaIp,
			},
			{
				Name:  "PRODUCT_OFFERING",
				Value: cr.Spec.Config.Elicensing.ProductOfferingId,
			},
		},
		Command: []string{"/bin/bash", "-c", "/bin/bash inithook.sh"},
	}
	if !isInitContainer {
		container.Lifecycle = &corev1.Lifecycle{
			PreStop: &corev1.Handler{
				Exec: &corev1.ExecAction{
					Command: []string{"/bin/bash", "-c", "curl 'http://localhost:8000/SendEndHook'"},
				},
			},
		}
		container.Name = "elma-sidecar-container"
		container.Command = []string{"/bin/bash", "-c", "/bin/bash start.sh"}
	}
	return container
}
