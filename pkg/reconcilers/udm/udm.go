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

package udm

import (
	"context"
	"fmt"
	"github.com/5GZORRO/issm-mec-cnmp-5g-operator/api/v1alpha1"
	fivegv1alpha1 "github.com/5GZORRO/issm-mec-cnmp-5g-operator/api/v1alpha1"
	"github.com/5GZORRO/issm-mec-cnmp-5g-operator/pkg/common"
	v12 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"strings"
)

// Udm is a type to manage k8s objects for Udm 5G
type Udm struct {
	*common.EntityHandlerImpl
	operationRunner common.OperationRunner
	k8sUtils        *common.K8sUtils
	client.Client
}

func NewUdm(operationRunner common.OperationRunner, k8sUtils *common.K8sUtils, client client.Client) *Udm {
	udm := Udm{
		EntityHandlerImpl: common.NewEntityHandler(),
		operationRunner:   operationRunner,
		k8sUtils:          k8sUtils,
		Client:            client,
	}

	udm.EntityHandlerImpl.AddTransitionFunction("start", udm.Start)
	udm.EntityHandlerImpl.AddTransitionFunction("stop", udm.Stop)

	return &udm
}

// UpdateStatus called to update CR status
func (a Udm) UpdateStatus(opResult controllerutil.OperationResult, instance common.Entity, reconcileErr error) {
	cr := instance.(*v1alpha1.Udm)

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

// Finalize method for Udm.  Executed on uninstall
func (a Udm) Finalize(request reconcile.Request, instance interface{}) {
	// Nothing to do
}

func (a Udm) Create(transition *v1alpha1.Transition) (common.Entity, error) {
	instance := &fivegv1alpha1.Udm{
		ObjectMeta: metav1.ObjectMeta{
			Name:      transition.Spec.Config.ResourceName,
			Namespace: transition.Spec.Config.ResourceNamespace,
		},
		Spec: fivegv1alpha1.UdmSpec{
			Config: &fivegv1alpha1.UdmConfig{
				PodSettings:  common.GetPodSettings(transition),
				ImageUrl:     transition.Spec.Config.Properties["image"],
				NrfIPAddress: transition.Spec.Config.Properties["nrf_ip_address"],
			},
		},
	}

	err := a.Client.Create(context.TODO(), instance)
	if err != nil {
		return nil, err
	}

	return instance, nil
}

func (a Udm) Start(transition *v1alpha1.Transition) (string, error) {
	cr, err := a.getCRForTransition(transition)
	if err != nil {
		return "", err
	}

	pod, err := a.getPodForCR(cr)
	if err != nil {
		return "", err
	}

	return a.operationRunner.Start(pod, transition, "udm", "udm")
}

func (a Udm) Stop(transition *v1alpha1.Transition) (string, error) {
	cr, err := a.getCRForTransition(transition)
	if err != nil {
		return "", err
	}

	pod, err := a.getPodForCR(cr)
	if err != nil {
		return "", err
	}

	return a.operationRunner.Stop(pod, transition, "udm", "udm")
}

func (a Udm) CheckReady(transitionCR *fivegv1alpha1.Transition, entity common.Entity) (bool, error) {
	return strings.ToLower(transitionCR.Spec.Config.TransitionName) == "start", nil
}

func (a Udm) GetCR(transitionCR *fivegv1alpha1.Transition) (common.Entity, error) {
	return a.k8sUtils.GetUdm(types.NamespacedName{
		Namespace: transitionCR.Spec.Config.ResourceNamespace,
		Name:      transitionCR.Spec.Config.ResourceName,
	})
}

func (a Udm) IsReconciled(transitionCR *fivegv1alpha1.Transition) (bool, error) {
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

func (a Udm) getCRForTransition(transition *v1alpha1.Transition) (*v1alpha1.Udm, error) {
	cr, err := a.k8sUtils.GetUdm(types.NamespacedName{
		Namespace: transition.Spec.Config.ResourceNamespace,
		Name:      transition.Spec.Config.ResourceName,
	})
	if err != nil {
		return nil, err
	}

	return cr, nil
}

func (a Udm) getPodForCR(cr *v1alpha1.Udm) (*v12.Pod, error) {
	pod, err := a.k8sUtils.GetPod(types.NamespacedName{
		Namespace: cr.Namespace,
		Name:      cr.Name,
	})
	if err != nil {
		return nil, err
	}

	return pod, nil
}
