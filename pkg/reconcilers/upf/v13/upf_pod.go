package v13

import (
	"fmt"
	fivegv1alpha1 "github.ibm.com/Steve-Glover/5GOperator/api/v1alpha1"
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
	return &corev1.Pod{
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
		Command:         []string{"/free5gc/free5gc-upfd/free5gc-upfd", "-f", "/free5gc/config/upf.yaml"},
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
