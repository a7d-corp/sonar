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
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func NewClusterRoleBinding(k8sClientSet *kubernetes.Clientset, ctx context.Context, sonarConfig sonarconfig.SonarConfig) (err error) {
	// Define the ClusterRoleBinding
	crb := &rbacv1.ClusterRoleBinding{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "rbac.authorization.k8s.io/v1",
			Kind:       "ClusterRoleBinding",
		},
		ObjectMeta: metav1.ObjectMeta{
			Labels: sonarConfig.Labels,
			Name:   sonarConfig.Name,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      sonarConfig.Name,
				Namespace: sonarConfig.Namespace,
			},
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     sonarConfig.Name,
		},
	}

	// Create the ClusterRoleBinding
	_, err = k8sClientSet.RbacV1().ClusterRoleBindings().Create(ctx, crb, metav1.CreateOptions{})

	return err
}
