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

package upf

import (
	"fmt"
	"github.com/5GZORRO/issm-mec-cnmp-5g-operator/api/v1alpha1"
	fivegv1alpha1 "github.com/5GZORRO/issm-mec-cnmp-5g-operator/api/v1alpha1"
	"github.com/5GZORRO/issm-mec-cnmp-5g-operator/pkg/common"
	v12 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"strings"
)

// Upf is a type to manage k8s objects for Upf 5G
type Upf struct {
	*common.EntityHandlerImpl
	operationRunner common.OperationRunner
	k8sUtils        *common.K8sUtils
	client.Client
}

func NewUpf(operationRunner common.OperationRunner, k8sUtils *common.K8sUtils, client client.Client) *Upf {
	upf := Upf{
		EntityHandlerImpl: common.NewEntityHandler(),
		operationRunner:   operationRunner,
		k8sUtils:          k8sUtils,
		Client:            client,
	}

	upf.EntityHandlerImpl.AddTransitionFunction("start", upf.Start)
	upf.EntityHandlerImpl.AddTransitionFunction("stop", upf.Stop)

	return &upf
}

// UpdateStatus called to update CR status
func (a Upf) UpdateStatus(opResult controllerutil.OperationResult, instance common.Entity, reconcileErr error) {
	cr := instance.(*v1alpha1.Upf)

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
			_, _ = a.k8sUtils.ExecCommand(obj, "upf", "iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE")
			// TODO: write this better
			if cr.Spec.Config.DataNetworkName != "" {
				_, _ = a.k8sUtils.ExecCommand(obj, "upf", "iptables -t nat -A POSTROUTING -o net3 -j MASQUERADE")
			}
		}
	}
}

// Finalize method for Upf.  Executed on uninstall
func (a Upf) Finalize(request reconcile.Request, instance interface{}) {
	// Nothing to do
}

func (a Upf) Create(transition *v1alpha1.Transition) (common.Entity, error) {
//	instance := &fivegv1alpha1.Upf{
//		ObjectMeta: metav1.ObjectMeta{
//			Name:      transition.Spec.Config.ResourceName,
//			Namespace: transition.Spec.Config.ResourceNamespace,
//		},
//		Spec: fivegv1alpha1.UpfSpec{
//			Config: &fivegv1alpha1.UpfConfig{
//				PodSettings: common.GetPodSettings(transition),
//				ImageUrl:    transition.Spec.Config.Properties["image"],
//				DnnName:     transition.Spec.Config.Properties["dnn_name"],
//				ApnCIDR:     transition.Spec.Config.Properties["apn_cidr"],
//			},
//		},
//	}
//
//	err := a.Client.Create(context.TODO(), instance)
//	if err != nil {
//		return nil, err
//	}
//
	return nil, nil
}

func (a Upf) Start(transition *v1alpha1.Transition) (string, error) {
	cr, err := a.getCRForTransition(transition)
	if err != nil {
		return "", err
	}

	pod, err := a.getPodForCR(cr)
	if err != nil {
		return "", err
	}

	return a.operationRunner.Start(pod, transition, "upf", "upf")
}

func (a Upf) Stop(transition *v1alpha1.Transition) (string, error) {
	cr, err := a.getCRForTransition(transition)
	if err != nil {
		return "", err
	}

	pod, err := a.getPodForCR(cr)
	if err != nil {
		return "", err
	}

	return a.operationRunner.Stop(pod, transition, "upf", "upf")
}

func (a Upf) CheckReady(transitionCR *fivegv1alpha1.Transition, entity common.Entity) (bool, error) {
	return strings.ToLower(transitionCR.Spec.Config.TransitionName) == "start", nil
}

func (a Upf) GetCR(transitionCR *fivegv1alpha1.Transition) (common.Entity, error) {
	return a.k8sUtils.GetUpf(types.NamespacedName{
		Namespace: transitionCR.Spec.Config.ResourceNamespace,
		Name:      transitionCR.Spec.Config.ResourceName,
	})
}

func (a Upf) IsReconciled(transitionCR *fivegv1alpha1.Transition) (bool, error) {
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

func (a Upf) getCRForTransition(transition *v1alpha1.Transition) (*v1alpha1.Upf, error) {
	cr, err := a.k8sUtils.GetUpf(types.NamespacedName{
		Namespace: transition.Spec.Config.ResourceNamespace,
		Name:      transition.Spec.Config.ResourceName,
	})
	if err != nil {
		return nil, err
	}

	return cr, nil
}

func (a Upf) getPodForCR(cr *v1alpha1.Upf) (*v12.Pod, error) {
	pod, err := a.k8sUtils.GetPod(types.NamespacedName{
		Namespace: cr.Namespace,
		Name:      cr.Name,
	})
	if err != nil {
		return nil, err
	}

	return pod, nil
}
