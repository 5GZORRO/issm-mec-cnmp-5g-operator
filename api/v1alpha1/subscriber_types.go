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

// SubscriberConfig config parameters supplied via Subscriber CR
type SubscriberConfig struct {
	ImageUrl            string `json:"image,omitempty"`
	WebconsoleIPAddress string `json:"webconsole_ip_address,omitempty"`
	PlmnID              string `json:"plmnID,omitempty"`
	IMSI                string `json:"IMSI,omitempty"`
	Dnn                 string `json:"dnn,omitempty"`
}

// SubscriberInternal config settings derived from other services
type SubscriberInternal struct {
	Version *string `json:"version,omitempty"`
}

// SubscriberSpec defines the desired state of Subscriber
type SubscriberSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Config   *SubscriberConfig   `json:"config,omitempty"`
	Internal *SubscriberInternal `json:"internal,omitempty"`
}

type SubscriberOutputs struct {
	IpAddr  string `json:"ip_address,omitempty"`
	PodName string `json:"podname,omitempty"`
}

// SubscriberStatus defines the observed state of Subscriber
type SubscriberStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	CRStatus `json:",inline"`
	Outputs  SubscriberOutputs `json:"outputs,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Subscriber is the Schema for the subscribers API
type Subscriber struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SubscriberSpec   `json:"spec,omitempty"`
	Status SubscriberStatus `json:"status,omitempty"`
}

func (a Subscriber) GetStatus() *CRStatus {
	return &a.Status.CRStatus
}

func (a *Subscriber) GetObject() client.Object {
	return a
}

func (a *Subscriber) IsReady() bool {
	return a.Status.IsReady()
}

func (a *Subscriber) IsReconciled() bool {
	return a.Status.IsReconciled()
}

func (a *Subscriber) SetError(err error) {
	a.GetStatus().SetError(err)
}

func (a *Subscriber) SetReady(ready bool) {
	a.GetStatus().SetReady(ready)
}

func (a *Subscriber) SetReconciled(reconciled bool) {
	a.GetStatus().SetReconciled(reconciled)
}

//+kubebuilder:object:root=true

// SubscriberList contains a list of Subscriber
type SubscriberList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Subscriber `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Subscriber{}, &SubscriberList{})
}
