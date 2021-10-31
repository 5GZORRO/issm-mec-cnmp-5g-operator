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

// NssfConfig config parameters supplied via Nssf CR
type NssfConfig struct {
	PodSettings  *PodSettingsSpec `json:"podSettings,omitempty"`
	ImageUrl     string           `json:"image,omitempty"`
	ImageUrlInit string           `json:"image_init,omitempty"`
	NrfIPAddress string           `json:"nrf_ip_address,omitempty"`
	NrfPort      string           `json:"nrf_port,omitempty"`
	Mnc          string           `json:"mnc,omitempty"`
	Mcc          string           `json:"mcc,omitempty"`
    Name         string           `json:"nssf_name,omitempty"`
}

// NssfInternal config settings derived from other services
type NssfInternal struct {
	Version *string `json:"version,omitempty"`
}

// NssfSpec defines the desired state of Nssf
type NssfSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Config   *NssfConfig   `json:"config,omitempty"`
	Internal *NssfInternal `json:"internal,omitempty"`
}

type NssfOutputs struct {
	IpAddr  string `json:"ip_address,omitempty"`
	PodName string `json:"podname,omitempty"`
}

// NssfStatus defines the observed state of Nssf
type NssfStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	CRStatus `json:",inline"`
	Outputs  NssfOutputs `json:"outputs,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Nssf is the Schema for the nssfs API
type Nssf struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NssfSpec   `json:"spec,omitempty"`
	Status NssfStatus `json:"status,omitempty"`
}

func (a Nssf) GetStatus() *CRStatus {
	return &a.Status.CRStatus
}

func (a *Nssf) GetObject() client.Object {
	return a
}

func (a *Nssf) IsReady() bool {
	return a.Status.IsReady()
}

func (a *Nssf) IsReconciled() bool {
	return a.Status.IsReconciled()
}

func (a *Nssf) SetError(err error) {
	a.GetStatus().SetError(err)
}

func (a *Nssf) SetReady(ready bool) {
	a.GetStatus().SetReady(ready)
}

func (a *Nssf) SetReconciled(reconciled bool) {
	a.GetStatus().SetReconciled(reconciled)
}

//+kubebuilder:object:root=true

// NssfList contains a list of Nssf
type NssfList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Nssf `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Nssf{}, &NssfList{})
}
