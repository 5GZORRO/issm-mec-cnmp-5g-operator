package v13

import (
	fivegv1alpha1 "github.ibm.com/Steve-Glover/5GOperator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func N3iwfPod(cr *fivegv1alpha1.N3iwf) *corev1.Pod {
	container := n3iwfContainer(cr)

	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			//Name: fmt.Sprintf("%s-%s", cr.Name, "n3iwf"),
			Name:      cr.Name,
			Namespace: cr.Namespace,
			Labels: map[string]string{
				"app": "n3iwf",
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				*container,
			},
			DNSPolicy:       "ClusterFirst",
			RestartPolicy:   "Always",
			SecurityContext: &corev1.PodSecurityContext{},
			Volumes: []corev1.Volume{
				{
					Name: "n3iwfconfig",
					VolumeSource: corev1.VolumeSource{
						ConfigMap: &corev1.ConfigMapVolumeSource{
							LocalObjectReference: corev1.LocalObjectReference{Name: cr.Name},
						},
					},
				},
				{
					Name: "n3iwffree5gcconfig",
					VolumeSource: corev1.VolumeSource{
						ConfigMap: &corev1.ConfigMapVolumeSource{
							LocalObjectReference: corev1.LocalObjectReference{Name: "n3iwf-free5gc"},
						},
					},
				},
			},
		},
	}
}

func n3iwfContainer(cr *fivegv1alpha1.N3iwf) *corev1.Container {
	return &corev1.Container{
		Name:                     "n3iwf",
		Image:                    cr.Spec.Config.ImageUrl,
		ImagePullPolicy:          corev1.PullAlways,
		Env:                      []corev1.EnvVar{},
		Command:                  []string{"/free5gc/n3iwf/n3iwf", "-n3iwfcfg", "/free5gc/config/n3iwf_create.yaml"},
		TerminationMessagePath:   "/dev/termination-log",
		TerminationMessagePolicy: "File",
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      "n3iwfconfig",
				MountPath: "/free5gc/config",
			},
			{
				Name:      "n3iwffree5gcconfig",
				MountPath: "/free5gc/config",
			},
		},
	}
}
