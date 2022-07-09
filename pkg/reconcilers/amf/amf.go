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

package amf

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
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"strings"
)

// Amf is a type to manage k8s objects for Amf 5G
type Amf struct {
	*common.EntityHandlerImpl
	client.Client
	operationRunner common.OperationRunner
	k8sUtils        *common.K8sUtils
}

func NewAmf(operationRunner common.OperationRunner, k8sUtils *common.K8sUtils, client client.Client) *Amf {
	amf := Amf{
		EntityHandlerImpl: common.NewEntityHandler(),
		k8sUtils:          k8sUtils,
		operationRunner:   operationRunner,
		Client:            client,
	}

	// these are the transitions supported on this CR
	amf.EntityHandlerImpl.AddTransitionFunction("start", amf.Start)
	amf.EntityHandlerImpl.AddTransitionFunction("stop", amf.Stop)

	return &amf
}

func (a Amf) IsReconciled(transitionCR *fivegv1alpha1.Transition) (bool, error) {
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

func (a Amf) CheckCreated(transitionCR *fivegv1alpha1.Transition) (bool, error) {
	cr, err := a.getCRForTransition(transitionCR)
	if err != nil {
		return false, err
	}

	pod, err := a.getPodForCR(cr)
	if err != nil {
		return false, err
	}

	podReadyCondition := common.GetPodReadyCondition(pod.Status)
	if podReadyCondition != nil && pod.Status.PodIP != "" {
		// the pod is ready, update the AMF parent CR
		cr.Status.Outputs.IpAddr = pod.Status.PodIP
		cr.Status.Outputs.PodName = cr.Name
	}

	return podReadyCondition != nil, nil
}

// UpdateStatus called to update CR status
func (a Amf) UpdateStatus(opResult controllerutil.OperationResult, instance common.Entity, reconcileErr error) {
	cr := instance.(*v1alpha1.Amf)

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

func (a Amf) Create(transition *v1alpha1.Transition) (common.Entity, error) {
	instance := &fivegv1alpha1.Amf{
		//TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      transition.Spec.Config.ResourceName,
			Namespace: transition.Spec.Config.ResourceNamespace,
		},
		Spec: fivegv1alpha1.AmfSpec{
			Config: &fivegv1alpha1.AmfConfig{
				PodSettings:  common.GetPodSettings(transition),
				ImageUrl:     transition.Spec.Config.Properties["image"],
				Name:         transition.Spec.Config.Properties["amf_name"],
				NrfIPAddress: transition.Spec.Config.Properties["nrf_ip_address"],
				NrfPort:      transition.Spec.Config.Properties["nrf_port"],
				Mnc:          transition.Spec.Config.Properties["mnc"],
				Mcc:          transition.Spec.Config.Properties["mcc"],
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

func (a Amf) Start(transition *v1alpha1.Transition) (string, error) {
	cr, err := a.getCRForTransition(transition)
	if err != nil {
		return "", err
	}

	pod, err := a.getPodForCR(cr)
	if err != nil {
		return "", err
	}

	return a.operationRunner.Start(pod, transition, "amf", "amf")
}

func (a Amf) Stop(transition *v1alpha1.Transition) (string, error) {
	cr, err := a.getCRForTransition(transition)
	if err != nil {
		return "", err
	}

	pod, err := a.getPodForCR(cr)
	if err != nil {
		return "", err
	}

	return a.operationRunner.Stop(pod, transition, "amf", "amf")
}

func (a Amf) GetCR(transitionCR *fivegv1alpha1.Transition) (common.Entity, error) {
	return a.k8sUtils.GetAmf(types.NamespacedName{
		Namespace: transitionCR.Spec.Config.ResourceNamespace,
		Name:      transitionCR.Spec.Config.ResourceName,
	})
}

func (a Amf) CheckReady(transitionCR *fivegv1alpha1.Transition, entity common.Entity) (bool, error) {
	// AMF is ready when the start transition has completed
	return strings.ToLower(transitionCR.Spec.Config.TransitionName) == "start", nil
}

// Finalize method for Amf.  Executed on uninstall
func (a Amf) Finalize(request reconcile.Request, instance interface{}) {
	// Nothing to do
}

func (a Amf) getCRForTransition(transition *v1alpha1.Transition) (*v1alpha1.Amf, error) {
	cr, err := a.k8sUtils.GetAmf(types.NamespacedName{
		Namespace: transition.Spec.Config.ResourceNamespace,
		Name:      transition.Spec.Config.ResourceName,
	})
	if err != nil {
		return nil, err
	}

	return cr, nil
}

func (a Amf) getPodForCR(amfCR *v1alpha1.Amf) (*v12.Pod, error) {
	pod, err := a.k8sUtils.GetPod(types.NamespacedName{
		Namespace: amfCR.Namespace,
		Name:      amfCR.Name,
	})
	if err != nil {
		return nil, err
	}

	return pod, nil
}
