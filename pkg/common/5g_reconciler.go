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
