package amf

import (
	fivegv1alpha1 "github.ibm.com/Steve-Glover/5GOperator/api/v1alpha1"
	"github.ibm.com/Steve-Glover/5GOperator/pkg/common"
	amfv13 "github.ibm.com/Steve-Glover/5GOperator/pkg/reconcilers/amf/v13"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// AmfReconcilerFactory factory class to get ApolloReconciler for required installation version
type AmfReconcilerFactory struct {
	K8sUtils *common.K8sUtils
	Client   client.Client
}

// Reconciler get a FiveGReconciler to reconcile the required version of Amf
func (f AmfReconcilerFactory) Reconciler() common.FiveGReconciler {
	// use cr.Spec.version to determine which version to install
	// use cr.Status.versions.reconciled to determine which version is currently installed

	amf := NewAmf(f.K8sUtils, f.K8sUtils, f.Client)
	reconcileBuilder := f.version13ReconcileBuilder()

	return common.FiveGReconcilerFacade{
		ReconcileFunc:    reconcileBuilder.Reconcile,
		UpdateStatusFunc: amf.UpdateStatus,
		FinalizeFunc:     amf.Finalize,
	}
}

func (f AmfReconcilerFactory) version13ReconcileBuilder() common.FiveGReconcileFuncBuilder {
	return common.FiveGReconcileFuncBuilder{
		K8sUtils: f.K8sUtils,
		PodFuncs: []common.ReconcilePod{
			func(instance interface{}) *corev1.Pod {
				cr := instance.(*fivegv1alpha1.Amf)
				return amfv13.AmfPod(cr)
			},
		},
		ConfigMapFuncs: []common.ReconcileConfigMap{
			func(instance interface{}) (*corev1.ConfigMap, error) {
				cr := instance.(*fivegv1alpha1.Amf)
				return amfv13.AmfConfigMap(cr)
			},
		},
		//StateFuncs: []common.ReconcileState{
		//	func(instance interface{}) (string, error) {
		//		cr := instance.(*fivegv1alpha1.Amf)
		//		return amfv13.AmfTargetState(cr)
		//	},
		//},
	}
}
