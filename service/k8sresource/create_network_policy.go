/*
Copyright © 2021 Simon Weald

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
package k8sresource

import (
	"context"

	"github.com/glitchcrab/sonar/internal/helpers"
	"github.com/glitchcrab/sonar/internal/sonarconfig"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func NewNetworkPolicy(k8sClientSet *kubernetes.Clientset, ctx context.Context, sonarConfig sonarconfig.SonarConfig) (err error) {
	// Define the NetworkPolicy
	np := &networkingv1.NetworkPolicy{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "networking.k8s.io/v1",
			Kind:       "NetworkPolicy",
		},
		ObjectMeta: metav1.ObjectMeta{
			Labels:    sonarConfig.Labels,
			Name:      sonarConfig.Name,
			Namespace: sonarConfig.Namespace,
		},
		Spec: networkingv1.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{
				MatchLabels: sonarConfig.Labels,
			},
			Ingress: []networkingv1.NetworkPolicyIngressRule{
				networkingv1.NetworkPolicyIngressRule{},
			},
			Egress: []networkingv1.NetworkPolicyEgressRule{
				networkingv1.NetworkPolicyEgressRule{},
			},
			PolicyTypes: []networkingv1.PolicyType{
				networkingv1.PolicyTypeIngress,
				networkingv1.PolicyTypeEgress,
			},
		},
	}

	// If dry-run is enabled, print the manifest and return
	if sonarConfig.DryRun {
		return helpers.PrintManifestYAML(np)
	}

	// Create the NetworkPolicy
	_, err = k8sClientSet.NetworkingV1().NetworkPolicies(sonarConfig.Namespace).Create(ctx, np, metav1.CreateOptions{})

	return err
}
