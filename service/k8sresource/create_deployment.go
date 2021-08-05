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
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var (
	replicas int32 = 1
)

func NewDeployment(k8sClientSet *kubernetes.Clientset, ctx context.Context, sonarConfig sonarconfig.SonarConfig) (err error) {
	// Define the Deployment
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Labels:    sonarConfig.Labels,
			Name:      sonarConfig.Name,
			Namespace: sonarConfig.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: sonarConfig.Labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: sonarConfig.Labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Image: sonarConfig.Image,
							Name:  "sonar",
							Resources: corev1.ResourceRequirements{
								Limits: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse("2"),
									corev1.ResourceMemory: resource.MustParse("250Mi"),
								},
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse("200m"),
									corev1.ResourceMemory: resource.MustParse("50Mi"),
								},
							},
						},
					},
					RestartPolicy:      corev1.RestartPolicyAlways,
					ServiceAccountName: sonarConfig.Name,
					SecurityContext: &corev1.PodSecurityContext{
						RunAsUser: &sonarConfig.PodUser,
					},
				},
			},
		},
	}

	// Create the Deployment
	_, err = k8sClientSet.AppsV1().Deployments(sonarConfig.Namespace).Create(ctx, deployment, metav1.CreateOptions{})

	return err
}
