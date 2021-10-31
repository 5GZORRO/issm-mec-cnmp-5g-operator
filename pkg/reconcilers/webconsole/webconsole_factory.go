package webconsole

import (
	fivegv1alpha1 "github.ibm.com/Steve-Glover/5GOperator/api/v1alpha1"
	"github.ibm.com/Steve-Glover/5GOperator/pkg/common"
	webconsolev13 "github.ibm.com/Steve-Glover/5GOperator/pkg/reconcilers/webconsole/v13"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// WebconsoleReconcilerFactory factory class to get WebconsoleReconciler for required installation version
type WebconsoleReconcilerFactory struct {
	K8sUtils *common.K8sUtils
	Client   client.Client
}

// Reconciler get a FiveGReconciler to reconcile the required version of Webconsole
func (f WebconsoleReconcilerFactory) Reconciler() common.FiveGReconciler {
	// use cr.Spec.version to determine which version to install
	// use cr.Status.versions.reconciled to determine which version is currently installed

	webconsole := NewWebconsole(f.K8sUtils, f.K8sUtils, f.Client)
	reconcileBuilder := f.version13ReconcileBuilder()

	return common.FiveGReconcilerFacade{
		ReconcileFunc:    reconcileBuilder.Reconcile,
		UpdateStatusFunc: webconsole.UpdateStatus,
		FinalizeFunc:     webconsole.Finalize,
	}
}

func (f WebconsoleReconcilerFactory) version13ReconcileBuilder() common.FiveGReconcileFuncBuilder {
	return common.FiveGReconcileFuncBuilder{
		K8sUtils: f.K8sUtils,
		PodFuncs: []common.ReconcilePod{
			func(instance interface{}) *corev1.Pod {
				cr := instance.(*fivegv1alpha1.Webconsole)
				return webconsolev13.WebconsolePod(cr)
			},
		},
		ServiceFuncs: []common.ReconcileService{
			func(instance interface{}) *corev1.Service {
				cr := instance.(*fivegv1alpha1.Webconsole)
				return webconsolev13.WebconsoleService(cr)
			},
		},
		ConfigMapFuncs: []common.ReconcileConfigMap{
			func(instance interface{}) (*corev1.ConfigMap, error) {
				cr := instance.(*fivegv1alpha1.Webconsole)
				return webconsolev13.WebconsoleConfigMap(cr)
			},
		},
	}
}
