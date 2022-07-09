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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

/**
 * Common 5G Core component status handling.
 *
 * The only necessary part is the ObservedGeneration handling which is required by CP4NA to detect
 * changes to the CR. The Operator implementation is free to choose its condition types. In this case,
 * each 5G Core component in this example Kubernetes Operator supports the following conditions:
 *
 * - ready
 * - reconciling
 *
 * The Transition CRD support could be provided as part of an SDK, in which case it would have to support specific
 * conditions so that CP4NA is able to interpret the status properly. For this example controller, the following conditions
 * are supported:
 *
 * - ready
 * - running
 * - reconciling
 */

const (
	Unknown      CRState = "NotFound"
	CRNotFound   CRState = "CRNotFound"
	CRFound      CRState = "CRFound"
	CRReconciled CRState = "CRReconciled"
	CRReady      CRState = "CRReady"
)

type CRState string

type CRStatus struct {
	ObservedGeneration int64               `json:"observedGeneration,omitempty"`
	Conditions         []*metav1.Condition `json:"conditions,omitempty"`
}

func (status *CRStatus) SetReady(ready bool) {
	condition := status.getReadyCondition()
	if ready {
		condition.Status = metav1.ConditionTrue
		condition.Reason = "ready"
	} else {
		condition.Status = metav1.ConditionFalse
		condition.Reason = "not_ready"
	}
}

func (status *CRStatus) SetRunning(isRunning bool) {
	condition := status.getRunningCondition()
	if isRunning {
		condition.Status = metav1.ConditionTrue
		condition.Reason = "running"
	} else {
		condition.Status = metav1.ConditionFalse
		condition.Reason = "not_running"
	}
}

func (status *CRStatus) SetReconciled(reconciled bool) {
	condition := status.getReconciledCondition()
	if reconciled {
		condition.Status = metav1.ConditionTrue
		condition.Reason = "reconciled"
	} else {
		condition.Status = metav1.ConditionFalse
		condition.Reason = "reconciling"
	}
}

func (status *CRStatus) SetError(err error) {
	if err == nil {
		return
	}
	readyCondition := status.getReadyCondition()
	readyCondition.Status = metav1.ConditionFalse
	readyCondition.Message = err.Error()
}

func (status CRStatus) IsReconciled() bool {
	return status.getReconciledCondition().Status == metav1.ConditionTrue
}

func (status CRStatus) IsReady() bool {
	return status.getReadyCondition().Status == metav1.ConditionTrue
}

func (status CRStatus) IsRunning() bool {
	return status.getRunningCondition().Status == metav1.ConditionTrue
}

func (status *CRStatus) IncGeneration(metadata metav1.ObjectMeta) {
	if status.ObservedGeneration < metadata.Generation {
		status.ObservedGeneration = metadata.Generation + 1
	}
}

func (status *CRStatus) Update(metadata metav1.ObjectMeta, isReconciled bool, opResult controllerutil.OperationResult) {
	log.Info("CRStatus Update1",
		"metadata", metadata,
		"opResult", opResult,
		"status", status,
		"isReconciled", isReconciled)

	reconciledCondition := status.getReconciledCondition()
	if isReconciled {
		reconciledCondition.Status = metav1.ConditionTrue
		reconciledCondition.Reason = "reconciled"
	} else {
		reconciledCondition.Status = metav1.ConditionFalse
		reconciledCondition.Reason = "not_reconciled"
	}

	log.Info("CRStatus Update2",
		"metadata", metadata,
		"opResult", opResult,
		"status", status,
		"isReconciled", isReconciled,
		"reconciledCondition", reconciledCondition)
}

func (status *CRStatus) getReadyCondition() *metav1.Condition {
	_, condition := status.getCondition("Ready")
	if condition == nil {
		condition = &metav1.Condition{
			Type:               "Ready",
			Status:             metav1.ConditionFalse,
			LastTransitionTime: metav1.Now(),
			Reason:             "not_ready",
			Message:            "",
		}
		status.Conditions = append(status.Conditions, condition)
	}
	return condition
}

func (status *CRStatus) getRunningCondition() *metav1.Condition {
	_, condition := status.getCondition("Running")
	if condition == nil {
		condition = &metav1.Condition{
			Type:               "Running",
			Status:             metav1.ConditionFalse,
			LastTransitionTime: metav1.Now(),
			Reason:             "not_running",
			Message:            "",
		}
		status.Conditions = append(status.Conditions, condition)
	}
	return condition
}

func (status *CRStatus) getReconciledCondition() *metav1.Condition {
	_, condition := status.getCondition("Reconciled")
	if condition == nil {
		condition = &metav1.Condition{
			Type:               "Reconciled",
			Status:             metav1.ConditionFalse,
			LastTransitionTime: metav1.Now(),
			Reason:             "not_reconciled",
			Message:            "",
		}
		status.Conditions = append(status.Conditions, condition)
	}
	return condition
}

func (status CRStatus) getCondition(conditionType string) (int, *metav1.Condition) {
	for i, _ := range status.Conditions {
		if status.Conditions[i].Type == conditionType {
			return i, status.Conditions[i]
		}
	}
	return -1, nil
}
