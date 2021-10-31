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
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// MongoConfig config parameters supplied via AMF CR
type MongoConfig struct {
	PodSettings *PodSettingsSpec `json:"podSettings,omitempty"`
	ImageUrl    string           `json:"image,omitempty"`
	ImageUrlInit string           `json:"image_init,omitempty"`
}

// MongoSpec defines the desired state of Mongo
type MongoSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Config *MongoConfig `json:"config,omitempty"`
}

type MongoOutputs struct {
	IpAddr  string `json:"ip_address,omitempty"`
	PodName string `json:"podname,omitempty"`
}

// MongoStatus defines the observed state of Mongo
type MongoStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	CRStatus `json:",inline"`
	Outputs  MongoOutputs `json:"outputs,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Mongo is the Schema for the mongoes API
type Mongo struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MongoSpec   `json:"spec,omitempty"`
	Status MongoStatus `json:"status,omitempty"`
}

func (m Mongo) GetStatus() *MongoStatus {
	return &m.Status
}

func (m *Mongo) GetObject() client.Object {
	return m
}

func (a *Mongo) IsReady() bool {
	return a.Status.IsReady()
}

func (a *Mongo) IsReconciled() bool {
	return a.Status.IsReconciled()
}

func (a *Mongo) SetError(err error) {
	a.GetStatus().SetError(err)
}

func (a *Mongo) SetReady(ready bool) {
	(&(a.Status)).SetReady(ready)
	log.Info("Mongo SetReady", "mongo", a)
}

func (a *Mongo) SetReconciled(reconciled bool) {
	a.GetStatus().SetReconciled(reconciled)
}

//+kubebuilder:object:root=true

// MongoList contains a list of Mongo
type MongoList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Mongo `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Mongo{}, &MongoList{})
}
