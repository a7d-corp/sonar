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

	"github.com/glitchcrab/sonar/internal/sonarconfig"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var (
	minRunAsID int64 = 1000
	maxRunAsID int64 = 65535
)

func NewPodSecurityPolicy(k8sClientSet *kubernetes.Clientset, ctx context.Context, sonarConfig sonarconfig.SonarConfig) (err error) {
	if sonarConfig.PodUser != minRunAsID {
		minRunAsID = sonarConfig.PodUser
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
			FSGroup: policyv1beta1.FSGroupStrategyOptions{
				Rule: policyv1beta1.FSGroupStrategyMustRunAs,
				Ranges: []policyv1beta1.IDRange{
					{
						Min: minRunAsID,
						Max: maxRunAsID,
					},
				},
			},
			HostIPC:                false,
			HostNetwork:            false,
			HostPID:                false,
			Privileged:             sonarConfig.Privileged,
			ReadOnlyRootFilesystem: false,
			RunAsGroup: &policyv1beta1.RunAsGroupStrategyOptions{
				Rule: policyv1beta1.RunAsGroupStrategyMustRunAs,
				Ranges: []policyv1beta1.IDRange{
					{
						Min: minRunAsID,
						Max: maxRunAsID,
					},
				},
			},
			RunAsUser: policyv1beta1.RunAsUserStrategyOptions{
				Rule: policyv1beta1.RunAsUserStrategyMustRunAs,
				Ranges: []policyv1beta1.IDRange{
					{
						Min: minRunAsID,
						Max: maxRunAsID,
					},
				},
			},
			SELinux: policyv1beta1.SELinuxStrategyOptions{
				Rule: policyv1beta1.SELinuxStrategyRunAsAny,
			},
			SupplementalGroups: policyv1beta1.SupplementalGroupsStrategyOptions{
				Rule: policyv1beta1.SupplementalGroupsStrategyRunAsAny,
			},
			Volumes: []policyv1beta1.FSType{
				"*",
			},
		},
	}

	// Create the PodSecurityPolicy
	_, err = k8sClientSet.PolicyV1beta1().PodSecurityPolicies().Create(ctx, psp, metav1.CreateOptions{})

	return err
}
