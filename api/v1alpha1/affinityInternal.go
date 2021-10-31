package v1alpha1

/**
 * An internal representation of corev1.Affinity object tree.
 * Functionally unnecessary, but corev1.Affinity object las long comments which get added to CRD's
 * resulting in CRD's larger than man allowable
 */

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/yaml"
)

var log = logf.Log.WithName("affinity")

type AffinityInternal struct {
	NodeAffinity    *NodeAffinityInternal    `json:"nodeAffinity,omitempty" protobuf:"bytes,1,opt,name=nodeAffinity"`
	PodAffinity     *PodAffinityInternal     `json:"podAffinity,omitempty" protobuf:"bytes,2,opt,name=podAffinity"`
	PodAntiAffinity *PodAntiAffinityInternal `json:"podAntiAffinity,omitempty" protobuf:"bytes,3,opt,name=podAntiAffinity"`
}

type NodeAffinityInternal struct {
	RequiredDuringSchedulingIgnoredDuringExecution  *NodeSelectorInternal             `json:"requiredDuringSchedulingIgnoredDuringExecution,omitempty" protobuf:"bytes,1,opt,name=requiredDuringSchedulingIgnoredDuringExecution"`
	PreferredDuringSchedulingIgnoredDuringExecution []PreferredSchedulingTermInternal `json:"preferredDuringSchedulingIgnoredDuringExecution,omitempty" protobuf:"bytes,2,rep,name=preferredDuringSchedulingIgnoredDuringExecution"`
}

type PodAffinityInternal struct {
	RequiredDuringSchedulingIgnoredDuringExecution  []PodAffinityTermInternal         `json:"requiredDuringSchedulingIgnoredDuringExecution,omitempty" protobuf:"bytes,1,rep,name=requiredDuringSchedulingIgnoredDuringExecution"`
	PreferredDuringSchedulingIgnoredDuringExecution []WeightedPodAffinityTermInternal `json:"preferredDuringSchedulingIgnoredDuringExecution,omitempty" protobuf:"bytes,2,rep,name=preferredDuringSchedulingIgnoredDuringExecution"`
}

type PodAntiAffinityInternal struct {
	RequiredDuringSchedulingIgnoredDuringExecution  []PodAffinityTermInternal         `json:"requiredDuringSchedulingIgnoredDuringExecution,omitempty" protobuf:"bytes,1,rep,name=requiredDuringSchedulingIgnoredDuringExecution"`
	PreferredDuringSchedulingIgnoredDuringExecution []WeightedPodAffinityTermInternal `json:"preferredDuringSchedulingIgnoredDuringExecution,omitempty" protobuf:"bytes,2,rep,name=preferredDuringSchedulingIgnoredDuringExecution"`
}

type NodeSelectorInternal struct {
	NodeSelectorTerms []NodeSelectorTermInternal `json:"nodeSelectorTerms" protobuf:"bytes,1,rep,name=nodeSelectorTerms"`
}

type NodeSelectorTermInternal struct {
	MatchExpressions []NodeSelectorRequirementInternal `json:"matchExpressions,omitempty" protobuf:"bytes,1,rep,name=matchExpressions"`
	MatchFields      []NodeSelectorRequirementInternal `json:"matchFields,omitempty" protobuf:"bytes,2,rep,name=matchFields"`
}

type NodeSelectorRequirementInternal struct {
	Key      string                      `json:"key" protobuf:"bytes,1,opt,name=key"`
	Operator corev1.NodeSelectorOperator `json:"operator" protobuf:"bytes,2,opt,name=operator,casttype=NodeSelectorOperator"`
	Values   []string                    `json:"values,omitempty" protobuf:"bytes,3,rep,name=values"`
}

type PreferredSchedulingTermInternal struct {
	Weight     int32                    `json:"weight" protobuf:"varint,1,opt,name=weight"`
	Preference NodeSelectorTermInternal `json:"preference" protobuf:"bytes,2,opt,name=preference"`
}

type PodAffinityTermInternal struct {
	LabelSelector *LabelSelectorInternal `json:"labelSelector,omitempty" protobuf:"bytes,1,opt,name=labelSelector"`
	Namespaces    []string               `json:"namespaces,omitempty" protobuf:"bytes,2,rep,name=namespaces"`
	TopologyKey   string                 `json:"topologyKey" protobuf:"bytes,3,opt,name=topologyKey"`
}

type WeightedPodAffinityTermInternal struct {
	Weight          int32                   `json:"weight" protobuf:"varint,1,opt,name=weight"`
	PodAffinityTerm PodAffinityTermInternal `json:"podAffinityTerm" protobuf:"bytes,2,opt,name=podAffinityTerm"`
}

type LabelSelectorInternal struct {
	MatchLabels      map[string]string                  `json:"matchLabels,omitempty" protobuf:"bytes,1,rep,name=matchLabels"`
	MatchExpressions []LabelSelectorRequirementInternal `json:"matchExpressions,omitempty" protobuf:"bytes,2,rep,name=matchExpressions"`
}

type LabelSelectorRequirementInternal struct {
	Key      string                       `json:"key" patchStrategy:"merge" patchMergeKey:"key" protobuf:"bytes,1,opt,name=key"`
	Operator metav1.LabelSelectorOperator `json:"operator" protobuf:"bytes,2,opt,name=operator,casttype=LabelSelectorOperator"`
	Values   []string                     `json:"values,omitempty" protobuf:"bytes,3,rep,name=values"`
}

func (a AffinityInternal) ConvertToAffinity() *corev1.Affinity {
	affinityYaml, err := yaml.Marshal(a)
	if err == nil {
		affinity := &corev1.Affinity{}
		err = yaml.Unmarshal(affinityYaml, affinity)
		if err == nil {
			return affinity
		}
	}
	log.Error(err, "Error converting affinity object in ConvertToAffinity")
	return nil
}
