package ausf

import (
	"context"
	"fmt"
	"github.ibm.com/Steve-Glover/5GOperator/api/v1alpha1"
	fivegv1alpha1 "github.ibm.com/Steve-Glover/5GOperator/api/v1alpha1"
	"github.ibm.com/Steve-Glover/5GOperator/pkg/common"
	v12 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"strings"
)

var log = logf.Log.WithName("ausf")

// Ausf is a type to manage k8s objects for Ausf 5G
type Ausf struct {
	common.EntityHandlerImpl
	k8sUtils        *common.K8sUtils
	operationRunner common.OperationRunner
	client.Client
}

func NewAusf(operationRunner common.OperationRunner, k8sUtils *common.K8sUtils, client client.Client) *Ausf {
	ausf := Ausf{
		EntityHandlerImpl: *common.NewEntityHandler(),
		operationRunner:   operationRunner,
		k8sUtils:          k8sUtils,
		Client:            client,
	}

	ausf.EntityHandlerImpl.AddTransitionFunction("start", ausf.Start)
	ausf.EntityHandlerImpl.AddTransitionFunction("stop", ausf.Stop)

	return &ausf
}

// UpdateStatus called to update CR status
func (a Ausf) UpdateStatus(opResult controllerutil.OperationResult, instance common.Entity,
	reconcileErr error) {
	cr := instance.(*v1alpha1.Ausf)

	if reconcileErr != nil {
		cr.Status.SetError(reconcileErr)
	} else {
		obj, err := a.k8sUtils.GetPod(
			types.NamespacedName{Name: cr.Name, Namespace: cr.Namespace},
		)

		// TODO check err
		if err == nil {
			podReadyCondition := common.GetPodReadyCondition(obj.Status)
			isReconciled := false
			if podReadyCondition != nil && obj.Status.PodIP != "" {
				// the pod is ready, update the Mongo CR
				cr.Status.Outputs.IpAddr = obj.Status.PodIP
				cr.Status.Outputs.PodName = cr.Name
				isReconciled = true
			}

			cr.Status.Update(cr.ObjectMeta, isReconciled, opResult)
		}
	}
}

// Finalize method for Ausf.  Executed on uninstall
func (a Ausf) Finalize(request reconcile.Request, instance interface{}) {
	// Nothing to do
}

func (a Ausf) Create(transition *v1alpha1.Transition) (common.Entity, error) {
	instance := &fivegv1alpha1.Ausf{
		ObjectMeta: metav1.ObjectMeta{
			Name:      transition.Spec.Config.ResourceName,
			Namespace: transition.Spec.Config.ResourceNamespace,
		},
		Spec: fivegv1alpha1.AusfSpec{
			Config: &fivegv1alpha1.AusfConfig{
				PodSettings:  common.GetPodSettings(transition),
				ImageUrl:     transition.Spec.Config.Properties["image"],
				Name:         transition.Spec.Config.Properties["amf_name"],
				NrfIPAddress: transition.Spec.Config.Properties["nrf_ip_address"],
				NrfPort:      transition.Spec.Config.Properties["nrf_port"],
			},
			Internal: nil,
		},
	}

	err := a.Client.Create(context.TODO(), instance)
	if err != nil {
		return nil, err
	}

	return instance, nil
}

func (a Ausf) Start(transition *v1alpha1.Transition) (string, error) {
	cr, err := a.getCRForTransition(transition)
	if err != nil {
		return "", err
	}

	pod, err := a.getPodForCR(cr)
	if err != nil {
		return "", err
	}

	return a.operationRunner.Start(pod, transition, "ausf", "ausf")
}

func (a Ausf) Stop(transition *v1alpha1.Transition) (string, error) {
	cr, err := a.getCRForTransition(transition)
	if err != nil {
		return "", err
	}

	pod, err := a.getPodForCR(cr)
	if err != nil {
		return "", err
	}

	return a.operationRunner.Stop(pod, transition, "ausf", "ausf")
}

func (a Ausf) GetCR(transitionCR *fivegv1alpha1.Transition) (common.Entity, error) {
	return a.k8sUtils.GetAusf(types.NamespacedName{
		Namespace: transitionCR.Spec.Config.ResourceNamespace,
		Name:      transitionCR.Spec.Config.ResourceName,
	})
}

func (a Ausf) CheckReady(transitionCR *fivegv1alpha1.Transition, entity common.Entity) (bool, error) {
	return strings.ToLower(transitionCR.Spec.Config.TransitionName) == "start", nil
}

func (a Ausf) IsReconciled(transitionCR *fivegv1alpha1.Transition) (bool, error) {
	cr, err := a.getCRForTransition(transitionCR)
	if err != nil {
		transitionCR.Status.SetError(err)
		return false, err
	}
	if cr == nil {
		err = fmt.Errorf("Cannot find cr")
		transitionCR.Status.SetError(err)
		return false, err
	}
	return cr.Status.IsReconciled(), nil
}

func (a Ausf) getCRForTransition(transition *v1alpha1.Transition) (*v1alpha1.Ausf, error) {
	cr, err := a.k8sUtils.GetAusf(types.NamespacedName{
		Namespace: transition.Spec.Config.ResourceNamespace,
		Name:      transition.Spec.Config.ResourceName,
	})
	if err != nil {
		return nil, err
	}

	return cr, nil
}

func (a Ausf) getPodForCR(cr *v1alpha1.Ausf) (*v12.Pod, error) {
	pod, err := a.k8sUtils.GetPod(types.NamespacedName{
		Namespace: cr.Namespace,
		Name:      cr.Name,
	})
	if err != nil {
		return nil, err
	}

	return pod, nil
}
