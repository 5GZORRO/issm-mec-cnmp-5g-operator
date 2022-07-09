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

// Pod resources
type ResourcesSpec struct {
	Requests RequestsSpec `json:"requests,omitempty"`
	Limits   LimitsSpec   `json:"limits,omitempty"`
}

// Pod resource requests
type RequestsSpec struct {
	Memory *string `json:"memory,omitempty"`
	CPU    *string `json:"cpu,omitempty"`
}

// Pod resource limits
type LimitsSpec struct {
	Memory *string `json:"memory,omitempty"`
	CPU    *string `json:"cpu,omitempty"`
}

// Pod liveness settings
type LivenessProbeSpec struct {
	InitialDelaySeconds *int32 `json:"initialDelaySeconds,omitempty"`
	PeriodSeconds       *int32 `json:"periodSeconds,omitempty"`
	TimeoutSeconds      *int32 `json:"timeoutSeconds,omitempty"`
	FailureThreshold    *int32 `json:"failureThreshold,omitempty"`
}

// Pod readiness settings
type ReadinessProbeSpec struct {
	InitialDelaySeconds *int32 `json:"initialDelaySeconds,omitempty"`
	PeriodSeconds       *int32 `json:"periodSeconds,omitempty"`
	SuccessThreshold    *int32 `json:"successThreshold,omitempty"`
	TimeoutSeconds      *int32 `json:"timeoutSeconds,omitempty"`
	FailureThreshold    *int32 `json:"failureThreshold,omitempty"`
}

// Persistence settings
type PersistanceSpec struct {
	StorageClassName *string `json:"storageClassName,omitempty"`
	StorageSize      *string `json:"storageSize,omitempty"`
}

// Pod settings
type PodSettingsSpec struct {
	Resources *ResourcesSpec `json:"resources,omitempty"`
	//LivenessProbe                 *LivenessProbeSpec   `json:"livenessProbe,omitempty"`
	//ReadinessProbe                *ReadinessProbeSpec  `json:"readinessProbe,omitempty"`
	//Affinity                      *AffinityInternal    `json:"affinity,omitempty"`
	//Tolerations                   *[]corev1.Toleration `json:"tolerations,omitempty"`
}
