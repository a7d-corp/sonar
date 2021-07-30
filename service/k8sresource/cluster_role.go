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
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func NewClusterRole(k8sClientSet *kubernetes.Clientset, ctx context.Context, sonarConfig sonarconfig.SonarConfig) (err error) {
	// Define the ClusterRole
	cr := &rbacv1.ClusterRole{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "rbac.authorization.k8s.io/v1",
			Kind:       "ClusterRole",
		},
		ObjectMeta: metav1.ObjectMeta{
			Labels: sonarConfig.Labels,
			Name:   sonarConfig.Name,
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups:     []string{"policy"},
				Resources:     []string{"podsecuritypolicies"},
				ResourceNames: []string{sonarConfig.Name},
				Verbs:         []string{"use"},
			},
		},
	}

	// Create the ClusterRole
	_, err = k8sClientSet.RbacV1().ClusterRoles().Create(ctx, cr, metav1.CreateOptions{})

	return err
}
