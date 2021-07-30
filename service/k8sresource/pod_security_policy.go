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

	"github.com/glitchcrab/sonar/pkg/sonarconfig"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var (
	minRunAsID = 2
	maxRunAsID = 65535
)

func NewPodSecurityPolicy(k8sClientSet *kubernetes.Clientset, ctx context.Context, sonarConfig sonarconfig.SonarConfig) (err error) {
	if sonarConfig.Privileged {
		minRunAsID = 1
	}

	// Define the PodSecurityPolicy
	psp := &policyv1beta1.PodSecurityPolicy{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "policy/v1beta1",
			Kind:       "PodSecurityPolicy",
		},
		ObjectMeta: metav1.ObjectMeta{
			Labels:    sonarConfig.Labels,
			Name:      sonarConfig.Name,
			Namespace: sonarConfig.Namespace,
		},
		Spec: policyv1beta1.PodSecurityPolicySpec{
			AllowPrivilegeEscalation: &sonarConfig.Privileged,
			//FSGroup:                  []policyv1beta1.FSGroupStrategyOptions{
			//Rule: FSGroupStrategyMustRunAs,
			//Ranges: []policyv1beta1.IDRange{
			//	Min: 2,
			//	Max: 65535,
			//},
			//},
			HostIPC:     false,
			HostNetwork: false,
			HostPID:     false,
			Privileged:  sonarConfig.Privileged,
			//RunAsGroup: []policyv1beta1.RunAsGroupStrategyOptions{
			//	Rule: RunAsGroupStrategyMayRunAs,
			//	Ranges: []policyv1beta1.IDRange{
			//		Min: 2,
			//		Max: 65535,
			//	},
			//},
			//RunAsUser: []policyv1beta1.RunAsUserStrategyOptions{
			//	Rule: RunAsUserStrategyMustRunAs,
			//	Ranges: []policyv1beta1.IDRange{
			//		Min: 2,
			//		Max: 65535,
			//	},
			//},
			//SELinux: []policyv1beta1.SELinuxStrategyOptions{
			//	Rule: SELinuxStrategyRunAsAny,
			//},
			//SupplementalGroups: []policyv1beta1.SupplementalGroupsStrategyOptions{
			//	Rule: SupplementalGroupsStrategyRunAsAny,
			//},
			Volumes: []policyv1beta1.FSType{
				"*",
			},
		},
	}

	// Create the PodSecurityPolicy
	_, err = k8sClientSet.PolicyV1beta1().PodSecurityPolicies().Create(ctx, psp, metav1.CreateOptions{})

	return err
}
