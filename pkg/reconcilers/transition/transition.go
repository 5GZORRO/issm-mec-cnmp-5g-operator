package transition

import (
	"github.ibm.com/Steve-Glover/5GOperator/pkg/common"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// Transition is a type to manage k8s objects for Transition 5G
type Transition struct {
	K8sUtils *common.K8sUtils
}

func NewTransition(k8sUtils *common.K8sUtils) *Transition {
	transition := Transition{
		K8sUtils: k8sUtils,
	}

	return &transition
}

// Finalize method for Transition.  Executed on uninstall
func (a Transition) Finalize(request reconcile.Request, instance interface{}) {
	// Nothing to do
}
