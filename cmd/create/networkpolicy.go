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
package create

import (
	"context"
	"fmt"

	"github.com/glitchcrab/sonar/internal/config"
	"github.com/glitchcrab/sonar/internal/helpers"
	log "github.com/sirupsen/logrus"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func createNetworkPolicy(k8sClientSet *kubernetes.Clientset, ctx context.Context, o config.CreateConfig) error {
	// Define the NetworkPolicy
	np := &networkingv1.NetworkPolicy{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "networking.k8s.io/v1",
			Kind:       "NetworkPolicy",
		},
		ObjectMeta: metav1.ObjectMeta{
			Labels:    o.Labels,
			Name:      o.Name,
			Namespace: o.Namespace,
		},
		Spec: networkingv1.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{
				MatchLabels: o.Labels,
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

	var err error

	if o.DryRun {
		err = helpers.PrintManifestYAML(np)
		if err != nil {
			return fmt.Errorf("networkpolicy \"%s\" manifest generation failed: %v\n", o.Name, err)
		}
	} else {
		_, err = k8sClientSet.NetworkingV1().NetworkPolicies(o.Namespace).Create(ctx, np, metav1.CreateOptions{})
	}

	if err != nil {
		if statusError, isStatus := err.(*errors.StatusError); isStatus && statusError.Status().Reason == metav1.StatusReasonAlreadyExists {
			return fmt.Errorf("networkpolicy \"%s\" already exists\n", o.Name)
		} else if err != nil {
			return fmt.Errorf("networkpolicy \"%s\" was not created: %w\n", o.Name, err)
		}
	} else {
		log.Infof("networkpolicy \"%s\" created\n", o.Name)
	}

	return nil
}
