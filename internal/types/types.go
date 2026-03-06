package types

// DiscoveredPod represents a pod which is a candidate for execing into.
type DiscoveredPod struct {
	Name      string
	Namespace string
}
