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
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func createServiceAccount(k8sClientSet *kubernetes.Clientset, ctx context.Context, o config.CreateConfig) error {
	// Define the ServiceAccount
	sa := &corev1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ServiceAccount",
		},
		ObjectMeta: metav1.ObjectMeta{
			Labels:    o.Labels,
			Name:      o.FullName,
			Namespace: o.Namespace,
		},
	}

	var err error

	// If dry-run is enabled, print the manifest and return
	if o.DryRun {
		err = helpers.PrintManifestYAML(sa)
		if err != nil {
			return fmt.Errorf("serviceaccount \"%s/%s\" manifest generation failed: %v\n", o.Namespace, o.Name, err)
		}
	} else {
		_, err = k8sClientSet.CoreV1().ServiceAccounts(o.Namespace).Create(ctx, sa, metav1.CreateOptions{})
	}

	if err != nil {
		if statusError, isStatus := err.(*errors.StatusError); isStatus && statusError.Status().Reason == metav1.StatusReasonAlreadyExists {
			return fmt.Errorf("serviceaccount \"%s/%s\" already exists\n", o.Namespace, o.Name)
		} else if err != nil {
			return fmt.Errorf("serviceaccount \"%s/%s\" was not created: %w\n", o.Namespace, o.Name, err)
		}
	} else {
		log.Infof("serviceaccount \"%s/%s\" created\n", o.Namespace, o.Name)
	}

	return nil
}
