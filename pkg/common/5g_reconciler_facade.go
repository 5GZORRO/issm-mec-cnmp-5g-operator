package common

import (
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// Reconcile run reconciliation
type Reconcile func(reconcile.Request, interface{}) (controllerutil.OperationResult, reconcile.Result, error)

// UpdateStatus update the status of a CR based on status of reconciled objects status
type UpdateStatus func(controllerutil.OperationResult, Entity, error)

// Finalize cleanup on deletion
type Finalize func(reconcile.Request, interface{})

// FiveGReconcilerFacade Facade object implementing FiveGReconciler
type FiveGReconcilerFacade struct {
	// Functions
	ReconcileFunc    Reconcile
	UpdateStatusFunc UpdateStatus
	FinalizeFunc     Finalize
}

// Reconcile Facade for FiveGReconcilerFacade.Reconcile function
func (f FiveGReconcilerFacade) Reconcile(request reconcile.Request, cr interface{}) (controllerutil.OperationResult, reconcile.Result, error) {
	return f.ReconcileFunc(request, cr)
}

// UpdateStatus Facade for FiveGReconcilerFacade.UpdateStatus function
func (f FiveGReconcilerFacade) UpdateStatus(opResult controllerutil.OperationResult,
	cr Entity, err error) {
	f.UpdateStatusFunc(opResult, cr, err)
}

// Finalize Facade for FiveGReconcilerFacade.Finalize function
func (f FiveGReconcilerFacade) Finalize(request reconcile.Request, cr interface{}) {
	f.FinalizeFunc(request, cr)
}
