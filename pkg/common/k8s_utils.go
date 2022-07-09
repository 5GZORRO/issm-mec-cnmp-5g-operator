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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/5GZORRO/issm-mec-cnmp-5g-operator/api/v1alpha1"
	"io"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	routev1 "github.com/openshift/api/route/v1"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	ingressv1b1 "k8s.io/api/extensions/v1beta1"
	networkingv1 "k8s.io/api/networking/v1"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	networking "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var log = logf.Log.WithName("common")

/*
 * Represents a CP4NA transition CR execution against one of the 5G CRs e.g. AMF, AUSF, etc
 */
type TransitionFunc func(transition *v1alpha1.Transition) (string, error)

// K8sUtils k8s utilities for interacting with k8s api's
type K8sUtils struct {
	RestConfig *rest.Config
	ClientSet  *kubernetes.Clientset
	Client     client.Client
	Scheme     *runtime.Scheme
}

func NewK8sUtils(restConfig *rest.Config, clientSet *kubernetes.Clientset, client client.Client,
	scheme *runtime.Scheme) *K8sUtils {
	return &K8sUtils{
		RestConfig: restConfig,
		ClientSet:  clientSet,
		Client:     client,
		Scheme:     scheme,
	}
}

// ReconcileConfigMap reconcile a ConfigMap with provided configuration
func (k8s K8sUtils) ReconcileConfigMap(cr metav1.Object, configmap *corev1.ConfigMap) (controllerutil.OperationResult, error) {
	log.Info("ReconcileConfigMap ", configmap.Name, configmap.Namespace)
	if err := setOwner(k8s.Scheme, cr, configmap); err != nil {
		return "", err
	}

	opResult := controllerutil.OperationResultNone

	// Check if this CongifMap already exists
	found := &corev1.ConfigMap{}
	err := k8s.Client.Get(context.TODO(), types.NamespacedName{Name: configmap.Name, Namespace: configmap.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating a new ConfigMap 1", "configMap.Namespace", configmap.Namespace)
		log.Info("Creating a new ConfigMap 2", "configMap.Name", configmap.Name)
		err = k8s.Client.Create(context.TODO(), configmap)
		if err != nil {
			return "", err
		}
		opResult = controllerutil.OperationResultCreated
	} else if err != nil {
		return "", err
	} else if !reflect.DeepEqual(configmap.Data, found.Data) || !reflect.DeepEqual(configmap.BinaryData, found.BinaryData) {
		log.Info("Updating ConfigMap 1", "configMap.Namespace", configmap.Namespace)
		log.Info("Updating ConfigMap 2", "configMap.Name", configmap.Name)
		configmap.ObjectMeta = found.ObjectMeta
		err = k8s.Client.Update(context.TODO(), configmap)
		if err != nil {
			return "", err
		} else {
			// restart the POD
			pod, err := k8s.GetPod(types.NamespacedName{
				Namespace: cr.GetNamespace(),
				Name:      cr.GetName(),
				})
			if err != nil {
				log.Info("Restart for ConfigMap 1: unable to get POD for this cm")
			} else {
				// TODO: hard code
				//if pod.Name == "smf-sample" {
				log.Info("Restart POD")
				// TODO: define criteria for pods that want to get restarted on cm changes
				// at most will fail if not SMF
				_, _ = k8s.ExecCommand(pod, "smf", "kill -TERM 1")
				//}
			}
		}
		opResult = controllerutil.OperationResultUpdated
	}
	return opResult, nil
}

// ReconcileService reconcile a Service with provided configuration
func (k8s K8sUtils) ReconcileService(cr metav1.Object, service *corev1.Service) (controllerutil.OperationResult, error) {
	log.Info("ReconcileService ", service.Name, service.Namespace)
	if err := setOwner(k8s.Scheme, cr, service); err != nil {
		return "", err
	}

	opResult := controllerutil.OperationResultNone

	// Check if this Service already exists
	found := &corev1.Service{}
	err := k8s.Client.Get(context.TODO(), types.NamespacedName{Name: service.Name, Namespace: service.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating a new Service", "Service", service.String())
		err = k8s.Client.Create(context.TODO(), service)
		if err != nil {
			return "", err
		}
		opResult = controllerutil.OperationResultCreated
	} else if err != nil {
		return "", err
	} else if service.Spec.Type != "NodePort" { // TODO not updating because it continnually recycles ports - how to avoid?
		// Handle update scenario
		service.Spec.ClusterIP = found.Spec.ClusterIP
		if !reflect.DeepEqual(service.Spec, found.Spec) {
			service.ObjectMeta = found.ObjectMeta
			log.Info("Updating Service", "Service", service.String())
			err = k8s.Client.Update(context.TODO(), service)
			if err != nil {
				return "", err
			}
			opResult = controllerutil.OperationResultUpdated
		}
	}

	return opResult, nil
}

// ReconcileStatefulSet reconcile a StatefulSet with provided configuration
func (k8s K8sUtils) ReconcileStatefulSet(cr metav1.Object, statefulSet *appsv1.StatefulSet) (controllerutil.OperationResult, error) {
	log.Info("ReconcileStatefulSet ", statefulSet.Name, statefulSet.Namespace)
	if err := setOwner(k8s.Scheme, cr, statefulSet); err != nil {
		return "", err
	}

	ss := &appsv1.StatefulSet{ObjectMeta: statefulSet.ObjectMeta}
	result, err := controllerutil.CreateOrUpdate(context.TODO(), k8s.Client, ss, func() error {
		ss.Spec = statefulSet.Spec
		return nil
	})

	log.Info("StatefulSet", "action", result)

	if err != nil {
		return "", err
	}

	return result, nil
}

// ReconcileDeployment reconcile a Deployment with provided configuration
func (k8s K8sUtils) ReconcileDeployment(cr metav1.Object, deployment *appsv1.Deployment) (controllerutil.OperationResult, error) {
	log.Info("ReconcileDeployment ", deployment.Name, deployment.Namespace)
	if err := setOwner(k8s.Scheme, cr, deployment); err != nil {
		return "", err
	}

	dp := &appsv1.Deployment{ObjectMeta: deployment.ObjectMeta}
	result, err := controllerutil.CreateOrUpdate(context.TODO(), k8s.Client, dp, func() error {
		dp.Spec = deployment.Spec
		return nil
	})

	log.Info("Deployment", "action", result)

	if err != nil {
		return "", err
	}

	return result, nil
}

// ReconcilePod reconcile a Pod with provided configuration
func (k8s K8sUtils) ReconcilePod(cr metav1.Object, pod *corev1.Pod) (controllerutil.OperationResult, error) {
	log.Info("ReconcilePod ", pod.Name, pod.Namespace)
	if err := setOwner(k8s.Scheme, cr, pod); err != nil {
		return "", err
	}

	dp := &corev1.Pod{ObjectMeta: pod.ObjectMeta}

	log.Info("ReconcilePod2 ", "pod.ObjectMeta", pod.ObjectMeta, "dp", dp)
	// dp is the actual retrieved pod
	result, err := controllerutil.CreateOrUpdate(context.TODO(), k8s.Client, dp, func() error {
		b, err := json.Marshal(dp)
		var asStr string
		if err == nil {
			asStr = string(b)
		}

		log.Info("ReconcilePod before update", "dp", asStr)

		// Parts of Pod spec are immutable so we set this value only if
		// a new object is going to be created
		if dp.ObjectMeta.CreationTimestamp.IsZero() {
			dp.Spec = pod.Spec
		} else {
			for _, podContainer := range pod.Spec.Containers {
				found := false
				for idx, dpContainer := range dp.Spec.Containers {
					if podContainer.Name == dpContainer.Name {
						log.Info("ReconcilePodx1", "dpContainer.Name", dpContainer.Name)
						dp.Spec.Containers[idx].Image = podContainer.Image
						found = true
						break
					}
				}
				if !found {
					dp.Spec.Containers = append(dp.Spec.Containers, podContainer)
				}
			}
			for _, podContainer := range pod.Spec.InitContainers {
				found := false
				for _, dpContainer := range dp.Spec.InitContainers {
					if podContainer.Name == dpContainer.Name {
						log.Info("ReconcilePodx2", "dpContainer.Name", dpContainer.Name)
						dpContainer.Image = podContainer.Image
						found = true
						break
					}
				}
				if !found {
					dp.Spec.InitContainers = append(dp.Spec.InitContainers, podContainer)
				}
			}
			for _, podToleration := range pod.Spec.Tolerations {
				found := false
				for _, dpToleration := range dp.Spec.Tolerations {
					if podToleration.Key == dpToleration.Key {
						found = true
						break
					}
				}
				if !found {
					log.Info("ReconcilePodx3", "dpContainer.Name", podToleration.Key)
					dp.Spec.Tolerations = append(dp.Spec.Tolerations, podToleration)
				}
			}
			log.Info("ReconcilePodx4", "db.cmHashpod", dp.Annotations["cmHash"])
			log.Info("ReconcilePodx4", "pod.cmHashpod", pod.Annotations["cmHash"])
//			if dp.Annotations["cmHash"] != pod.Annotations["cmHash"] {
//				// TODO: hard code
//				if pod.Name == "smf-sample" {
//					log.Info("Restart POD")
//					_, _ = k8s.ExecCommand(pod, "smf", "kill -TERM 1")
//				}
//				// update dp with cmHash with most uptodate one..
//				//dp.Annotations["cmHash"] = pod.Annotations["cmHash"]
//			}
		}

		b, err = json.Marshal(dp)
		if err == nil {
			asStr = string(b)
		}

		log.Info("ReconcilePod after update", "dp", asStr)

		return nil
	})

	log.Info("Pod", "action", result, "err", err)

	if err != nil {
		return "", err
	}

	return result, nil
}

// ReconcileDaemonSet reconcile a DaemonSet with provided configuration
func (k8s K8sUtils) ReconcileDaemonSet(cr metav1.Object, daemonset *appsv1.DaemonSet) (controllerutil.OperationResult, error) {
	log.Info("ReconcileDaemonSet ", daemonset.Name, daemonset.Namespace)
	if err := setOwner(k8s.Scheme, cr, daemonset); err != nil {
		return "", err
	}

	ds := &appsv1.DaemonSet{ObjectMeta: daemonset.ObjectMeta}
	result, err := controllerutil.CreateOrUpdate(context.TODO(), k8s.Client, ds, func() error {
		ds.Spec = daemonset.Spec
		return nil
	})

	log.Info("DaemonSet", "action", result)

	if err != nil {
		return "", err
	}

	return result, nil
}

// ReconcileJob reconcile a Job with provided configuration
func (k8s K8sUtils) ReconcileJob(cr metav1.Object, job *batchv1.Job) (controllerutil.OperationResult, error) {
	log.Info("ReconcileJob ", job.Name, job.Namespace)
	if err := setOwner(k8s.Scheme, cr, job); err != nil {
		return "", err
	}

	opResult := controllerutil.OperationResultNone

	// Check if this Job already exists
	found := &batchv1.Job{}
	err := k8s.Client.Get(context.TODO(), types.NamespacedName{Name: job.Name, Namespace: job.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating a new Job", "Job.Namespace", job.Namespace, "Job.Name", job.Name)
		err = k8s.Client.Create(context.TODO(), job)
		if err != nil {
			return "", err
		}
		opResult = controllerutil.OperationResultCreated
	} else if err != nil {
		return "", err
	}

	return opResult, nil
}

// ReconcileRoute reconcile a Route with provided configuration
func (k8s K8sUtils) ReconcileRoute(cr metav1.Object, route *routev1.Route) (controllerutil.OperationResult, error) {
	log.Info("ReconcileRoute ", route.Name, route.Namespace)
	if err := setOwner(k8s.Scheme, cr, route); err != nil {
		return "", err
	}

	opResult := controllerutil.OperationResultNone

	// Check if this Route already exists
	found := &routev1.Route{}
	err := k8s.Client.Get(context.TODO(), types.NamespacedName{Name: route.Name, Namespace: route.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating a new Route", "Route.Namespace", route.Namespace, "Route.Name", route.Name)
		err = k8s.Client.Create(context.TODO(), route)
		if err != nil {
			return "", err
		}
		opResult = controllerutil.OperationResultCreated
	} else if err != nil {
		return "", err
	}

	route.Spec.Host = found.Spec.Host
	if !reflect.DeepEqual(route.Spec, found.Spec) {
		log.Info("Updating Route", "Route.Namespace", route.Namespace, "Route.Name", route.Name)
		route.ObjectMeta = found.ObjectMeta
		err = k8s.Client.Update(context.TODO(), route)
		if err != nil {
			return "", err
		}
		opResult = controllerutil.OperationResultUpdated
	}

	return opResult, nil
}

// ReconcileIngress reconcile a Ingress with provided configuration
func (k8s K8sUtils) ReconcileIngress(cr metav1.Object, ingress *ingressv1b1.Ingress) (controllerutil.OperationResult, error) {
	log.Info("ReconcileIngress ", ingress.Name, ingress.Namespace)
	if err := setOwner(k8s.Scheme, cr, ingress); err != nil {
		return "", err
	}

	opResult := controllerutil.OperationResultNone

	// Check if this Ingress already exists
	found := &ingressv1b1.Ingress{}
	err := k8s.Client.Get(context.TODO(), types.NamespacedName{Name: ingress.Name, Namespace: ingress.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating a new Ingress", "Ingress.Namespace", ingress.Namespace, "Ingress.Name", ingress.Name)
		err = k8s.Client.Create(context.TODO(), ingress)
		if err != nil {
			return "", err
		}
		opResult = controllerutil.OperationResultCreated
	} else if err != nil {
		return "", err
	}

	if !reflect.DeepEqual(ingress.Spec, found.Spec) {
		log.Info("Updating Ingress", "Ingress.Namespace", ingress.Namespace, "Ingress.Name", ingress.Name)
		ingress.ObjectMeta = found.ObjectMeta
		err = k8s.Client.Update(context.TODO(), ingress)
		if err != nil {
			return "", err
		}
		opResult = controllerutil.OperationResultUpdated
	}

	return opResult, nil
}

// ReconcilePodDisruptionBudget reconcile a PodDisruptionBudget with provided configuration
func (k8s K8sUtils) ReconcilePodDisruptionBudget(cr metav1.Object, podDisruptionBudget *policyv1beta1.PodDisruptionBudget) (controllerutil.OperationResult, error) {
	log.Info("ReconcilePodDisruptionBudget ", podDisruptionBudget.Name, podDisruptionBudget.Namespace)
	if err := setOwner(k8s.Scheme, cr, podDisruptionBudget); err != nil {
		return "", err
	}

	opResult := controllerutil.OperationResultNone

	// Check if this PodDisruptionBudget already exists
	found := &policyv1beta1.PodDisruptionBudget{}
	err := k8s.Client.Get(context.TODO(), types.NamespacedName{Name: podDisruptionBudget.Name, Namespace: podDisruptionBudget.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating a new PodDisruptionBudget", "PodDisruptionBudget.Namespace", podDisruptionBudget.Namespace, "PodDisruptionBudget.Name", podDisruptionBudget.Name)
		err = k8s.Client.Create(context.TODO(), podDisruptionBudget)
		if err != nil {
			return "", err
		}
		opResult = controllerutil.OperationResultCreated
	} else if err != nil {
		return "", err
	} else if !reflect.DeepEqual(podDisruptionBudget.Spec, found.Spec) {
		log.Info("Updating PodDisruptionBudget", "PodDisruptionBudget.Namespace", podDisruptionBudget.Namespace, "PodDisruptionBudget.Name", podDisruptionBudget.Name)
		podDisruptionBudget.ObjectMeta = found.ObjectMeta
		err = k8s.Client.Update(context.TODO(), podDisruptionBudget)
		if err != nil {
			return "", err
		}
		opResult = controllerutil.OperationResultUpdated
	}

	return opResult, nil
}

func (k8s K8sUtils) ReconcilePVC(cr metav1.Object, pvc *corev1.PersistentVolumeClaim) (controllerutil.OperationResult, error) {
	log.Info("ReconcilePVC ", pvc.Name, pvc.Namespace)
	if err := setOwner(k8s.Scheme, cr, pvc); err != nil {
		return "", err
	}

	opResult := controllerutil.OperationResultNone

	// Check if this PVC already exists
	found := &corev1.PersistentVolumeClaim{}
	err := k8s.Client.Get(context.TODO(), types.NamespacedName{Name: pvc.Name, Namespace: pvc.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating a new PVC", "PVC.Namespace", pvc.Namespace, "PVC.Name", pvc.Name)
		err = k8s.Client.Create(context.TODO(), pvc)
		if err != nil {
			return "", err
		}
		opResult = controllerutil.OperationResultCreated
	} else if err != nil {
		return "", err
	}

	// Note - Can't update a PVC
	return opResult, nil
}

// ReconcileRole reconcile a Role with provided configuration
func (k8s K8sUtils) ReconcileRole(cr metav1.Object, role *rbacv1.Role) (controllerutil.OperationResult, error) {
	log.Info("ReconcileRole ", role.Name, role.Namespace)
	if err := setOwner(k8s.Scheme, cr, role); err != nil {
		return "", err
	}

	opResult := controllerutil.OperationResultNone

	// Check if this Role already exists
	found := &rbacv1.Role{}
	err := k8s.Client.Get(context.TODO(), types.NamespacedName{Name: role.Name, Namespace: role.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating a new Role", "Role.Namespace", role.Namespace, "Role.Name", role.Name)
		err = k8s.Client.Create(context.TODO(), role)
		if err != nil {
			return "", err
		}
		opResult = controllerutil.OperationResultCreated
	} else if err != nil {
		return "", err
	} else if !reflect.DeepEqual(role.Rules, found.Rules) {
		log.Info("Updating Role", "Role.Namespace", role.Namespace, "Role.Name", role.Name)
		role.ObjectMeta = found.ObjectMeta
		err = k8s.Client.Update(context.TODO(), role)
		if err != nil {
			return "", err
		}
		opResult = controllerutil.OperationResultUpdated
	}

	return opResult, nil
}

// ReconcileRoleBinding reconcile a RoleBinding with provided configuration
func (k8s K8sUtils) ReconcileRoleBinding(cr metav1.Object, rolebinding *rbacv1.RoleBinding) (controllerutil.OperationResult, error) {
	log.Info("ReconcileRoleBinding ", rolebinding.Name, rolebinding.Namespace)
	if err := setOwner(k8s.Scheme, cr, rolebinding); err != nil {
		return "", err
	}

	opResult := controllerutil.OperationResultNone

	// Check if this RoleBinding already exists
	found := &rbacv1.RoleBinding{}
	err := k8s.Client.Get(context.TODO(), types.NamespacedName{Name: rolebinding.Name, Namespace: rolebinding.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating a new RoleBinding", "RoleBinding.Namespace", rolebinding.Namespace, "RoleBinding.Name", rolebinding.Name)
		err = k8s.Client.Create(context.TODO(), rolebinding)
		if err != nil {
			return "", err
		}
		opResult = controllerutil.OperationResultCreated
	} else if err != nil {
		return "", err
	} else if !reflect.DeepEqual(rolebinding.Subjects, found.Subjects) || !reflect.DeepEqual(rolebinding.RoleRef, found.RoleRef) {
		log.Info("Updating RoleBinding", "RoleBinding.Namespace", rolebinding.Namespace, "RoleBinding.Name", rolebinding.Name)
		rolebinding.ObjectMeta = found.ObjectMeta
		err = k8s.Client.Update(context.TODO(), rolebinding)
		if err != nil {
			return "", err
		}
		opResult = controllerutil.OperationResultUpdated
	}

	return opResult, nil
}

// ReconcileClusterRole reconcile a ClusterRole with provided configuration
func (k8s K8sUtils) ReconcileClusterRole(cr metav1.Object, clusterrole *rbacv1.ClusterRole) (controllerutil.OperationResult, error) {
	log.Info("ReconcileRole ", " ClusterRole ", clusterrole.Name)

	opResult := controllerutil.OperationResultNone

	// Check if this ClusterRole already exists
	found := &rbacv1.ClusterRole{}
	err := k8s.Client.Get(context.TODO(), types.NamespacedName{Name: clusterrole.Name}, found)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating a new ClusterRole", "ClusterRole.Name", clusterrole.Name)
		err = k8s.Client.Create(context.TODO(), clusterrole)
		if err != nil {
			return "", err
		}
		opResult = controllerutil.OperationResultCreated
	} else if err != nil {
		return "", err
	} else if !reflect.DeepEqual(clusterrole.Rules, found.Rules) {
		log.Info("Updating ClusterRole", "Role.Name", clusterrole.Name)
		clusterrole.ObjectMeta = found.ObjectMeta
		err = k8s.Client.Update(context.TODO(), clusterrole)
		if err != nil {
			return "", err
		}
		opResult = controllerutil.OperationResultUpdated
	}

	return opResult, nil
}

// ReconcileClusterRoleBinding reconcile a ClusterRoleBinding with provided configuration
func (k8s K8sUtils) ReconcileClusterRoleBinding(cr metav1.Object, clusterrolebinding *rbacv1.ClusterRoleBinding) (controllerutil.OperationResult, error) {
	log.Info("ReconcileRoleBinding ", "ClusterRoleBinding ", clusterrolebinding.Name)

	opResult := controllerutil.OperationResultNone

	// Check if this ClusterRoleBinding already exists
	found := &rbacv1.ClusterRoleBinding{}
	err := k8s.Client.Get(context.TODO(), types.NamespacedName{Name: clusterrolebinding.Name}, found)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating a new ClusterRoleBinding", "ClusterRoleBinding.Name", clusterrolebinding.Name)
		err = k8s.Client.Create(context.TODO(), clusterrolebinding)
		if err != nil {
			return "", err
		}
		opResult = controllerutil.OperationResultCreated
	} else if err != nil {
		return "", err
	} else if !reflect.DeepEqual(clusterrolebinding.Subjects, found.Subjects) || !reflect.DeepEqual(clusterrolebinding.RoleRef, found.RoleRef) {
		log.Info("Updating ClusterRoleBinding", "ClusterRoleBinding.Name", clusterrolebinding.Name)
		clusterrolebinding.ObjectMeta = found.ObjectMeta
		err = k8s.Client.Update(context.TODO(), clusterrolebinding)
		if err != nil {
			return "", err
		}
		opResult = controllerutil.OperationResultUpdated
	}

	return opResult, nil
}

// ReconcileServiceAccount reconcile a Service Account with provided configuration
func (k8s K8sUtils) ReconcileServiceAccount(cr metav1.Object, sa *corev1.ServiceAccount) (controllerutil.OperationResult, error) {
	log.Info("ReconcileServiceAccount ", sa.Name, sa.Namespace)
	if err := setOwner(k8s.Scheme, cr, sa); err != nil {
		return "", err
	}

	opResult := controllerutil.OperationResultNone

	// Check if this SA already exists
	found := &corev1.ServiceAccount{}
	err := k8s.Client.Get(context.TODO(), types.NamespacedName{Name: sa.Name, Namespace: sa.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating a new Service Account", "ServiceAccount.Namespace", sa.Namespace, "ServiceAccount.Name", sa.Name)
		err = k8s.Client.Create(context.TODO(), sa)
		if err != nil {
			return "", err
		}
		opResult = controllerutil.OperationResultCreated
	} else if err != nil {
		return "", err
	}

	// TODO ServiceAccount update not implemented - review
	return opResult, nil
}

func (k8s K8sUtils) ReconcileNetworkPolicy(cr metav1.Object, np *networkingv1.NetworkPolicy) (controllerutil.OperationResult, error) {
	log.Info("ReconcileNetworkPolicy ", np.Name, np.Namespace)
	if err := setOwner(k8s.Scheme, cr, np); err != nil {
		return "", err
	}

	opResult := controllerutil.OperationResultNone

	// Check if this NetworkPolicy already exists
	found := &networkingv1.NetworkPolicy{}
	err := k8s.Client.Get(context.TODO(), types.NamespacedName{Name: np.Name, Namespace: np.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating a new NetworkPolicy", "NetworkPolicy.Namespace", np.Namespace, "NetworkPolicy.Name", np.Name)
		err = k8s.Client.Create(context.TODO(), np)
		if err != nil {
			return "", err
		}
		opResult = controllerutil.OperationResultCreated
	} else if err != nil {
		return "", err
	} else if !reflect.DeepEqual(np.Spec, found.Spec) {
		log.Info("Updating NetworkPolicy", "NetworkPolicy.Namespace", np.Namespace, "NetworkPolicy.Name", np.Name)
		np.ObjectMeta = found.ObjectMeta
		err = k8s.Client.Update(context.TODO(), np)
		if err != nil {
			return "", err
		}
		opResult = controllerutil.OperationResultUpdated
	}

	return opResult, nil
}

// Delete deletes the specified object
func (k8s K8sUtils) Delete(obj *unstructured.Unstructured) {
	policy := metav1.DeletePropagationBackground
	options := &client.DeleteAllOfOptions{
		DeleteOptions: client.DeleteOptions{
			PropagationPolicy: &policy,
		},
	}
	err := k8s.Client.Delete(context.TODO(), obj, options)
	if err != nil {
		log.Error(err, "Error deleting unstuctured object")
	}
}

// DeleteByNameSpaceLabels deletes an object in the specified namespace by labels
func (k8s K8sUtils) DeleteByNameSpaceLabels(obj *unstructured.Unstructured, namespace string, labels map[string]string) {
	policy := metav1.DeletePropagationBackground
	options := &client.DeleteAllOfOptions{
		DeleteOptions: client.DeleteOptions{
			PropagationPolicy: &policy,
		},
	}
	err := k8s.Client.DeleteAllOf(context.TODO(), obj, client.InNamespace(namespace), client.MatchingLabels(labels), options)
	if err != nil {
		log.Error(err, "Error in DeleteByNameSpaceLabels")
	}
}

// GetNetworkPolicy get a NetworkPolicy object
func (k8s K8sUtils) GetNetworkPolicy(namespacedName types.NamespacedName) (*networking.NetworkPolicy, error) {
	found := networking.NetworkPolicy{}
	err := k8s.Client.Get(context.TODO(), namespacedName, &found)
	if err != nil {
		return nil, err
	}
	return &found, nil
}

// GetJob get a Job object
func (k8s K8sUtils) GetJob(namespacedName types.NamespacedName) (*batchv1.Job, error) {
	found := batchv1.Job{}
	err := k8s.Client.Get(context.TODO(), namespacedName, &found)
	if err != nil {
		return nil, err
	}
	return &found, nil
}

// GetDeployment get a Deployment object
func (k8s K8sUtils) GetDeployment(namespacedName types.NamespacedName) (*appsv1.Deployment, error) {
	found := appsv1.Deployment{}
	err := k8s.Client.Get(context.TODO(), namespacedName, &found)
	if err != nil {
		return nil, err
	}
	return &found, nil
}

// GetPod get a Pod object
func (k8s K8sUtils) GetPod(namespacedName types.NamespacedName) (*corev1.Pod, error) {
	found := corev1.Pod{}
	err := k8s.Client.Get(context.TODO(), namespacedName, &found)
	if err != nil {
		return nil, err
	}
	return &found, nil
}

func (k8s K8sUtils) UpdateCR(cr client.Object) error {
	err := k8s.Client.Update(context.TODO(), cr)
	if err != nil {
		return err
	}
	return nil
}

// GetAmf get a Amf object
func (k8s K8sUtils) GetAmf(namespacedName types.NamespacedName) (*v1alpha1.Amf, error) {
	found := v1alpha1.Amf{}
	err := k8s.Client.Get(context.TODO(), namespacedName, &found)
	if err != nil {
		return nil, err
	}
	return &found, nil
}

// GetMongo get a Mongo object
func (k8s K8sUtils) GetMongo(namespacedName types.NamespacedName) (*v1alpha1.Mongo, error) {
	found := v1alpha1.Mongo{}
	err := k8s.Client.Get(context.TODO(), namespacedName, &found)
	if err != nil {
		return nil, err
	}
	return &found, nil
}

// GetUdm get a Udm object
func (k8s K8sUtils) GetUdm(namespacedName types.NamespacedName) (*v1alpha1.Udm, error) {
	found := v1alpha1.Udm{}
	err := k8s.Client.Get(context.TODO(), namespacedName, &found)
	if err != nil {
		return nil, err
	}
	return &found, nil
}

// GetUdr get a Udr object
func (k8s K8sUtils) GetUdr(namespacedName types.NamespacedName) (*v1alpha1.Udr, error) {
	found := v1alpha1.Udr{}
	err := k8s.Client.Get(context.TODO(), namespacedName, &found)
	if err != nil {
		return nil, err
	}
	return &found, nil
}

// GetUpf get a Upf object
func (k8s K8sUtils) GetUpf(namespacedName types.NamespacedName) (*v1alpha1.Upf, error) {
	found := v1alpha1.Upf{}
	err := k8s.Client.Get(context.TODO(), namespacedName, &found)
	if err != nil {
		return nil, err
	}
	return &found, nil
}

// GetWebconsole get a Webconsole object
func (k8s K8sUtils) GetWebconsole(namespacedName types.NamespacedName) (*v1alpha1.Webconsole, error) {
	found := v1alpha1.Webconsole{}
	err := k8s.Client.Get(context.TODO(), namespacedName, &found)
	if err != nil {
		return nil, err
	}
	return &found, nil
}

// GetSmf get a Smf object
func (k8s K8sUtils) GetSmf(namespacedName types.NamespacedName) (*v1alpha1.Smf, error) {
	found := v1alpha1.Smf{}
	err := k8s.Client.Get(context.TODO(), namespacedName, &found)
	if err != nil {
		return nil, err
	}
	return &found, nil
}

// GetN3iwf get a N3iwf object
func (k8s K8sUtils) GetN3iwf(namespacedName types.NamespacedName) (*v1alpha1.N3iwf, error) {
	found := v1alpha1.N3iwf{}
	err := k8s.Client.Get(context.TODO(), namespacedName, &found)
	if err != nil {
		return nil, err
	}
	return &found, nil
}

// GetSmf get a Subscriber object
func (k8s K8sUtils) GetSubscriber(namespacedName types.NamespacedName) (*v1alpha1.Subscriber, error) {
	found := v1alpha1.Subscriber{}
	err := k8s.Client.Get(context.TODO(), namespacedName, &found)
	if err != nil {
		return nil, err
	}
	return &found, nil
}

// GetNrf get a Nrf object
func (k8s K8sUtils) GetNrf(namespacedName types.NamespacedName) (*v1alpha1.Nrf, error) {
	found := v1alpha1.Nrf{}
	err := k8s.Client.Get(context.TODO(), namespacedName, &found)
	if err != nil {
		return nil, err
	}
	return &found, nil
}

// GetNssf get a Nssf object
func (k8s K8sUtils) GetNssf(namespacedName types.NamespacedName) (*v1alpha1.Nssf, error) {
	found := v1alpha1.Nssf{}
	err := k8s.Client.Get(context.TODO(), namespacedName, &found)
	if err != nil {
		return nil, err
	}
	return &found, nil
}

// GetAusf get a Ausf object
func (k8s K8sUtils) GetAusf(namespacedName types.NamespacedName) (*v1alpha1.Ausf, error) {
	found := v1alpha1.Ausf{}
	err := k8s.Client.Get(context.TODO(), namespacedName, &found)
	if err != nil {
		return nil, err
	}
	return &found, nil
}

// GetPcf get a Pcf object
func (k8s K8sUtils) GetPcf(namespacedName types.NamespacedName) (*v1alpha1.Pcf, error) {
	found := v1alpha1.Pcf{}
	err := k8s.Client.Get(context.TODO(), namespacedName, &found)
	if err != nil {
		return nil, err
	}
	return &found, nil
}

// GetDaemonSet get a DaemonSet object
func (k8s K8sUtils) GetDaemonSet(namespacedName types.NamespacedName) (*appsv1.DaemonSet, error) {
	found := appsv1.DaemonSet{}
	err := k8s.Client.Get(context.TODO(), namespacedName, &found)
	if err != nil {
		return nil, err
	}
	return &found, nil
}

// GetStatefulSet get a StatefulSet object
func (k8s K8sUtils) GetStatefulSet(namespacedName types.NamespacedName) (*appsv1.StatefulSet, error) {
	found := appsv1.StatefulSet{}
	err := k8s.Client.Get(context.TODO(), namespacedName, &found)
	if err != nil {
		return nil, err
	}
	return &found, nil
}

// Exist checks if an object exists
func (k8s K8sUtils) Exist(namespacedName types.NamespacedName, kind schema.GroupVersionKind) bool {
	// Check if this Object exists
	found := unstructured.Unstructured{}
	found.SetGroupVersionKind(kind)
	err := k8s.Client.Get(context.TODO(), namespacedName, &found)
	if err != nil && errors.IsNotFound(err) {
		return false
	}
	if found.GetName() == namespacedName.Name && found.GetNamespace() == namespacedName.Namespace {
		return true
	}
	return false
}

func setOwner(scheme *runtime.Scheme, owner metav1.Object, controlled metav1.Object) error {
	if owner.GetNamespace() != controlled.GetNamespace() {
		log.Info("Controlled object has different namespace to owenr, ignoring", "Owner Namespace", owner.GetNamespace(), "Object Namespace", controlled.GetNamespace())
		return nil
	}
	// Set TNCO instance as the owner and controller
	if err := controllerutil.SetControllerReference(owner, controlled, scheme); err != nil {
		return err
	}
	return nil
}

// GetPodReadyCondition extracts the pod ready condition from the given status and returns that.
// Returns nil if the condition is not present.
func GetPodReadyCondition(status corev1.PodStatus) *corev1.PodCondition {
	_, condition := GetPodCondition(&status, corev1.PodReady)
	return condition
}

// GetPodCondition extracts the provided condition from the given status and returns that.
// Returns nil and -1 if the condition is not present, and the index of the located condition.
func GetPodCondition(status *corev1.PodStatus, conditionType corev1.PodConditionType) (int, *corev1.PodCondition) {
	if status == nil {
		return -1, nil
	}
	for i := range status.Conditions {
		if status.Conditions[i].Type == conditionType {
			return i, &status.Conditions[i]
		}
	}
	return -1, nil
}

/**
 * Executes a sequence of commands inside the given pod container
 */
func (k8s K8sUtils) ExecCommands(pod *corev1.Pod, containerName string, commands [][]string) (string, error) {
	for _, command := range commands {
		_, err := k8s.ExecCommand(pod, containerName, command...)
		if err != nil {
			return "", err
		}
	}

	return "", nil
}

/**
 * Executes a command inside the given pod container
 */
func (k8s K8sUtils) ExecCommand(pod *corev1.Pod, containerName string, command ...string) (string, error) {
	commands := []string{"sh", "-c"}
	for idx, _ := range command {
		commands = append(commands, command[idx])
	}
	log.Info("Executing command", "commands", commands)
	stdout, err := k8s.Exec(pod.Namespace, pod.Name, containerName, commands, nil, false)
	return stdout, err
}

/**
 * Uploads a file (given by 'reader') to the given pod
 */
func (k8s K8sUtils) UploadToK8s(pod *corev1.Pod, toPath, containerName string, reader io.Reader) (string, error) {
	b := strings.Builder{}

	dir, _ := filepath.Split(toPath)
	command := []string{"mkdir", "-p", dir}
	stdout, err := k8s.Exec(pod.Namespace, pod.Name, containerName, command, nil, false)
	b.WriteString(stdout)
	if err != nil {
		return stdout, err
	}

	command = []string{"cp", "/dev/stdin", toPath}
	stdout, err = k8s.Exec(pod.Namespace, pod.Name, containerName, command, reader, false)
	b.WriteString(stdout)
	if err != nil {
		return stdout, err
	}

	command = []string{"chmod", "+x", toPath}
	stdout, err = k8s.Exec(pod.Namespace, pod.Name, containerName, command, reader, false)
	b.WriteString(stdout)
	if err != nil {
		return "", err
	}

	return b.String(), nil
}

/*
 * Executes a command (array) in a given pod container, returning stderr as a string
 */
func (k8s K8sUtils) Exec(namespace, podName, containerName string, command []string, stdin io.Reader, tty bool) (string, error) {
	stdout := &bytes.Buffer{}

	req := k8s.ClientSet.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec").
		Timeout(10 * time.Second)
	scheme := runtime.NewScheme()
	if err := corev1.AddToScheme(scheme); err != nil {
		return "", fmt.Errorf("Exec error adding to scheme: %v", err)
	}

	req.VersionedParams(&corev1.PodExecOptions{
		Command:   command,
		Container: containerName,
		Stdin:     stdin != nil,
		Stdout:    true,
		Stderr:    true,
		TTY:       tty,
	}, runtime.NewParameterCodec(scheme))

	exec, err := remotecommand.NewSPDYExecutor(k8s.RestConfig, "POST", req.URL())
	if err != nil {
		return "", fmt.Errorf("Exec error while creating Executor: %v", err)
	}

	iteration := 0
	var lastError error
	for {
		if iteration > 3 {
			break
		}

		log.Info("Executing K8s command", "pod", podName, "container", containerName,
			"namespace", namespace, "command", command, "iteration", iteration, "tty", tty)

		err = exec.Stream(remotecommand.StreamOptions{
			Stdin:  stdin,
			Stdout: stdout,
			Stderr: stdout,
		})
		if err != nil {
			log.Info("Error executing K8s command", "pod", podName, "container", containerName,
				"namespace", namespace, "command", command, "iteration", iteration)

			if strings.Contains(err.Error(), "container not found") {
				lastError = err
				iteration++
				time.Sleep(1 * time.Second)
				continue
			}
			return stdout.String(), fmt.Errorf("Exec error for namespace %s pod %s container %s command %v in Stream: %v", namespace,
				podName, containerName, command, err)
		}

		log.Info("Executed K8s command", "pod", podName, "container", containerName,
			"namespace", namespace, "command", command, "iteration", iteration)

		break
	}

	if lastError != nil {
		return stdout.String(), fmt.Errorf("Exec error for namespace %s pod %s container %s command %v in Stream: %v", namespace,
			podName, containerName, command, lastError)
	}

	return stdout.String(), nil
}

/*
 * Starts a component by running the start transition script for the component type
 */
func (k8s K8sUtils) Start(pod *corev1.Pod, transition *v1alpha1.Transition, componentType string, containerName string) (string, error) {
	b := strings.Builder{}

	startCmd, err := LoadResource("/workspace/resources/" + componentType + "/start.sh")
	defer startCmd.Close()
	if err != nil {
		return b.String(), err
	}

	stdout, err := k8s.UploadToK8s(pod, "/tmp/start.sh", containerName, startCmd)
	b.WriteString(stdout)
	if err != nil {
		return b.String(), err
	}

	log.Info("AMF Start1")
	stdout, err = k8s.ExecCommand(pod, containerName, "/tmp/start.sh")
	log.Info("AMF Start2", "stdout", stdout, "err", err)
	b.WriteString(stdout)
	if err != nil {
		return b.String(), err
	}

	return b.String(), nil
}

/*
 * Stops a component by running the stop transition script for the component type
 */
func (k8s K8sUtils) Stop(pod *corev1.Pod, transition *v1alpha1.Transition, componentType string, containerName string) (string, error) {
	b := strings.Builder{}

	stopCmd, err := LoadResource("/workspace/resources/" + componentType + "/stop.sh")
	if err != nil {
		return b.String(), err
	}

	stdout, err := k8s.UploadToK8s(pod, "/tmp/stop.sh", containerName, stopCmd)
	b.WriteString(stdout)
	if err != nil {
		return b.String(), err
	}

	stdout, err = k8s.ExecCommand(pod, containerName, "/tmp/stop.sh")
	b.WriteString(stdout)
	if err != nil {
		return b.String(), err
	}

	return b.String(), nil
}

func (k8s K8sUtils) RunOperation(pod *corev1.Pod, transition *v1alpha1.Transition, componentType string, containerName string, commands []string) (string, error) {
	b := strings.Builder{}
	for _, command := range commands {
		stdout, err := k8s.ExecCommand(pod, containerName, command)
		b.WriteString(command)
		b.WriteString(" : ")
		b.WriteString(stdout)
		b.WriteString("\n")
		if err != nil {
			return b.String(), err
		}
	}

	return b.String(), nil
}
