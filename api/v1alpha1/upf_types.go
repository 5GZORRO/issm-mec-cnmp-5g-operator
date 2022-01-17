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


// UpfSpec defines the desired state of Upf
type Dnn struct {
	DnnName     string           `json:"dnn_name,omitempty"`
	ApnCIDR     string           `json:"apn_cidr,omitempty"`
}

// UpfConfig config parameters supplied via Upf CR
type UpfConfig struct {
	PodSettings *PodSettingsSpec `json:"podSettings,omitempty"`
	ImageUrl    string           `json:"image,omitempty"`
	ImageUrlInit string          `json:"image_init,omitempty"`
	Dnns        []Dnn            `json:"dnns,omitempty"`
	Name        string           `json:"upf_name,omitempty"`
	DataNetworkName string       `json:"data_network_name,omitempty"`
	Elicensing *Elicensing   	 `json:"elicensing,omitempty"`
}

// Elicensing defines the parameters required to request POP check 
type Elicensing struct {
	ElmaIp				string	`json:"elma_ip,omitempty"`
	Image				string	`json:"image,omitempty"`		
	ProductOfferingId	string	`json:"product_offering_id,omitempty"`
}

// UpfInternal config settings derived from other services
type UpfInternal struct {
	Version *string `json:"version,omitempty"`
}

// UpfSpec defines the desired state of Upf
type UpfSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Config   *UpfConfig   `json:"config,omitempty"`
	Internal *UpfInternal `json:"internal,omitempty"`
}

type UpfOutputs struct {
	IpAddr  string `json:"ip_address,omitempty"`
	PodName string `json:"podname,omitempty"`
}

// UpfStatus defines the observed state of Upf
type UpfStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	CRStatus `json:",inline"`
	Outputs  UpfOutputs `json:"outputs,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Upf is the Schema for the upfs API
type Upf struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   UpfSpec   `json:"spec,omitempty"`
	Status UpfStatus `json:"status,omitempty"`
}

func (a Upf) GetStatus() *CRStatus {
	return &a.Status.CRStatus
}

func (a *Upf) GetObject() client.Object {
	return a
}

func (a *Upf) IsReady() bool {
	return a.Status.IsReady()
}

func (a *Upf) IsReconciled() bool {
	return a.Status.IsReconciled()
}

func (a *Upf) SetError(err error) {
	a.GetStatus().SetError(err)
}

func (a *Upf) SetReady(ready bool) {
	a.GetStatus().SetReady(ready)
}

func (a *Upf) SetReconciled(reconciled bool) {
	a.GetStatus().SetReconciled(reconciled)
}

//+kubebuilder:object:root=true

// UpfList contains a list of Upf
type UpfList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Upf `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Upf{}, &UpfList{})
}
