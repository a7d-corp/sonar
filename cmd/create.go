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
package cmd

import (
	"context"
	"fmt"
	"regexp"

	"github.com/glitchcrab/sonar/internal/sonarconfig"
	"github.com/glitchcrab/sonar/service/k8sclient"
	"github.com/glitchcrab/sonar/service/k8sresource"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	imageRegex = "^[a-z0-9/.-]*[:][a-z0-9.-]*$"
)

var (
	dryRun            bool
	image             string
	networkPolicy     bool
	nodeExec          bool
	nodeName          string
	podArgs           string
	podCommand        string
	podSecurityPolicy bool
	podUser           int64
	privileged        bool

	createCmd = &cobra.Command{
		Use:     "create",
		Aliases: []string{"deploy"},
		Short:   "Create applies a debug deployment to a Kubernetes cluster",
		Long: `Create will attempt to create a debugging deployment and all supporting
resources in the provided kubectl context (or the current context if
none is provided). Sonar assumes that your context has the required
privileges to create the necessary resources.

All flags are optional as defaults are provided.

Note: it is safe to run "sonar create" multiple times; if a resource
already exists then it will be skipped. For example, this can be used
to add a NetworkPolicy to an existing Sonar deployment which was
created without it.

Global flags:

Run "sonar help" in order to see flags which apply to all subcommands.

Flags:

--dry-run (default: False)

Prints the generated manifests to stdout only.

--image (default: 'busybox:latest')

Name of the image to use. Image names may be provided with or without a
tag; if no tag is detected then 'latest' is automatically used.

--pod-cmd (default: 'sleep')

Command to use as the entrypoint.

--pod-args (default: '24h')

Args to pass to the command.

--pod-userid (default: 1000)

User ID to run the container as (set in deployment's SecurityContext).

--podsecuritypolicy (default: false)

Create a PodSecurityPolicy and the associated ClusterRole and Binding.
The PSP will inherit the value set via --pod-userid and configure the
minimum value of the RunAs range accordingly.

--privileged (default: false)

Allow the pod to run as a privileged pod; must be provided at the same
time as --podsecuritypolicy to have any effect.

--networkpolicy (default: false)

Apply a NetworkPolicy which allows all ingress and egress traffic.

--node-name (default: none)

Attempt to schedule the pod on the named node.i

--node-exec (default: false)

Create a privileged pod in the node's PID & network namespaces. A node
name to schedule onto must also be provided. Note that the following
flags will be ignored: networkpolicy, podsecuritypolicy, privileged.`,
		Example: `
"sonar create" - accept all defaults. Creates a deployment in namespace
'default' called 'sonar-debug'.  The pod image will be 'busybox:latest'
with 'sleep 24h' as the initial command.

"sonar create --image glitchcrab/ubuntu-debug:v1.0 --pod-cmd sleep \
    --pod-args 1h --node-name worker10" - uses the provided image,
command and args, and attempts to schedule the pod on node 'worker10'.

"sonar create --podsecuritypolicy --pod-userid 0 --privileged" - creates
a deployment which runs as root. Also creates a PodSecurityPolicy
(and associated RBAC) which allows the pod to run as root/privileged.

"sonar create --networkpolicy" - creates a NetworkPolicy which allows
all ingress and traffic to the Sonar pod.

"sonar create --node-exec true --node-name worker2 \
    --pod-userid 0" - creates a pod with root access to the node named
worker2.`,
		Run: createSonarDeployment,
	}
)

func init() {
	rootCmd.AddCommand(createCmd)
	cobra.OnInitialize(validateFlags)

	createCmd.Flags().BoolVarP(&dryRun, "dry-run", "d", false, "print manifests to stdout (default \"false\")")
	createCmd.Flags().StringVarP(&image, "image", "i", "busybox:latest", "image name (e.g. glitchcrab/ubuntu-debug:latest)")
	createCmd.Flags().BoolVar(&networkPolicy, "networkpolicy", false, "create NetworkPolicy (default \"false\")")
	createCmd.Flags().BoolVar(&nodeExec, "node-exec", false, "spawn a container with root access to the node (default \"false\")")
	createCmd.Flags().StringVarP(&nodeName, "node-name", "", "", "node name to attempt to schedule the pod on")
	createCmd.Flags().StringVarP(&podArgs, "pod-args", "a", "24h", "args to pass to pod command")
	createCmd.Flags().StringVarP(&podCommand, "pod-command", "c", "sleep", "pod command (aka image entrypoint)")
	createCmd.Flags().BoolVar(&podSecurityPolicy, "podsecuritypolicy", false, "create PodSecurityPolicy (default \"false\")")
	createCmd.Flags().Int64VarP(&podUser, "pod-userid", "u", 1000, "userID to run the pod as")
	createCmd.Flags().BoolVar(&privileged, "privileged", false, "run the container as root (assumes userID of 0) (default \"false\")")
}

func validateFlags() {
	// If the user has provided an image name then validate that it looks sane.
	if createCmd.Flags().Lookup("image").Changed {
		// Validate image to see if a tag has been provided; if not then
		// use :latest. Does not validate full image name, just whether a
		// tag was provided.
		ok, _ := regexp.MatchString(imageRegex, image)
		if !ok {
			image = fmt.Sprintf("%s:latest", image)
		}
	}
}

