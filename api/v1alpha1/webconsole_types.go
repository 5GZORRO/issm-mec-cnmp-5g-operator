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

// WebconsoleConfig config parameters supplied via Webconsole CR
type WebconsoleConfig struct {
	PodSettings *PodSettingsSpec `json:"podSettings,omitempty"`
	ImageUrl    string           `json:"image,omitempty"`
	ImageUrlInit string           `json:"image_init,omitempty"`	
	MongoIPAddr string           `json:"mongo_ip_address,omitempty"`
}

// WebconsoleInternal config settings derived from other services
type WebconsoleInternal struct {
	Version *string `json:"version,omitempty"`
}

type WebconsoleOutputs struct {
	IpAddr  string `json:"ip_address,omitempty"`
	PodName string `json:"podname,omitempty"`
}

// WebconsoleStatus defines the observed state of Webconsole
type WebconsoleStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	CRStatus `json:",inline"`
	Outputs  WebconsoleOutputs `json:"outputs,omitempty"`
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// WebconsoleSpec defines the desired state of Webconsole
type WebconsoleSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Config   *WebconsoleConfig   `json:"config,omitempty"`
	Internal *WebconsoleInternal `json:"internal,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Webconsole is the Schema for the webconsoles API
type Webconsole struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   WebconsoleSpec   `json:"spec,omitempty"`
	Status WebconsoleStatus `json:"status,omitempty"`
}

func (a Webconsole) GetStatus() *CRStatus {
	return &a.Status.CRStatus
}

func (a *Webconsole) GetObject() client.Object {
	return a
}

func (a *Webconsole) IsReady() bool {
	return a.Status.IsReady()
}

func (a *Webconsole) IsReconciled() bool {
	return a.Status.IsReconciled()
}

func (a *Webconsole) SetError(err error) {
	a.GetStatus().SetError(err)
}

func (a *Webconsole) SetReady(ready bool) {
	a.GetStatus().SetReady(ready)
}

func (a *Webconsole) SetReconciled(reconciled bool) {
	a.GetStatus().SetReconciled(reconciled)
}

//+kubebuilder:object:root=true

// WebconsoleList contains a list of Webconsole
type WebconsoleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Webconsole `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Webconsole{}, &WebconsoleList{})
}
