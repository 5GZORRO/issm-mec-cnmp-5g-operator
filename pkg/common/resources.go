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

package common

import (
	"fmt"
	"github.ibm.com/Steve-Glover/5GOperator/api/v1alpha1"
	fivegv1alpha1 "github.ibm.com/Steve-Glover/5GOperator/api/v1alpha1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"strings"
)

type EntityHandler interface {
	Create(transition *v1alpha1.Transition) (Entity, error)
	GetCR(transitionCR *fivegv1alpha1.Transition) (Entity, error)
	IsReconciled(transitionCR *fivegv1alpha1.Transition) (bool, error)
	IsReady(transitionCR *fivegv1alpha1.Transition) (bool, error)
	RunTransition(transition *v1alpha1.Transition, entity Entity, updateFunc func(Entity, *fivegv1alpha1.Transition, error) error) error
	CheckReady(transitionCR *fivegv1alpha1.Transition, entity Entity) (bool, error)
	AddTransitionFunction(name string, transitionFunc TransitionFunc)
}

type EntityHandlerImpl struct {
	TransitionFuncs map[string]TransitionFunc
}

func NewEntityHandler() *EntityHandlerImpl {
	r := &EntityHandlerImpl{
		TransitionFuncs: make(map[string]TransitionFunc),
	}
	return r
}

func GetCRState(r EntityHandler, transitionCR *fivegv1alpha1.Transition) (fivegv1alpha1.CRState, error) {
	entity, err := r.GetCR(transitionCR)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return fivegv1alpha1.CRNotFound, nil
		}
		return fivegv1alpha1.Unknown, err
	}
	if entity == nil {
		return fivegv1alpha1.CRNotFound, nil
	}

	if entity.IsReady() {
		return fivegv1alpha1.CRReady, nil
	} else if entity.IsReconciled() {
		return fivegv1alpha1.CRReconciled, nil
	} else {
		return fivegv1alpha1.CRFound, nil
	}
}

func (r EntityHandlerImpl) AddTransitionFunction(name string, transitionFunc TransitionFunc) {
	r.TransitionFuncs[name] = transitionFunc
}

func (r EntityHandlerImpl) GetCR(transitionCR *fivegv1alpha1.Transition) (Entity, error) {
	log.Info("GetCR", "transitionCR", transitionCR)
	return nil, nil
}

func (r EntityHandlerImpl) IsReconciled(transitionCR *fivegv1alpha1.Transition) (bool, error) {
	cr, err := r.GetCR(transitionCR)
	if err != nil {
		return false, err
	}
	if cr == nil {
		return false, fmt.Errorf("Cannot find cr")
	}
	return cr.IsReconciled(), nil
}

func (r EntityHandlerImpl) IsReady(transitionCR *fivegv1alpha1.Transition) (bool, error) {
	cr, err := r.GetCR(transitionCR)
	if err != nil {
		return false, err
	}
	if cr == nil {
		err = fmt.Errorf("Cannot find cr")
		return false, err
	}
	return cr.IsReady(), nil
}

func (r EntityHandlerImpl) RunTransition(transition *v1alpha1.Transition, entity Entity, updateFunc func(Entity, *fivegv1alpha1.Transition, error) error) error {
	go func() {
		log.Info("RunTransition", "resourceType", r.Type(), "transition", transition)

		transitionName := strings.ToLower(transition.Spec.Config.TransitionName)
		var err error
		switch transitionName {
		case "create":
			err = fmt.Errorf("Create transition is being handled in the wrong place %+v", transition)
		default:
			transitionFunc, ok := r.TransitionFuncs[transitionName]
			if ok {
				// the transition execution runs synchronously (hence in a Goroutine)
				var output string
				output, err = transitionFunc(transition)
				log.Info("RunTransition complete", "output", output)
			} else {
				err = fmt.Errorf("Invalid transition %s", transitionName)
			}
		}

		// run the callback to update the CR status
		updateFunc(entity, transition, err)
	}()

	return nil
}

func (r EntityHandlerImpl) Create(transition *v1alpha1.Transition) (Entity, error) {
	return nil, nil
}

func (r EntityHandlerImpl) CheckReady(transitionCR *fivegv1alpha1.Transition, entity Entity) error {
	log.Info("CheckReady", "transitionCR", transitionCR, "entity", entity)
	return nil
}

func (r EntityHandlerImpl) Type() string {
	return "generic"
}