func createSonarDeployment(cmd *cobra.Command, args []string) {
	// Set sane options if we're exec-ing into a node.
	if nodeExec {
		// Error out if node name was not provided.
		if nodeName == "" {
			log.Fatal("--node-exec also requires --node-name to be provided")
		}

		networkPolicy = false
		podSecurityPolicy = true
		privileged = true
	}

	// Create a SonarConfig
	sonarConfig := sonarconfig.SonarConfig{
		DryRun:            dryRun,
		Image:             image,
		Labels:            labels,
		Name:              name,
		Namespace:         namespace,
		NetworkPolicy:     networkPolicy,
		NodeExec:          nodeExec,
		NodeName:          nodeName,
		PodArgs:           podArgs,
		PodCommand:        podCommand,
		PodSecurityPolicy: podSecurityPolicy,
		PodUser:           podUser,
		Privileged:        privileged,
	}

	// Create a clientset to interact with the cluster.
	k8sClientSet, err := k8sclient.New(kubeContext, kubeConfig)
	if err != nil {
		log.Fatal(err) // TODO: better logging
	}

	// Create a context
	ctx := context.TODO()

	{
		// Create a ServiceAccount
		err := k8sresource.NewServiceAccount(k8sClientSet, ctx, sonarConfig)
		// Handle the response
		if statusError, isStatus := err.(*errors.StatusError); isStatus && statusError.Status().Reason == metav1.StatusReasonAlreadyExists {
			log.Warnf("serviceaccount \"%s/%s\" already exists\n", sonarConfig.Namespace, sonarConfig.Name)
		} else if err != nil {
			log.Warnf("serviceaccount \"%s/%s\" was not created: %w\n", sonarConfig.Namespace, sonarConfig.Name, err)
		} else {
			log.Infof("serviceaccount \"%s/%s\" created\n", sonarConfig.Namespace, sonarConfig.Name)
		}
	}

	if sonarConfig.PodSecurityPolicy {
		{
			// Create a PodSecurityPolicy
			err := k8sresource.NewPodSecurityPolicy(k8sClientSet, ctx, sonarConfig)
			// Handle the response
			if statusError, isStatus := err.(*errors.StatusError); isStatus && statusError.Status().Reason == metav1.StatusReasonAlreadyExists {
				log.Warnf("podsecuritypolicy \"%s\" already exists\n", sonarConfig.Name)
			} else if err != nil {
				log.Warnf("podsecuritypolicy \"%s\" was not created: %w\n", sonarConfig.Name, err)
			} else {
				log.Infof("podsecuritypolicy \"%s\" created\n", sonarConfig.Name)
			}
		}

		{
			// Create a ClusterRole
			err := k8sresource.NewClusterRole(k8sClientSet, ctx, sonarConfig)
			// Handle the response
			if statusError, isStatus := err.(*errors.StatusError); isStatus && statusError.Status().Reason == metav1.StatusReasonAlreadyExists {
				log.Warnf("clusterrole \"%s\" already exists\n", sonarConfig.Name)
			} else if err != nil {
				log.Warnf("clusterrole \"%s\" was not created: %w\n", sonarConfig.Name, err)
			} else {
				log.Infof("clusterrole \"%s\" created\n", sonarConfig.Name)
			}
		}

		{
			// Create a ClusterRoleBinding
			err := k8sresource.NewClusterRoleBinding(k8sClientSet, ctx, sonarConfig)
			// Handle the response
			if statusError, isStatus := err.(*errors.StatusError); isStatus && statusError.Status().Reason == metav1.StatusReasonAlreadyExists {
				log.Warnf("clusterrolebinding \"%s\" already exists\n", sonarConfig.Name)
			} else if err != nil {
				log.Warnf("clusterrolebinding \"%s\" was not created: %w\n", sonarConfig.Name, err)
			} else {
				log.Infof("clusterrolebinding \"%s\" created\n", sonarConfig.Name)
			}
		}
	}

	if sonarConfig.NetworkPolicy {
		// Create a NetworkPolicy
		err := k8sresource.NewNetworkPolicy(k8sClientSet, ctx, sonarConfig)
		// Handle the response
		if statusError, isStatus := err.(*errors.StatusError); isStatus && statusError.Status().Reason == metav1.StatusReasonAlreadyExists {
			log.Warnf("networkpolicy \"%s\" already exists\n", sonarConfig.Name)
		} else if err != nil {
			log.Warnf("networkpolicy \"%s\" was not created: %w\n", sonarConfig.Name, err)
		} else {
			log.Infof("networkpolicy \"%s\" created\n", sonarConfig.Name)
		}
	}

	{
		// Create a Deployment
		err := k8sresource.NewDeployment(k8sClientSet, ctx, sonarConfig)
		// Handle the response
		if statusError, isStatus := err.(*errors.StatusError); isStatus && statusError.Status().Reason == metav1.StatusReasonAlreadyExists {
			log.Warnf("deployment \"%s/%s\" already exists\n", sonarConfig.Namespace, sonarConfig.Name)
		} else if err != nil {
			log.Warnf("deployment \"%s/%s\" was not created: %w\n", sonarConfig.Namespace, sonarConfig.Name, err)
		} else {
			log.Infof("deployment \"%s/%s\" created\n", sonarConfig.Namespace, sonarConfig.Name)
		}
	}
}
