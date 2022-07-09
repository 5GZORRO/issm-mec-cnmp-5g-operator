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

package udr

import (
	fivegv1alpha1 "github.com/5GZORRO/issm-mec-cnmp-5g-operator/api/v1alpha1"
	"github.com/5GZORRO/issm-mec-cnmp-5g-operator/pkg/common"
	udrv13 "github.com/5GZORRO/issm-mec-cnmp-5g-operator/pkg/reconcilers/udr/v13"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// UdrReconcilerFactory factory class to get UdrReconciler for required installation version
type UdrReconcilerFactory struct {
	K8sUtils *common.K8sUtils
	Client   client.Client
}

// Reconciler get a FiveGReconciler to reconcile the required version of Udr
func (f UdrReconcilerFactory) Reconciler() common.FiveGReconciler {
	// use cr.Spec.version to determine which version to install
	// use cr.Status.versions.reconciled to determine which version is currently installed

	udr := NewUdr(f.K8sUtils, f.K8sUtils, f.Client)
	reconcileBuilder := f.version13ReconcileBuilder()

	return common.FiveGReconcilerFacade{
		ReconcileFunc:    reconcileBuilder.Reconcile,
		UpdateStatusFunc: udr.UpdateStatus,
		FinalizeFunc:     udr.Finalize,
	}
}

func (f UdrReconcilerFactory) version13ReconcileBuilder() common.FiveGReconcileFuncBuilder {
	return common.FiveGReconcileFuncBuilder{
		K8sUtils: f.K8sUtils,
		PodFuncs: []common.ReconcilePod{
			func(instance interface{}) *corev1.Pod {
				cr := instance.(*fivegv1alpha1.Udr)
				return udrv13.UdrPod(cr)
			},
		},
		ConfigMapFuncs: []common.ReconcileConfigMap{
			func(instance interface{}) (*corev1.ConfigMap, error) {
				cr := instance.(*fivegv1alpha1.Udr)
				return udrv13.UdrConfigMap(cr)
			},
		},
	}
}
