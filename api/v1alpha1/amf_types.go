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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// Global configurations
type Global struct {
	ImagePullPolicy *corev1.PullPolicy `json:"imagePullPolicy,omitempty"`
}

// AmfConfig config parameters supplied via AMF CR
type AmfConfig struct {
	PodSettings  *PodSettingsSpec `json:"podSettings,omitempty"`
	ImageUrl     string           `json:"image,omitempty"`
	ImageUrlInit string           `json:"image_init,omitempty"`
	Name         string           `json:"amf_name,omitempty"`
	NrfIPAddress string           `json:"nrf_ip_address,omitempty"`
	NrfPort      string           `json:"nrf_port,omitempty"`
	Mnc          string           `json:"mnc,omitempty"`
	Mcc          string           `json:"mcc,omitempty"`
}

// AmfInternal config settings derived from other services
type AmfInternal struct {
	Version *string `json:"version,omitempty"`
}

// AmfSpec defines the desired state of Amf
//+kubebuilder:subresource:status
type AmfSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	//IPAddress string `json:"ip_address,omitempty"`
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Config   *AmfConfig   `json:"config,omitempty"`
	Internal *AmfInternal `json:"internal,omitempty"`
}

type AmfOutputs struct {
	IpAddr  string `json:"ip_address,omitempty"`
	PodName string `json:"podname,omitempty"`
}

// AmfStatus defines the observed state of Amf
type AmfStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	CRStatus `json:",inline"`
	Outputs  AmfOutputs `json:"outputs,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Amf is the Schema for the amfs API
type Amf struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AmfSpec   `json:"spec,omitempty"`
	Status AmfStatus `json:"status,omitempty"`
}

func (a Amf) GetStatus() *CRStatus {
	return &a.Status.CRStatus
}

func (a *Amf) GetObject() client.Object {
	return a
}

func (a *Amf) IsReady() bool {
	return a.Status.IsReady()
}

func (a *Amf) IsReconciled() bool {
	return a.Status.IsReconciled()
}

func (a *Amf) SetError(err error) {
	a.GetStatus().SetError(err)
}

func (a *Amf) SetReady(ready bool) {
	a.GetStatus().SetReady(ready)
}

func (a *Amf) SetReconciled(reconciled bool) {
	a.GetStatus().SetReconciled(reconciled)
}

//+kubebuilder:object:root=true

// AmfList contains a list of Amf
type AmfList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Amf `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Amf{}, &AmfList{})
}
