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

package smf

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

// Smf is a type to manage k8s objects for Smf 5G
type Smf struct {
	*common.EntityHandlerImpl
	operationRunner common.OperationRunner
	k8sUtils        *common.K8sUtils
	client.Client
}

func NewSmf(operationRunner common.OperationRunner, k8sUtils *common.K8sUtils, client client.Client) *Smf {
	smf := Smf{
		EntityHandlerImpl: common.NewEntityHandler(),
		operationRunner:   operationRunner,
		k8sUtils:          k8sUtils,
		Client:            client,
	}

	smf.EntityHandlerImpl.AddTransitionFunction("start", smf.Start)
	smf.EntityHandlerImpl.AddTransitionFunction("stop", smf.Stop)
	smf.EntityHandlerImpl.AddTransitionFunction("attachnode", smf.AddNode)
	smf.EntityHandlerImpl.AddTransitionFunction("detachnode", smf.RmNode)
	smf.EntityHandlerImpl.AddTransitionFunction("addlink", smf.AddLink)
	smf.EntityHandlerImpl.AddTransitionFunction("rmlink", smf.RmLink)
	smf.EntityHandlerImpl.AddTransitionFunction("addueroute", smf.AddUERoute)

	return &smf
}

// UpdateStatus called to update CR status
func (a Smf) UpdateStatus(opResult controllerutil.OperationResult, instance common.Entity, reconcileErr error) {
	cr := instance.(*v1alpha1.Smf)

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

// Finalize method for Smf.  Executed on uninstall
func (a Smf) Finalize(request reconcile.Request, instance interface{}) {
	// Nothing to do
}

func (a Smf) Create(transition *v1alpha1.Transition) (common.Entity, error) {
	instance := &fivegv1alpha1.Smf{
		ObjectMeta: metav1.ObjectMeta{
			Name:      transition.Spec.Config.ResourceName,
			Namespace: transition.Spec.Config.ResourceNamespace,
		},
		Spec: fivegv1alpha1.SmfSpec{
			Config: &fivegv1alpha1.SmfConfig{
				PodSettings:  common.GetPodSettings(transition),
				ImageUrl:     transition.Spec.Config.Properties["image"],
				NrfIPAddress: transition.Spec.Config.Properties["nrf_ip_address"],
				NrfPort:      transition.Spec.Config.Properties["nrf_port"],
				Nodes:        make([]fivegv1alpha1.UpNode, 0, 10),
				Links:        make([]fivegv1alpha1.Link, 0, 10),
				UEList:       fivegv1alpha1.UEList{UES: make([]fivegv1alpha1.UE, 0, 10)},
			},
		},
	}

	err := a.Client.Create(context.TODO(), instance)
	if err != nil {
		return nil, err
	}

	return instance, nil
}

func (a Smf) Start(transition *v1alpha1.Transition) (string, error) {
	cr, err := a.getCRForTransition(transition)
	if err != nil {
		return "", err
	}

	pod, err := a.getPodForCR(cr)
	if err != nil {
		return "", err
	}

	return a.operationRunner.Start(pod, transition, "smf", "smf")
}

func (a Smf) Stop(transition *v1alpha1.Transition) (string, error) {
	cr, err := a.getCRForTransition(transition)
	if err != nil {
		return "", err
	}

	pod, err := a.getPodForCR(cr)
	if err != nil {
		return "", err
	}

	return a.operationRunner.Stop(pod, transition, "smf", "smf")
}

func (a Smf) AddNode(transition *v1alpha1.Transition) (string, error) {
	cr, err := a.getCRForTransition(transition)
	if err != nil {
		return "", err
	}

	nodeName, ok := transition.Spec.Config.Properties["nodeName"]
	if !ok {
		return "", fmt.Errorf("Missing nodeName")
	}
	nodeType, ok := transition.Spec.Config.Properties["nodeType"]
	if !ok {
		return "", fmt.Errorf("Missing nodeType")
	}
	nodeIdSbi, ok := transition.Spec.Config.Properties["nodeIdSbi"]
	if !ok {
		return "", fmt.Errorf("Missing nodeIdSbi")
	}
	nodeIdUp, ok := transition.Spec.Config.Properties["nodeIdUp"]
	if !ok {
		return "", fmt.Errorf("Missing nodeIdUp")
	}

	sd, ok := transition.Spec.Config.Properties["sd"]
	if !ok {
		return "", fmt.Errorf("Missing sd")
	}
	sst, ok := transition.Spec.Config.Properties["sst"]
	if !ok {
		return "", fmt.Errorf("Missing sst")
	}

	pool, ok := transition.Spec.Config.Properties["pool"]
	if !ok {
		// update the SMF CR, which should restart the pod with the new configuration
		cr.Spec.Config.Nodes = append(cr.Spec.Config.Nodes, v1alpha1.UpNode{
			Name:   nodeName,
			Type:   nodeType,
			NodeIdUp: nodeIdUp,
			NodeIdSbi: nodeIdSbi,
			Sd:     sd,
			Sst:    sst,
		})
	} else {
		// update the SMF CR, which should restart the pod with the new configuration
		cr.Spec.Config.Nodes = append(cr.Spec.Config.Nodes, v1alpha1.UpNode{
			Name:   nodeName,
			Type:   nodeType,
			NodeIdUp: nodeIdUp,
			NodeIdSbi: nodeIdSbi,
			Sd:     sd,
			Sst:    sst,
			Pool:   pool,
		})
	}

	err = a.k8sUtils.UpdateCR(cr)
	return "", err
}

func (a Smf) RmNode(transition *v1alpha1.Transition) (string, error) {
	cr, err := a.getCRForTransition(transition)
	if err != nil {
		return "", err
	}

	nodeName, ok := transition.Spec.Config.Properties["nodeName"]
	if !ok {
		return "", fmt.Errorf("Missing nodeName")
	}

	i := 0
	for range cr.Spec.Config.Nodes {
		log.Info("Name", "cr.Spec.Config.Nodes[i].Name", cr.Spec.Config.Nodes[i].Name)
		if cr.Spec.Config.Nodes[i].Name == nodeName {
			break
		}
		i++
	}
	log.Info("I", "i", i)
	if i < len(cr.Spec.Config.Nodes) {
		// remove
		cr.Spec.Config.Nodes[i] = cr.Spec.Config.Nodes[len(cr.Spec.Config.Nodes)-1] // Copy last element to index i
		cr.Spec.Config.Nodes[len(cr.Spec.Config.Nodes)-1] = v1alpha1.UpNode{} // Erase last element (write zero value)
		cr.Spec.Config.Nodes = cr.Spec.Config.Nodes[:len(cr.Spec.Config.Nodes)-1] // Truncate slice.
	}
	err = a.k8sUtils.UpdateCR(cr)
	return "", err
}


func (a Smf) AddLink(transition *v1alpha1.Transition) (string, error) {
	cr, err := a.getCRForTransition(transition)
	if err != nil {
		return "", err
	}

	nodeId, ok := transition.Spec.Config.Properties["nodeId"]
	if !ok {
		return "", fmt.Errorf("Missing nodeId")
	}
	connectedTo, ok := transition.Spec.Config.Properties["connectedTo"]
	if !ok {
		return "", fmt.Errorf("Missing connectedTo")
	}

	// update the SMF CR, which should restart the pod with the new configuration
	cr.Spec.Config.Links = append(cr.Spec.Config.Links, v1alpha1.Link{
		AEnd:   nodeId,
		BEnd:   connectedTo,
	})

	err = a.k8sUtils.UpdateCR(cr)
	return "", err
}

func (a Smf) RmLink(transition *v1alpha1.Transition) (string, error) {
	// this version assumes single occurance of this UPF in links list
	cr, err := a.getCRForTransition(transition)
	if err != nil {
		return "", err
	}

	nodeId, ok := transition.Spec.Config.Properties["nodeId"]
	if !ok {
		return "", fmt.Errorf("Missing nodeId")
	}

	i := 0
	for range cr.Spec.Config.Links {
		if cr.Spec.Config.Links[i].AEnd == nodeId || cr.Spec.Config.Links[i].BEnd == nodeId {
			break
		}
		i++
	}
	if i < len(cr.Spec.Config.Links) {
		// remove
		cr.Spec.Config.Links[i] = cr.Spec.Config.Links[len(cr.Spec.Config.Links)-1] // Copy last element to index i
		cr.Spec.Config.Links[len(cr.Spec.Config.Links)-1] = v1alpha1.Link{} // Erase last element (write zero value)
		cr.Spec.Config.Links = cr.Spec.Config.Links[:len(cr.Spec.Config.Links)-1] // Truncate slice.
	}
	err = a.k8sUtils.UpdateCR(cr)
	return "", err
}

func (a Smf) AddUERoute(transition *v1alpha1.Transition) (string, error) {
	cr, err := a.getCRForTransition(transition)
	if err != nil {
		return "", err
	}

	SUPI, ok := transition.Spec.Config.Properties["SUPI"]
	if !ok {
		return "", fmt.Errorf("Missing SUPI")
	}
	AN, ok := transition.Spec.Config.Properties["AN"]
	if !ok {
		return "", fmt.Errorf("Missing AN")
	}
	DestinationIP, ok := transition.Spec.Config.Properties["DestinationIP"]
	if !ok {
		return "", fmt.Errorf("Missing DestinationIP")
	}
	DestinationPort, ok := transition.Spec.Config.Properties["DestinationPort"]
	if !ok {
		return "", fmt.Errorf("Missing DestinationPort")
	}
	UPFPathAsString, ok := transition.Spec.Config.Properties["UPFPath"]
	if !ok {
		return "", fmt.Errorf("Missing UPFPath")
	}

	cr.Spec.Config.UEList.UES = append(cr.Spec.Config.UEList.UES, v1alpha1.UE{
		SUPI:            SUPI,
		AN:              AN,
		DestinationIP:   DestinationIP,
		DestinationPort: DestinationPort,
		UPFPath:         parseUPFPath(UPFPathAsString),
	})

	return "", a.k8sUtils.UpdateCR(cr)
}

func (a Smf) CheckReady(transitionCR *fivegv1alpha1.Transition, entity common.Entity) (bool, error) {
	return strings.ToLower(transitionCR.Spec.Config.TransitionName) == "start", nil
}

func (a Smf) GetCR(transitionCR *fivegv1alpha1.Transition) (common.Entity, error) {
	return a.k8sUtils.GetSmf(types.NamespacedName{
		Namespace: transitionCR.Spec.Config.ResourceNamespace,
		Name:      transitionCR.Spec.Config.ResourceName,
	})
}

func (a Smf) IsReconciled(transitionCR *fivegv1alpha1.Transition) (bool, error) {
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

func (a Smf) getCRForTransition(transition *v1alpha1.Transition) (*v1alpha1.Smf, error) {
	cr, err := a.k8sUtils.GetSmf(types.NamespacedName{
		Namespace: transition.Spec.Config.ResourceNamespace,
		Name:      transition.Spec.Config.ResourceName,
	})
	if err != nil {
		return nil, err
	}

	return cr, nil
}

func (a Smf) getPodForCR(amfCR *v1alpha1.Smf) (*v12.Pod, error) {
	pod, err := a.k8sUtils.GetPod(types.NamespacedName{
		Namespace: amfCR.Namespace,
		Name:      amfCR.Name,
	})
	if err != nil {
		return nil, err
	}

	return pod, nil
}

func parseUPFPath(UPFPathAsString string) []string {
	l := strings.Split(UPFPathAsString, ",")
	UPFPath := make([]string, 0, len(l))
	for idx, _ := range l {
		UPFPath = append(UPFPath, l[idx])
	}

	return UPFPath
}
