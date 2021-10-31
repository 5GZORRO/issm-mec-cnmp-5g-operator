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

// UdmConfig config parameters supplied via Udm CR
type UdmConfig struct {
	PodSettings  *PodSettingsSpec `json:"podSettings,omitempty"`
	ImageUrl     string           `json:"image,omitempty"`
	ImageUrlInit string           `json:"image_init,omitempty"`	
	NrfIPAddress string           `json:"nrf_ip_address,omitempty"`
	Name         string           `json:"udm_name,omitempty"`
}

// UdmInternal config settings derived from other services
type UdmInternal struct {
	Version *string `json:"version,omitempty"`
}

// UdmSpec defines the desired state of Udm
type UdmSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Config   *UdmConfig   `json:"config,omitempty"`
	Internal *UdmInternal `json:"internal,omitempty"`
}

type UdmOutputs struct {
	IpAddr  string `json:"ip_address,omitempty"`
	PodName string `json:"podname,omitempty"`
}

// UdmStatus defines the observed state of Udm
type UdmStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	CRStatus `json:",inline"`
	Outputs  UdmOutputs `json:"outputs,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Udm is the Schema for the udms API
type Udm struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   UdmSpec   `json:"spec,omitempty"`
	Status UdmStatus `json:"status,omitempty"`
}

func (a Udm) GetStatus() *CRStatus {
	return &a.Status.CRStatus
}

func (a *Udm) GetObject() client.Object {
	return a
}

func (a *Udm) IsReady() bool {
	return a.Status.IsReady()
}

func (a *Udm) IsReconciled() bool {
	return a.Status.IsReconciled()
}

func (a *Udm) SetError(err error) {
	a.GetStatus().SetError(err)
}

func (a *Udm) SetReady(ready bool) {
	a.GetStatus().SetReady(ready)
}

func (a *Udm) SetReconciled(reconciled bool) {
	a.GetStatus().SetReconciled(reconciled)
}

//+kubebuilder:object:root=true

// UdmList contains a list of Udm
type UdmList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Udm `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Udm{}, &UdmList{})
}
