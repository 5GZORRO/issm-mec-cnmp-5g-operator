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
	routev1 "github.com/openshift/api/route/v1"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	ingressv1b1 "k8s.io/api/extensions/v1beta1"
	networking "k8s.io/api/networking/v1"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var isOpenshift bool

func init() {
	isOpenshift, _ = VerifyRouteAPI()
}

type ReconcileConfigMap func(interface{}) (*corev1.ConfigMap, error)

type ReconcileService func(interface{}) *corev1.Service

type ReconcilePod func(interface{}) *corev1.Pod

type ReconcileDeploymennt func(interface{}) *appsv1.Deployment

type ReconcileStatefulSet func(interface{}) *appsv1.StatefulSet

type ReconcileRoute func(interface{}) *routev1.Route

type ReconcileIngress func(interface{}) *ingressv1b1.Ingress

type ReconcileJob func(interface{}) *batchv1.Job

type ReconcilePodDisruptionBudget func(interface{}) *policyv1beta1.PodDisruptionBudget

type ReconcilePVC func(interface{}) *corev1.PersistentVolumeClaim

type ReconcileDaemonSet func(interface{}) *appsv1.DaemonSet

type ReconcileServiceAccount func(interface{}) *corev1.ServiceAccount

type ReconcileRoleBinding func(interface{}) *rbacv1.RoleBinding

type ReconcileNetworkPolicy func(interface{}) *networking.NetworkPolicy

type ReconcileTransition func(interface{}) error

//type ReconcileState func(interface{}) (string, error)

type FiveGReconcileFuncBuilder struct {
	K8sUtils                 *K8sUtils
	ConfigMapFuncs           []ReconcileConfigMap
	ServiceFuncs             []ReconcileService
	PodFuncs                 []ReconcilePod
	DeploymentFuncs          []ReconcileDeploymennt
	DaemonSetFuncs           []ReconcileDaemonSet
	StatefulSetFuncs         []ReconcileStatefulSet
	RouteFuncs               []ReconcileRoute
	IngressFuncs             []ReconcileIngress
	JobFuncs                 []ReconcileJob
	PodDisruptionBudgetFuncs []ReconcilePodDisruptionBudget
	PVCFuncs                 []ReconcilePVC
	ServiceAccountFuncs      []ReconcileServiceAccount
	RoleBindingFuncs         []ReconcileRoleBinding
	NetworkPolicyFuncs       []ReconcileNetworkPolicy
	TransitionFuncs          []ReconcileTransition
	//StateFuncs          	 []ReconcileState
}

func updateOperationResult(current controllerutil.OperationResult, update controllerutil.OperationResult) controllerutil.OperationResult {
	if current == controllerutil.OperationResultNone && update != "" && update != controllerutil.OperationResultNone {
		return update
	}
	return current
}

