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

package subscriber

import (
	"github.com/5GZORRO/issm-mec-cnmp-5g-operator/pkg/common"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// SubscriberReconcilerFactory factory class to get SmfReconciler for required installation version
type SubscriberReconcilerFactory struct {
	K8sUtils *common.K8sUtils
	Client   client.Client
}

// Reconciler get a SubscriberReconciler to reconcile the required version of Subscriber
func (f SubscriberReconcilerFactory) Reconciler() common.FiveGReconciler {
	// use cr.Spec.version to determine which version to install
	// use cr.Status.versions.reconciled to determine which version is currently installed

	subscriber := NewSubscriber(f.K8sUtils, f.K8sUtils, f.Client)
	reconcileBuilder := f.version13ReconcileBuilder()

	return common.FiveGReconcilerFacade{
		ReconcileFunc:    reconcileBuilder.Reconcile,
		UpdateStatusFunc: subscriber.UpdateStatus,
		FinalizeFunc:     subscriber.Finalize,
	}
}

func (f SubscriberReconcilerFactory) version13ReconcileBuilder() common.FiveGReconcileFuncBuilder {
	return common.FiveGReconcileFuncBuilder{
		K8sUtils: f.K8sUtils,
	}
}
