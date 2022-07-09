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

package nrf

import (
	fivegv1alpha1 "github.com/5GZORRO/issm-mec-cnmp-5g-operator/api/v1alpha1"
	common "github.com/5GZORRO/issm-mec-cnmp-5g-operator/pkg/common"
	nrfv13 "github.com/5GZORRO/issm-mec-cnmp-5g-operator/pkg/reconcilers/nrf/v13"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// NrfReconcilerFactory factory class to get NrfReconciler for required installation version
type NrfReconcilerFactory struct {
	K8sUtils *common.K8sUtils
	Client   client.Client
}

// Reconciler get a FiveGReconciler to reconcile the required version of Nrf
func (f NrfReconcilerFactory) Reconciler() common.FiveGReconciler {
	// use cr.Spec.version to determine which version to install
	// use cr.Status.versions.reconciled to determine which version is currently installed

	nrf := NewNrf(f.K8sUtils, f.K8sUtils, f.Client)
	reconcileBuilder := f.version13ReconcileBuilder()

	return common.FiveGReconcilerFacade{
		ReconcileFunc:    reconcileBuilder.Reconcile,
		UpdateStatusFunc: nrf.UpdateStatus,
		FinalizeFunc:     nrf.Finalize,
	}
}

func (f NrfReconcilerFactory) version13ReconcileBuilder() common.FiveGReconcileFuncBuilder {
	return common.FiveGReconcileFuncBuilder{
		K8sUtils: f.K8sUtils,
		PodFuncs: []common.ReconcilePod{
			func(instance interface{}) *corev1.Pod {
				cr := instance.(*fivegv1alpha1.Nrf)
				return nrfv13.NrfPod(cr)
			},
		},
		ConfigMapFuncs: []common.ReconcileConfigMap{
			func(instance interface{}) (*corev1.ConfigMap, error) {
				cr := instance.(*fivegv1alpha1.Nrf)
				return nrfv13.NrfConfigMap(cr)
			},
		},
	}
}
