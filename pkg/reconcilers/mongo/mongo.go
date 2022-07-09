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

package mongo

import (
	"context"
	"fmt"
	"github.ibm.com/Steve-Glover/5GOperator/api/v1alpha1"
	fivegv1alpha1 "github.ibm.com/Steve-Glover/5GOperator/api/v1alpha1"
	"github.ibm.com/Steve-Glover/5GOperator/pkg/common"
	v12 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"strings"
)

var log = logf.Log.WithName("mongo")

// Mongo is a type to manage k8s objects for Mongo 5G
type Mongo struct {
	*common.EntityHandlerImpl
	operationRunner common.OperationRunner
	k8sUtils        *common.K8sUtils
	client.Client
}

func NewMongo(operationRunner common.OperationRunner, k8sUtils *common.K8sUtils, client client.Client) *Mongo {
	mongo := Mongo{
		EntityHandlerImpl: common.NewEntityHandler(),
		operationRunner:   operationRunner,
		k8sUtils:          k8sUtils,
		Client:            client,
	}

	mongo.AddTransitionFunction("start", mongo.Start)
	mongo.AddTransitionFunction("stop", mongo.Stop)

	return &mongo
}

func (a Mongo) Create(transition *v1alpha1.Transition) (common.Entity, error) {
	desiredInstance := &fivegv1alpha1.Mongo{
		ObjectMeta: metav1.ObjectMeta{
			Name:      transition.Spec.Config.ResourceName,
			Namespace: transition.Spec.Config.ResourceNamespace,
		},
		Spec: fivegv1alpha1.MongoSpec{
			Config: &fivegv1alpha1.MongoConfig{
				PodSettings: common.GetPodSettings(transition),
				ImageUrl:    transition.Spec.Config.Properties["image"],
			},
		},
	}

	instance := &fivegv1alpha1.Mongo{}
	err := a.Client.Get(context.TODO(), types.NamespacedName{
		Namespace: transition.Spec.Config.ResourceNamespace,
		Name:      transition.Spec.Config.ResourceName,
	}, instance)

	created := false
	if err != nil {
		if k8serrors.IsNotFound(err) {
			err = a.Client.Create(context.TODO(), desiredInstance)
			if err != nil {
				transition.Status.SetError(err)
				return nil, err
			}
			created = true
		}

		// Error reading the object - requeue the request.
		return nil, err
	}

	if !created {
		err = a.Client.Update(context.TODO(), desiredInstance)
		if err != nil {
			transition.Status.SetError(err)
			return nil, err
		}
	}

	return instance, nil
}

func (a Mongo) IsReconciled(transitionCR *fivegv1alpha1.Transition) (bool, error) {
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

func (a Mongo) Start(transition *v1alpha1.Transition) (string, error) {
	cr, err := a.getCRForTransition(transition)
	if err != nil {
		return "", err
	}

	pod, err := a.getPodForCR(cr)
	if err != nil {
		return "", err
	}

	return a.operationRunner.Start(pod, transition, "mongo", "mongodb")
}

func (a Mongo) Stop(transition *v1alpha1.Transition) (string, error) {
	cr, err := a.getCRForTransition(transition)
	if err != nil {
		return "", err
	}

	pod, err := a.getPodForCR(cr)
	if err != nil {
		return "", err
	}

	return a.operationRunner.Stop(pod, transition, "mongo", "mongodb")
}

// UpdateStatus called to update CR status
func (a Mongo) UpdateStatus(opResult controllerutil.OperationResult, instance common.Entity, reconcileErr error) {
	cr := instance.(*v1alpha1.Mongo)
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

func (a Mongo) CheckReady(transitionCR *fivegv1alpha1.Transition, entity common.Entity) (bool, error) {
	return strings.ToLower(transitionCR.Spec.Config.TransitionName) == "start", nil
}

func (a Mongo) GetCR(transitionCR *fivegv1alpha1.Transition) (common.Entity, error) {
	return a.k8sUtils.GetMongo(types.NamespacedName{
		Namespace: transitionCR.Spec.Config.ResourceNamespace,
		Name:      transitionCR.Spec.Config.ResourceName,
	})
}

// Finalize method for Mongo.  Executed on uninstall
func (a Mongo) Finalize(request reconcile.Request, instance interface{}) {
	// Nothing to do
}

func (a Mongo) getCRForTransition(transition *v1alpha1.Transition) (*v1alpha1.Mongo, error) {
	cr, err := a.k8sUtils.GetMongo(types.NamespacedName{
		Namespace: transition.Spec.Config.ResourceNamespace,
		Name:      transition.Spec.Config.ResourceName,
	})
	if err != nil {
		return nil, err
	}

	return cr, nil
}

func (a Mongo) getPodForCR(amfCR *v1alpha1.Mongo) (*v12.Pod, error) {
	pod, err := a.k8sUtils.GetPod(types.NamespacedName{
		Namespace: amfCR.Namespace,
		Name:      amfCR.Name,
	})
	if err != nil {
		return nil, err
	}

	return pod, nil
}
