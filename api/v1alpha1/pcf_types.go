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

// PcfConfig config parameters supplied via Pcf CR
type PcfConfig struct {
	PodSettings    *PodSettingsSpec `json:"podSettings,omitempty"`
	ImageUrl       string           `json:"image,omitempty"`
	ImageUrlInit string           `json:"image_init,omitempty"`
	Name           string           `json:"pcf_name,omitempty"`
	NrfIPAddress   string           `json:"nrf_ip_address,omitempty"`
	NrfPort        string           `json:"nrf_port,omitempty"`
	MongoIPAddress string           `json:"mongo_ip_address,omitempty"`
}

// PcfInternal config settings derived from other services
type PcfInternal struct {
	Version *string `json:"version,omitempty"`
}

// PcfSpec defines the desired state of Pcf
type PcfSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Config   *PcfConfig   `json:"config,omitempty"`
	Internal *PcfInternal `json:"internal,omitempty"`
}

type PcfOutputs struct {
	IpAddr  string `json:"ip_address,omitempty"`
	PodName string `json:"podname,omitempty"`
}

// PcfStatus defines the observed state of Pcf
type PcfStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	CRStatus `json:",inline"`
	Outputs  PcfOutputs `json:"outputs,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Pcf is the Schema for the pcfs API
type Pcf struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PcfSpec   `json:"spec,omitempty"`
	Status PcfStatus `json:"status,omitempty"`
}

func (a Pcf) GetStatus() *CRStatus {
	return &a.Status.CRStatus
}

func (a *Pcf) GetObject() client.Object {
	return a
}

func (a *Pcf) IsReady() bool {
	return a.Status.IsReady()
}

func (a *Pcf) IsReconciled() bool {
	return a.Status.IsReconciled()
}

func (a *Pcf) SetError(err error) {
	a.GetStatus().SetError(err)
}

func (a *Pcf) SetReady(ready bool) {
	a.GetStatus().SetReady(ready)
}

func (a *Pcf) SetReconciled(reconciled bool) {
	a.GetStatus().SetReconciled(reconciled)
}

//+kubebuilder:object:root=true

// PcfList contains a list of Pcf
type PcfList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Pcf `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Pcf{}, &PcfList{})
}
