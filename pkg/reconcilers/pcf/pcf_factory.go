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

package pcf

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	fivegv1alpha1 "github.com/5GZORRO/issm-mec-cnmp-5g-operator/api/v1alpha1"
	"github.com/5GZORRO/issm-mec-cnmp-5g-operator/pkg/common"
	pcfv13 "github.com/5GZORRO/issm-mec-cnmp-5g-operator/pkg/reconcilers/pcf/v13"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// PcfReconcilerFactory factory class to get PcfReconciler for required installation version
type PcfReconcilerFactory struct {
	K8sUtils *common.K8sUtils
	Client   client.Client
}

// Reconciler get a FiveGReconciler to reconcile the required version of Pcf
func (f PcfReconcilerFactory) Reconciler() common.FiveGReconciler {
	// use cr.Spec.version to determine which version to install
	// use cr.Status.versions.reconciled to determine which version is currently installed

	pcf := NewPcf(f.K8sUtils, f.K8sUtils, f.Client)
	reconcileBuilder := f.version13ReconcileBuilder()

	return common.FiveGReconcilerFacade{
		ReconcileFunc:    reconcileBuilder.Reconcile,
		UpdateStatusFunc: pcf.UpdateStatus,
		FinalizeFunc:     pcf.Finalize,
	}
}

func NewSHA256(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

func (f PcfReconcilerFactory) version13ReconcileBuilder() common.FiveGReconcileFuncBuilder {
	return common.FiveGReconcileFuncBuilder{
		K8sUtils: f.K8sUtils,
		PodFuncs: []common.ReconcilePod{
			func(instance interface{}) *corev1.Pod {
				cr := instance.(*fivegv1alpha1.Pcf)

				cm := &corev1.ConfigMap{}
				err := f.K8sUtils.Client.Get(context.TODO(), types.NamespacedName{Name: cr.Name, Namespace: cr.Namespace}, cm)
				if err != nil && errors.IsNotFound(err) {
					// TODO
				} else if err != nil {
					// TODO
				}
				return pcfv13.PcfPod(cr, cm)
			},
		},
		ConfigMapFuncs: []common.ReconcileConfigMap{
			func(instance interface{}) (*corev1.ConfigMap, error) {
				cr := instance.(*fivegv1alpha1.Pcf)
				return pcfv13.PcfConfigMap(cr)
			},
		},
	}
}
