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

// AusfConfig config parameters supplied via AMF CR
type AusfConfig struct {
	PodSettings  *PodSettingsSpec `json:"podSettings,omitempty"`
	ImageUrl     string           `json:"image,omitempty"`
	ImageUrlInit string           `json:"image_init,omitempty"`
	Name         string           `json:"ausf_name,omitempty"`
	NrfIPAddress string           `json:"nrf_ip_address,omitempty"`
	NrfPort      string           `json:"nrf_port,omitempty"`
}

// AusfInternal config settings derived from other services
type AusfInternal struct {
	Version *string `json:"version,omitempty"`
}

// AusfSpec defines the desired state of Ausf
type AusfSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Config   *AusfConfig   `json:"config,omitempty"`
	Internal *AusfInternal `json:"internal,omitempty"`
}

type AusfOutputs struct {
	IpAddr  string `json:"ip_address,omitempty"`
	PodName string `json:"podname,omitempty"`
}

// AusfStatus defines the observed state of Ausf
type AusfStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	CRStatus `json:",inline"`
	Outputs  AusfOutputs `json:"outputs,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Ausf is the Schema for the ausfs API
type Ausf struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AusfSpec   `json:"spec,omitempty"`
	Status AusfStatus `json:"status,omitempty"`
}

func (a Ausf) GetStatus() *CRStatus {
	return &a.Status.CRStatus
}

func (a *Ausf) GetObject() client.Object {
	return a
}

func (a *Ausf) IsReady() bool {
	return a.Status.IsReady()
}

func (a *Ausf) IsReconciled() bool {
	return a.Status.IsReconciled()
}

func (a *Ausf) SetError(err error) {
	a.GetStatus().SetError(err)
}

func (a *Ausf) SetReady(ready bool) {
	a.GetStatus().SetReady(ready)
}

func (a *Ausf) SetReconciled(reconciled bool) {
	a.GetStatus().SetReconciled(reconciled)
}

//+kubebuilder:object:root=true

// AusfList contains a list of Ausf
type AusfList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Ausf `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Ausf{}, &AusfList{})
}
