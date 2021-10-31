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
