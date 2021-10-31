package n3iwf

import (
	fivegv1alpha1 "github.ibm.com/Steve-Glover/5GOperator/api/v1alpha1"
	"github.ibm.com/Steve-Glover/5GOperator/pkg/common"
	n3iwfv13 "github.ibm.com/Steve-Glover/5GOperator/pkg/reconcilers/n3iwf/v13"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// N3iwfReconcilerFactory factory class to get N3iwfReconciler for required installation version
type N3iwfReconcilerFactory struct {
	K8sUtils *common.K8sUtils
	Client   client.Client
}

// Reconciler get a N3iwfReconciler to reconcile the required version of N3iwf
func (f N3iwfReconcilerFactory) Reconciler() common.FiveGReconciler {
	// use cr.Spec.version to determine which version to install
	// use cr.Status.versions.reconciled to determine which version is currently installed

	n3iwf := NewN3iwf(f.K8sUtils, f.K8sUtils, f.Client)
	reconcileBuilder := f.version13ReconcileBuilder()

	return common.FiveGReconcilerFacade{
		ReconcileFunc:    reconcileBuilder.Reconcile,
		UpdateStatusFunc: n3iwf.UpdateStatus,
		FinalizeFunc:     n3iwf.Finalize,
	}
}

func (f N3iwfReconcilerFactory) version13ReconcileBuilder() common.FiveGReconcileFuncBuilder {
	return common.FiveGReconcileFuncBuilder{
		K8sUtils: f.K8sUtils,
		PodFuncs: []common.ReconcilePod{
			func(instance interface{}) *corev1.Pod {
				cr := instance.(*fivegv1alpha1.N3iwf)
				return n3iwfv13.N3iwfPod(cr)
			},
		},
		ConfigMapFuncs: []common.ReconcileConfigMap{
			func(instance interface{}) (*corev1.ConfigMap, error) {
				cr := instance.(*fivegv1alpha1.N3iwf)
				return n3iwfv13.N3iwfConfigMap(cr)
			},
		},
	}
}