func (r FiveGReconcileFuncBuilder) Reconcile(request reconcile.Request, instance interface{}) (controllerutil.OperationResult, reconcile.Result, error) {
	cr := instance.(metav1.Object)

	opResult := controllerutil.OperationResultNone

	if r.ConfigMapFuncs != nil {
		for _, c := range r.ConfigMapFuncs {
			configmap, err := c(instance)
			if err != nil {
				return "", reconcile.Result{}, err
			}
			if configmap != nil {
				result, err := r.K8sUtils.ReconcileConfigMap(cr, configmap)

				if result != controllerutil.OperationResultNone || err != nil {
					log.Info("ReconcileConfigmap", "configmap", configmap, "result", result, "err", err, "opResult", opResult)
				}

				if err != nil {
					return "", reconcile.Result{}, err
				}
				opResult = updateOperationResult(opResult, result)
			}
		}
	}

	if r.ServiceFuncs != nil {
		for _, s := range r.ServiceFuncs {
			service := s(instance)
			if service != nil {
				result, err := r.K8sUtils.ReconcileService(cr, service)

				if result != controllerutil.OperationResultNone || err != nil {
					log.Info("ReconcileService", "service", service, "result", result, "err", err, "opResult", opResult)
				}

				if err != nil {
					return "", reconcile.Result{}, err
				}

				opResult = updateOperationResult(opResult, result)
			}
		}
	}

	if r.DeploymentFuncs != nil {
		for _, d := range r.DeploymentFuncs {
			deployment := d(instance)
			if deployment != nil {
				result, err := r.K8sUtils.ReconcileDeployment(cr, deployment)

				if result != controllerutil.OperationResultNone || err != nil {
					log.Info("ReconcileDeployment", "deployment", deployment, "result", result, "err", err, "opResult", opResult)
				}

				if err != nil {
					return "", reconcile.Result{}, err
				}

				opResult = updateOperationResult(opResult, result)
			}
		}
	}

	if r.PodFuncs != nil {
		for _, d := range r.PodFuncs {
			pod := d(instance)
			if pod != nil {
				result, err := r.K8sUtils.ReconcilePod(cr, pod)

				if result != controllerutil.OperationResultNone || err != nil {
					log.Info("ReconcilePod", "pod", pod, "result", result, "err", err, "opResult", opResult)
				}

				if err != nil {
					return "", reconcile.Result{}, err
				}

				opResult = updateOperationResult(opResult, result)
			}
		}
	}

	if r.DaemonSetFuncs != nil {
		for _, ds := range r.DaemonSetFuncs {
			daemonset := ds(instance)
			if daemonset != nil {
				result, err := r.K8sUtils.ReconcileDaemonSet(cr, daemonset)
				if err != nil {
					return "", reconcile.Result{}, err
				}
				opResult = updateOperationResult(opResult, result)
			}
		}
	}

	if r.StatefulSetFuncs != nil {
		for _, s := range r.StatefulSetFuncs {
			statefulset := s(instance)
			if statefulset != nil {
				result, err := r.K8sUtils.ReconcileStatefulSet(cr, statefulset)
				if err != nil {
					return "", reconcile.Result{}, err
				}
				opResult = updateOperationResult(opResult, result)
			}
		}
	}

	if isOpenshift {
		if r.RouteFuncs != nil {
			for _, i := range r.RouteFuncs {
				route := i(instance)
				if route != nil {
					result, err := r.K8sUtils.ReconcileRoute(cr, route)
					if err != nil {
						return "", reconcile.Result{}, err
					}
					opResult = updateOperationResult(opResult, result)
				}
			}
		}
	}

	if !isOpenshift {
		if r.IngressFuncs != nil {
			for _, i := range r.IngressFuncs {
				ingress := i(instance)
				if ingress != nil {
					result, err := r.K8sUtils.ReconcileIngress(cr, ingress)
					if err != nil {
						return "", reconcile.Result{}, err
					}
					opResult = updateOperationResult(opResult, result)
				}
			}
		}
	}

	if r.JobFuncs != nil {
		for _, j := range r.JobFuncs {
			job := j(instance)
			if job != nil {
				result, err := r.K8sUtils.ReconcileJob(cr, job)
				if err != nil {
					return "", reconcile.Result{}, err
				}
				opResult = updateOperationResult(opResult, result)
			}
		}
	}

	if r.PodDisruptionBudgetFuncs != nil {
		for _, p := range r.PodDisruptionBudgetFuncs {
			pdb := p(instance)
			if pdb != nil {
				result, err := r.K8sUtils.ReconcilePodDisruptionBudget(cr, pdb)
				if err != nil {
					return "", reconcile.Result{}, err
				}
				opResult = updateOperationResult(opResult, result)
			}
		}
	}

	if r.PVCFuncs != nil {
		for _, c := range r.PVCFuncs {
			pvc := c(instance)
			if pvc != nil {
				result, err := r.K8sUtils.ReconcilePVC(cr, pvc)
				if err != nil {
					return "", reconcile.Result{}, err
				}
				opResult = updateOperationResult(opResult, result)
			}
		}
	}

	if r.ServiceAccountFuncs != nil {
		for _, sa := range r.ServiceAccountFuncs {
			serviceaccount := sa(instance)
			if serviceaccount != nil {
				result, err := r.K8sUtils.ReconcileServiceAccount(cr, serviceaccount)
				if err != nil {
					return "", reconcile.Result{}, err
				}
				opResult = updateOperationResult(opResult, result)
			}
		}
	}

	if r.RoleBindingFuncs != nil {
		for _, rb := range r.RoleBindingFuncs {
			rolebinding := rb(instance)
			if rolebinding != nil {
				result, err := r.K8sUtils.ReconcileRoleBinding(cr, rolebinding)
				if err != nil {
					return "", reconcile.Result{}, err
				}
				opResult = updateOperationResult(opResult, result)
			}
		}
	}

	if r.NetworkPolicyFuncs != nil {
		for _, c := range r.NetworkPolicyFuncs {
			pvc := c(instance)
			if pvc != nil {
				result, err := r.K8sUtils.ReconcileNetworkPolicy(cr, pvc)
				if err != nil {
					return "", reconcile.Result{}, err
				}
				opResult = updateOperationResult(opResult, result)
			}
		}
	}

	if r.TransitionFuncs != nil {
		for _, c := range r.TransitionFuncs {
			err := c(instance)
			if err != nil {
				return "", reconcile.Result{}, err
			}
		}
	}

	//if r.StateFuncs != nil {
	//	for _, c := range r.StateFuncs {
	//		targetState, err := c(instance)
	//		if err != nil {
	//			return "", reconcile.Result{}, err
	//		}
	//
	//		switch targetState {
	//		case strings.ToLower("created"):
	//
	//		}
	//		result, err := r.K8sUtils.ReconcileCRState(cr, statefulset)
	//		if err != nil {
	//			return "", reconcile.Result{}, err
	//		}
	//		opResult = updateOperationResult(opResult, result)
	//
	//
	//	}
	//}

	return opResult, reconcile.Result{}, nil
}
