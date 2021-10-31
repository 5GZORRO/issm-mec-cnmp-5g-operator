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

// NrfInternal config settings derived from other services
type NrfInternal struct {
	Version *string `json:"version,omitempty"`
}

// NrfConfig config parameters supplied via AMF CR
type NrfConfig struct {
	PodSettings *PodSettingsSpec `json:"podSettings,omitempty"`
	ImageUrl    string           `json:"image,omitempty"`
	ImageUrlInit string           `json:"image_init,omitempty"`
	MongoIPAddr string           `json:"mongo_ip_address,omitempty"`
	Mnc         string           `json:"mnc,omitempty"`
	Mcc         string           `json:"mcc,omitempty"`
	Port        string           `json:"port,omitempty"`
	Name        string           `json:"nrf_name,omitempty"`
}

// NrfSpec defines the desired state of Nrf
type NrfSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Config   *NrfConfig   `json:"config,omitempty"`
	Internal *NrfInternal `json:"internal,omitempty"`
}

type NrfOutputs struct {
	IpAddr  string `json:"ip_address,omitempty"`
	PodName string `json:"podname,omitempty"`
}

// NrfStatus defines the observed state of Nrf
type NrfStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	CRStatus `json:",inline"`
	Outputs  NrfOutputs `json:"outputs,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Nrf is the Schema for the nrves API
type Nrf struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NrfSpec   `json:"spec,omitempty"`
	Status NrfStatus `json:"status,omitempty"`
}

func (a Nrf) GetStatus() *CRStatus {
	return &a.Status.CRStatus
}

func (a *Nrf) GetObject() client.Object {
	return a
}

func (a *Nrf) IsReady() bool {
	return a.Status.IsReady()
}

func (a *Nrf) IsReconciled() bool {
	return a.Status.IsReconciled()
}

func (a *Nrf) SetError(err error) {
	a.GetStatus().SetError(err)
}

func (a *Nrf) SetReady(ready bool) {
	a.GetStatus().SetReady(ready)
}

func (a *Nrf) SetReconciled(reconciled bool) {
	a.GetStatus().SetReconciled(reconciled)
}

//+kubebuilder:object:root=true

// NrfList contains a list of Nrf
type NrfList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Nrf `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Nrf{}, &NrfList{})
}
