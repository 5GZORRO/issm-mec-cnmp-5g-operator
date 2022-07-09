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
	"github.com/5GZORRO/issm-mec-cnmp-5g-operator/api/model"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// TransitionConfig config parameters supplied via Transition CR
type TransitionConfig struct {
	ResourceType      string           `json:"resourceType,omitempty"`
	ResourceNamespace string           `json:"resourceNamespace,omitempty"`
	ResourceName      string           `json:"resourceName,omitempty"`
	TransitionName    string           `json:"transition,omitempty"`
	Timestamp         string           `json:"ts,omitempty"`
	Properties        model.Properties `json:"properties,omitempty"`
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// TransitionSpec defines the desired state of Transition
type TransitionSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Config *TransitionConfig `json:"config,omitempty"`
}

// TransitionStatus defines the observed state of Transition
type TransitionStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	CRStatus `json:",inline"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Transition is the Schema for the transitions API
type Transition struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TransitionSpec   `json:"spec,omitempty"`
	Status TransitionStatus `json:"status,omitempty"`
}

func (t Transition) GetStatus() *TransitionStatus {
	return &t.Status
}

func (t *Transition) GetObject() client.Object {
	return t
}

func (t *Transition) IsReady() bool {
	return t.Status.IsReady()
}

func (t *Transition) IsReconciled() bool {
	return t.Status.IsReconciled()
}

func (t *Transition) SetError(err error) {
	t.GetStatus().SetError(err)
}

func (t *Transition) SetReady(ready bool) {
	t.GetStatus().SetReady(ready)
}

func (t *Transition) SetReconciled(reconciled bool) {
	t.GetStatus().SetReconciled(reconciled)
}

//+kubebuilder:object:root=true

// TransitionList contains a list of Transition
type TransitionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Transition `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Transition{}, &TransitionList{})
}
