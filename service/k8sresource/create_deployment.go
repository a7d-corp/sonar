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
	"strings"

	"github.com/glitchcrab/sonar/internal/sonarconfig"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var (
	hostIPC        = false
	hostNet        = false
	hostPID        = false
	replicas int32 = 1
)

func NewDeployment(k8sClientSet *kubernetes.Clientset, ctx context.Context, sonarConfig sonarconfig.SonarConfig) (err error) {
	// Create container in the host namespaces if node-exec is set.
	if sonarConfig.NodeExec {
		hostIPC = true
		hostNet = true
		hostPID = true
	}

	securityContext := &corev1.SecurityContext{
		Capabilities: &corev1.Capabilities{
			Drop: []corev1.Capability{"ALL"},
		},
		Privileged:               &sonarConfig.Privileged,
		RunAsUser:                &sonarConfig.PodUser,
		RunAsGroup:               &sonarConfig.PodGroup,
		RunAsNonRoot:             &sonarConfig.NonRoot,
		AllowPrivilegeEscalation: &sonarConfig.PrivilegeEscalation,
		SeccompProfile: &corev1.SeccompProfile{
			Type: "RuntimeDefault",
		},
	}

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
							SecurityContext: securityContext,
						},
					},
					HostIPC:            hostIPC,
					HostNetwork:        hostNet,
					HostPID:            hostPID,
					RestartPolicy:      corev1.RestartPolicyAlways,
					ServiceAccountName: sonarConfig.Name,
				},
			},
		},
	}

	// Update the deployment's command if one was provided.
	if sonarConfig.PodCommand != "" {
		command := strings.Fields(sonarConfig.PodCommand)
		deployment.Spec.Template.Spec.Containers[0].Command = command
	}

	// Update the deployment's command if one was provided.
	if sonarConfig.PodArgs != "" {
		cmdargs := strings.Fields(sonarConfig.PodArgs)
		deployment.Spec.Template.Spec.Containers[0].Args = cmdargs
	}

	// Add the NodeName if one was provided.
	if sonarConfig.NodeName != "" {
		deployment.Spec.Template.Spec.NodeName = sonarConfig.NodeName
	}

	// Mount the hosts's filesystem if exec-ing into a node.
	if sonarConfig.NodeExec {
		// create the volume.
		deployment.Spec.Template.Spec.Volumes = []corev1.Volume{
			{
				Name: "host-rootfs",
				VolumeSource: corev1.VolumeSource{
					HostPath: &corev1.HostPathVolumeSource{
						Path: "/",
					},
				},
			},
		}

		// attach it to the container
		deployment.Spec.Template.Spec.Containers[0].VolumeMounts = []corev1.VolumeMount{
			{
				Name:      "host-rootfs",
				MountPath: "/host",
			},
		}
	}

	// Create the Deployment
	_, err = k8sClientSet.AppsV1().Deployments(sonarConfig.Namespace).Create(ctx, deployment, metav1.CreateOptions{})

	return err
}
