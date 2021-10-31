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

// N3iwfConfig config parameters supplied via N3iwf CR
type N3iwfConfig struct {
	PodSettings  *PodSettingsSpec `json:"podSettings,omitempty"`
	ImageUrl     string           `json:"image,omitempty"`
	ImageUrlInit string           `json:"image_init,omitempty"`
	NrfIPAddress string           `json:"nrf_ip_address,omitempty"`
	NrfPort      string           `json:"nrf_port,omitempty"`
	Mnc          string           `json:"mnc,omitempty"`
	Mcc          string           `json:"mcc,omitempty"`
	AmfIPAddress string           `json:"amf_ip_address,omitempty"`
	IPSecAddress string           `json:"ipsec_address,omitempty"`
	//IPAddress string           	  `json:"ip_address,omitempty"`
	UECIDR string `json:"ue_cidr,omitempty"`
}

// N3iwfInternal config settings derived from other services
type N3iwfInternal struct {
	Version *string `json:"version,omitempty"`
}

// N3iwfSpec defines the desired state of N3iwf
type N3iwfSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Config   *N3iwfConfig   `json:"config,omitempty"`
	Internal *N3iwfInternal `json:"internal,omitempty"`
}

type N3iwfOutputs struct {
	IpAddr  string `json:"ip_address,omitempty"`
	PodName string `json:"podname,omitempty"`
}

// N3iwfStatus defines the observed state of N3iwf
type N3iwfStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	CRStatus `json:",inline"`
	Outputs  N3iwfOutputs `json:"outputs,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// N3iwf is the Schema for the n3iwfs API
type N3iwf struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   N3iwfSpec   `json:"spec,omitempty"`
	Status N3iwfStatus `json:"status,omitempty"`
}

func (a N3iwf) GetStatus() *CRStatus {
	return &a.Status.CRStatus
}

func (a *N3iwf) GetObject() client.Object {
	return a
}

func (a *N3iwf) IsReady() bool {
	return a.Status.IsReady()
}

func (a *N3iwf) IsReconciled() bool {
	return a.Status.IsReconciled()
}

func (a *N3iwf) SetError(err error) {
	a.GetStatus().SetError(err)
}

func (a *N3iwf) SetReady(ready bool) {
	a.GetStatus().SetReady(ready)
}

func (a *N3iwf) SetReconciled(reconciled bool) {
	a.GetStatus().SetReconciled(reconciled)
}

//+kubebuilder:object:root=true

// N3iwfList contains a list of N3iwf
type N3iwfList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []N3iwf `json:"items"`
}

func init() {
	SchemeBuilder.Register(&N3iwf{}, &N3iwfList{})
}
