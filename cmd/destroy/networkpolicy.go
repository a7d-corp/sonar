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
package destroy

import (
	"context"
	"fmt"

	"github.com/glitchcrab/sonar/internal/config"
	"github.com/glitchcrab/sonar/internal/utils"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func deleteNetworkPolicy(k8sClientSet *kubernetes.Clientset, ctx context.Context, o config.DeleteConfig, force bool) (string, error) {
	// Set foreground deletion so the client waits for confirmation before proceeding
	deletePolicy := metav1.DeletePropagationForeground
	deleteOptions := metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}

	resourceType := "networkpolicy"
	name := fmt.Sprintf("%s/%s", o.Namespace, o.Name)

	var err error
	var ok bool
	if !force {
		ok = utils.ConfirmationPrompt(resourceType, name)
	} else {
		ok = true
	}
	if ok {
		err = k8sClientSet.NetworkingV1().NetworkPolicies(o.Namespace).Delete(ctx, o.Name, deleteOptions)
	}

	// Handle errors
	if statusError, isStatus := err.(*errors.StatusError); isStatus && statusError.Status().Reason == metav1.StatusReasonNotFound {
		// Skip deletion of this resource
		log.Info("no matching networkpolicy found; skipping deletion")
		return "", nil
	} else if err != nil {
		// Only return an error if the resource was not deleted
		return "", fmt.Errorf("networkpolicy \"%s\" failed deletion: %w", o.Name, err)
	} else {
		// Inform the user that the resource was deleted
		log.Infof("deleting networkpolicy")
	}

	return resourceType, nil
}
