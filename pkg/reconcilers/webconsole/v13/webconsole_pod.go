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

func WebconsolePod(cr *fivegv1alpha1.Webconsole) *corev1.Pod {
	container := webconsoleContainer(cr)
	init_container := webconsoleInitContainer(cr)

	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			//Name: fmt.Sprintf("%s-%s", cr.Name, "nrf"),
			Name:      cr.Name,
			Namespace: cr.Namespace,
			Labels: map[string]string{
				"app": "webconsole",
			},
			Annotations: map[string]string{
				// Note: important for network names to be consistent with configmap
				"k8s.v1.cni.cncf.io/networks": fmt.Sprintf("[{\"name\" : \"sbi\"}]"),
			},
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
					Name: "webconsole",
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
}
func webconsoleInitContainer(cr *fivegv1alpha1.Webconsole) *corev1.Container {
	return &corev1.Container{
		Name: "webconsole-init",
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

func webconsoleContainer(cr *fivegv1alpha1.Webconsole) *corev1.Container {
	return &corev1.Container{
		Name:                     "webconsole",
		Image:                    cr.Spec.Config.ImageUrl,
		ImagePullPolicy:          corev1.PullAlways,
		Env:                      []corev1.EnvVar{},
		Command:                  []string{"./webui"},
		TerminationMessagePath:   "/dev/termination-log",
		TerminationMessagePolicy: "File",
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      "webconsole",
				MountPath: "/free5gc/config",
			},
		},
	}
}
