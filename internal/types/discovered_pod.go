package types

import (
	corev1 "k8s.io/api/core/v1"
)

// DiscoveredPod represents a pod which is a candidate for execing into.
type DiscoveredPod struct {
	Name      string
	Namespace string
	Status    corev1.PodPhase
}
