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
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// 5GReconciler interface to be implemented by types which need to create application k8s objects
type FiveGReconciler interface {
	// Reconcile run reconciliation
	Reconcile(request reconcile.Request, cr interface{}) (controllerutil.OperationResult, reconcile.Result, error)

	// UpdateStatus  update the status of a CR based on status of reconciled objects status
	UpdateStatus(opResult controllerutil.OperationResult, cr Entity, reconcileErr error)

	// Finalize cleanup on deletion
	Finalize(request reconcile.Request, cr interface{})
}
