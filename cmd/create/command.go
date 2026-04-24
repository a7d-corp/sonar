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
	"errors"
	"log"

	"github.com/glitchcrab/sonar/internal/app"
	"github.com/glitchcrab/sonar/internal/config"
	"github.com/glitchcrab/sonar/internal/k8sclient"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/client-go/kubernetes"
)

var (
	dryRun              bool
	image               string
	networkPolicy       bool
	nodeExec            bool
	nodeName            string
	podArgs             string
	podCommand          string
	podGroup            int64
	podUser             int64
	privileged          bool
	privilegeEscalation bool
	runAsNonRoot        bool
	unprivilegedPing    bool
)

func NewCommand() *cobra.Command {
	command := &cobra.Command{
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

--privileged (default: false)

Allow the pod to run as a privileged pod; must be provided at the same
time as --podsecuritypolicy to have any effect.

--networkpolicy (default: false)

Apply a NetworkPolicy which allows all ingress and egress traffic.

--node-name (default: none)

Attempt to schedule the pod on the named node.

--node-exec (default: false)

--unprivileged-ping (default: false)

Sets the 'net.ipv4.ping_group_range' sysctl to allow ping to be used
without root privileges.

Examples:

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
worker2.

"sonar create --dry-run" - prints the generated Kubernetes manifests
to stdout without applying them to the cluster.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			//err = runCreateCommand(cmd, args, v)
			err = runCreateCommand(cmd, args)
			if err != nil {
				return err
			} else {
				return nil
			}
		},
	}

	command.Flags().BoolVarP(&dryRun, "dry-run", "d", false, "print generated manifests to stdout only")
	command.Flags().StringVarP(&image, "image", "i", "busybox:latest", "image name (e.g. glitchcrab/ubuntu-debug:latest)")
	command.Flags().BoolVar(&networkPolicy, "networkpolicy", false, "create NetworkPolicy")
	command.Flags().BoolVar(&nodeExec, "node-exec", false, "spawn a container with root access to the node")
	command.Flags().StringVarP(&nodeName, "node-name", "", "", "node name to attempt to schedule the pod on")
	command.Flags().StringVarP(&podArgs, "pod-args", "a", "24h", "args to pass to pod command")
	command.Flags().StringVarP(&podCommand, "pod-command", "c", "sleep", "pod command (aka image entrypoint)")
	command.Flags().Int64VarP(&podGroup, "pod-groupid", "g", 1000, "groupID to run the pod as")
	command.Flags().Int64VarP(&podUser, "pod-userid", "u", 1000, "userID to run the pod as")
	command.Flags().BoolVar(&privileged, "privileged", false, "run a privileged container (assumes userID of 0)")
	command.Flags().BoolVar(&privilegeEscalation, "privilege-escalation", false, "allow privilege escalation")
	command.Flags().BoolVar(&runAsNonRoot, "non-root", true, "run the container as non-root (assumes userID of 0)")
	command.Flags().BoolVar(&unprivilegedPing, "unprivileged-ping", false, "allow a non-root user to use ping")

	return command
}

func runCreateCommand(command *cobra.Command, args []string) error {
	// Get the App instance from the command context
	a, err := app.GetApp(command)
	if err != nil {
		return err
	}

	v, err := app.GetViper(command)
	if err != nil {
		return err
	}

	// Wire in Viper.
	v, err = updateViperConfig(command, v)
	if err != nil {
		log.Fatalf("Error updating Viper config: %v", err)
	}

	opts := config.CreateConfig{
		DryRun:              dryRun,
		FullName:            a.Globals.FullName,
		Image:               v.GetString("image"),
		Labels:              a.Globals.Labels,
		Name:                a.Globals.Name,
		Namespace:           a.Globals.Namespace,
		NetworkPolicy:       v.GetBool("networkpolicy"),
		NodeExec:            nodeExec,
		NodeName:            nodeName,
		NonRoot:             v.GetBool("non-root"),
		PodArgs:             v.GetString("pod-args"),
		PodCommand:          v.GetString("pod-command"),
		PodGroup:            v.GetInt64("pod-groupid"),
		PodUser:             v.GetInt64("pod-userid"),
		Privileged:          v.GetBool("privileged"),
		PrivilegeEscalation: v.GetBool("privilege-escalation"),
		UnprivilegedPing:    v.GetBool("unprivileged-ping"),
	}

	if err := validateCreateConfig(command, &opts); err != nil {
		return err
	}

	// Create a Kubernetes clientset.
	var k8sClientSet *kubernetes.Clientset
	if !opts.DryRun {
		k8sClientSet, err = k8sclient.New(a.Globals.KubeContext, a.Globals.KubeConfig)
		if err != nil {
			return err
		}
	}

	ctx := context.TODO()

	var errs []error

	// Create the ServiceAccount
	saErr := createServiceAccount(k8sClientSet, ctx, opts)
	if saErr != nil {
		errs = append(errs, saErr)
	}

	// If set, create a NetworkPolicy
	if opts.NetworkPolicy {
		npErr := createNetworkPolicy(k8sClientSet, ctx, opts)
		if npErr != nil {
			errs = append(errs, npErr)
		}
	}

	// Create the Deployment
	deployErr := createDeployment(k8sClientSet, ctx, opts)
	if deployErr != nil {
		errs = append(errs, deployErr)
	}

	// If there were any validation errors, return them as a single error.
	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

// updateViperConfig updates a Viper instance with some create command flags.
func updateViperConfig(command *cobra.Command, v *viper.Viper) (*viper.Viper, error) {
	// Bind some more flags to Viper.

	flagsToBind := []string{
		"image",
		"networkpolicy",
		"pod-args",
		"pod-command",
		"pod-groupid",
		"pod-userid",
		"privileged",
		"privilege-escalation",
		"non-root",
		"unprivileged-ping",
	}

	var err error
	for _, flag := range flagsToBind {
		if err = v.BindPFlag(flag, command.Flags().Lookup(flag)); err != nil {
			return nil, err
		}
	}

	return v, nil
}
