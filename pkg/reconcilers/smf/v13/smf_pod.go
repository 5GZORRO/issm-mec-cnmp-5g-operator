package v13

import (
    "fmt"
	"crypto/sha256"
	"encoding/hex"
	fivegv1alpha1 "github.ibm.com/Steve-Glover/5GOperator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func cmHash(cm *corev1.ConfigMap) string {
	hash := sha256.Sum256([]byte(cm.String()))
	return hex.EncodeToString(hash[:])
}

func SmfPod(cr *fivegv1alpha1.Smf, cm *corev1.ConfigMap) *corev1.Pod {
	container := smfContainer(cr)
	init_container := smfInitContainer(cr)

	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
			Labels: map[string]string{
				"app": "smf",
			},
			Annotations: map[string]string{
				// annotate with cm hash so that if the cm changes, the pod is recreated
				"cmHash": cmHash(cm),
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
					Name: "smfconfig",
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
func smfInitContainer(cr *fivegv1alpha1.Smf) *corev1.Container {
    return &corev1.Container{
        Name: "smf-init",
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

func smfContainer(cr *fivegv1alpha1.Smf) *corev1.Container {
	return &corev1.Container{
		Name:            "smf",
		Image:           cr.Spec.Config.ImageUrl,
		ImagePullPolicy: corev1.PullAlways,
		Env:             []corev1.EnvVar{},
		//Command:         []string{"/bin/sh", "-ec", "while :; do echo '.'; sleep 5 ; done"},
		Command:                  []string{"/free5gc/smf/smf", "--smfcfg", "/free5gc/config/smf.yaml", "--uerouting", "/free5gc/config/uerouting.yaml"},
		TerminationMessagePath:   "/dev/termination-log",
		TerminationMessagePolicy: "File",
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      "smfconfig",
				MountPath: "/free5gc/config",
			},
		},
	}
}
