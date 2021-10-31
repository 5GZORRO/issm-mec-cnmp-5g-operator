package mongo

import (
	fivegv1alpha1 "github.ibm.com/Steve-Glover/5GOperator/api/v1alpha1"
	common "github.ibm.com/Steve-Glover/5GOperator/pkg/common"
	mongov13 "github.ibm.com/Steve-Glover/5GOperator/pkg/reconcilers/mongo/v13"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// MongoReconcilerFactory factory class to get MongoReconciler for required installation version
type MongoReconcilerFactory struct {
	K8sUtils *common.K8sUtils
	Client   client.Client
}

// Reconciler get a FiveGReconciler to reconcile the required version of Mongo
func (f MongoReconcilerFactory) Reconciler() common.FiveGReconciler {
	// use cr.Spec.version to determine which version to install
	// use cr.Status.versions.reconciled to determine which version is currently installed

	mongo := NewMongo(f.K8sUtils, f.K8sUtils, f.Client)
	reconcileBuilder := f.version13ReconcileBuilder()

	return common.FiveGReconcilerFacade{
		ReconcileFunc:    reconcileBuilder.Reconcile,
		UpdateStatusFunc: mongo.UpdateStatus,
		FinalizeFunc:     mongo.Finalize,
	}
}

func (f MongoReconcilerFactory) version13ReconcileBuilder() common.FiveGReconcileFuncBuilder {
	return common.FiveGReconcileFuncBuilder{
		K8sUtils: f.K8sUtils,
		PodFuncs: []common.ReconcilePod{
			func(instance interface{}) *corev1.Pod {
				cr := instance.(*fivegv1alpha1.Mongo)
				return mongov13.MongoPod(cr)
			},
		},
		ServiceFuncs: []common.ReconcileService{
			func(instance interface{}) *corev1.Service {
				cr := instance.(*fivegv1alpha1.Mongo)
				return mongov13.MongoService(cr)
			},
		},
	}
}
