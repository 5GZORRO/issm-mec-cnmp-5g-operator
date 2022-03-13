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

type UpNode struct {
	Name   string `json:"name,omitempty"`
	Type   string `json:"type,omitempty"`
	AnIp   string `json:"an_ip,omitempty"`
	NodeId string `json:"node_id,omitempty"`
	NodeIdUp string `json:"node_id_up,omitempty"`
	NodeIdSbi string `json:"node_id_sbi,omitempty"`
	Sst    string `json:"sst,omitempty"`
	Sd     string `json:"sd,omitempty"`
	Pool   string `json:"pool,omitempty"`
}

type Link struct {
	AEnd string `json:"aEnd,omitempty"`
	BEnd string `json:"bEnd,omitempty"`
}

type UE struct {
	SUPI            string   `json:"SUPI,omitempty"`
	AN              string   `json:"AN,omitempty"`
	DestinationIP   string   `json:"DestinationIP,omitempty"`
	DestinationPort string   `json:"DestinationPort,omitempty"`
	UPFPath         []string `json:"UPFPath,omitempty"`
}

type UEList struct {
	UES []UE `json:"ues,omitempty"`
}

// SmfConfig config parameters supplied via Smf CR
type SmfConfig struct {
	PodSettings  *PodSettingsSpec `json:"podSettings,omitempty"`
	ImageUrl     string           `json:"image,omitempty"`
	ImageUrlInit string           `json:"image_init,omitempty"`
	ImageExtUrl  string           `json:"image_ext,omitempty"`
	NrfIPAddress string           `json:"nrf_ip_address,omitempty"`
	NrfPort      string           `json:"nrf_port,omitempty"`
	Nodes        []UpNode         `json:"up_nodes,omitempty"`
	Links        []Link           `json:"up_links,omitempty"`
	UEList       UEList           `json:"ue_list,omitempty"`
	Name         string           `json:"smf_name,omitempty"`	
}

// SmfInternal config settings derived from other services
type SmfInternal struct {
	Version *string `json:"version,omitempty"`
}

// SmfSpec defines the desired state of Smf
type SmfSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Config   *SmfConfig   `json:"config,omitempty"`
	Internal *SmfInternal `json:"internal,omitempty"`
}

type SmfOutputs struct {
	IpAddr  string `json:"ip_address,omitempty"`
	PodName string `json:"podname,omitempty"`
}

// SmfStatus defines the observed state of Smf
type SmfStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this fil
	CRStatus `json:",inline"`
	Outputs  SmfOutputs `json:"outputs,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Smf is the Schema for the smfs API
type Smf struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SmfSpec   `json:"spec,omitempty"`
	Status SmfStatus `json:"status,omitempty"`
}

func (a Smf) GetStatus() *CRStatus {
	return &a.Status.CRStatus
}

func (a *Smf) GetObject() client.Object {
	return a
}

func (a *Smf) IsReady() bool {
	return a.Status.IsReady()
}

func (a *Smf) IsReconciled() bool {
	return a.Status.IsReconciled()
}

func (a *Smf) SetError(err error) {
	a.GetStatus().SetError(err)
}

func (a *Smf) SetReady(ready bool) {
	a.GetStatus().SetReady(ready)
}

func (a *Smf) SetReconciled(reconciled bool) {
	a.GetStatus().SetReconciled(reconciled)
}

//+kubebuilder:object:root=true

// SmfList contains a list of Smf
type SmfList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Smf `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Smf{}, &SmfList{})
}
