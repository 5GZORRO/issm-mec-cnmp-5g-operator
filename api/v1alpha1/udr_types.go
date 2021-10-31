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

// UdrConfig config parameters supplied via Udr CR
type UdrConfig struct {
	PodSettings    *PodSettingsSpec `json:"podSettings,omitempty"`
	ImageUrl       string           `json:"image,omitempty"`
	ImageUrlInit   string           `json:"image_init,omitempty"`	
	Name           string           `json:"udr_name,omitempty"`
	NrfIPAddress   string           `json:"nrf_ip_address,omitempty"`
	NrfPort        string           `json:"nrf_port,omitempty"`
	MongoIPAddress string           `json:"mongo_ip_address,omitempty"`
}

// UdrInternal config settings derived from other services
type UdrInternal struct {
	Version *string `json:"version,omitempty"`
}

// UdrSpec defines the desired state of Udr
type UdrSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Config   *UdrConfig   `json:"config,omitempty"`
	Internal *UdrInternal `json:"internal,omitempty"`
}

type UdrOutputs struct {
	IpAddr  string `json:"ip_address,omitempty"`
	PodName string `json:"podname,omitempty"`
}

// UdrStatus defines the observed state of Udr
type UdrStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	CRStatus `json:",inline"`
	Outputs  UdrOutputs `json:"outputs,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Udr is the Schema for the udrs API
type Udr struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   UdrSpec   `json:"spec,omitempty"`
	Status UdrStatus `json:"status,omitempty"`
}

func (a Udr) GetStatus() *CRStatus {
	return &a.Status.CRStatus
}

func (a *Udr) GetObject() client.Object {
	return a
}

func (a *Udr) IsReady() bool {
	return a.Status.IsReady()
}

func (a *Udr) IsReconciled() bool {
	return a.Status.IsReconciled()
}

func (a *Udr) SetError(err error) {
	a.GetStatus().SetError(err)
}

func (a *Udr) SetReady(ready bool) {
	a.GetStatus().SetReady(ready)
}

func (a *Udr) SetReconciled(reconciled bool) {
	a.GetStatus().SetReconciled(reconciled)
}

//+kubebuilder:object:root=true

// UdrList contains a list of Udr
type UdrList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Udr `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Udr{}, &UdrList{})
}
