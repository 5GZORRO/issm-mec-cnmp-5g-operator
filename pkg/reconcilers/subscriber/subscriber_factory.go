package subscriber

import (
	"github.ibm.com/Steve-Glover/5GOperator/pkg/common"
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
