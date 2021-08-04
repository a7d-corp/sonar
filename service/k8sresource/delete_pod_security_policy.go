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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func DeletePodSecurityPolicy(k8sClientSet *kubernetes.Clientset, ctx context.Context, sonarConfig sonarconfig.SonarConfig, force bool) (err error) {
	// Set foreground deletion so the client waits for confirmation before proceeding
	deletePolicy := metav1.DeletePropagationForeground
	deleteOptions := metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}

	resourceType := "podsecuritypolicy"

	var ok bool
	if !force {
		ok = helpers.ConfirmationPrompt(resourceType, sonarConfig.Name)
	} else {
		ok = true
	}
	if ok {
		err = k8sClientSet.PolicyV1beta1().PodSecurityPolicies().Delete(ctx, sonarConfig.Name, deleteOptions)
	}

	return err
}
