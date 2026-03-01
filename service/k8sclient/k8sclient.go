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
package k8sclient

import (
	"errors"
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/clientcmd"

	log "github.com/sirupsen/logrus"
)

// findKubeConfig discovers the current kubeconfig
func findKubeConfig() (string, error) {
	// Check if KUBECONFIG environment variable is set
	if envKubeConfig := os.Getenv("KUBECONFIG"); envKubeConfig != "" {
		return envKubeConfig, nil
	}

	path, err := homedir.Expand("~/.kube/config")
	if err != nil {
		return "", err
	}

	return path, nil
}

func New(kubeContext, kubeConfigPath string) (*kubernetes.Clientset, error) {
	var err error

	// Discover the kubeconfig if an explicit path wasn't provided
	if kubeConfigPath == "" {
		kubeConfigPath, err = findKubeConfig()
		if err != nil {
			return nil, err
		}
	}

	// Inform the user which kubeconfig file is being used
	log.Infof("using kubeconfig: %s", kubeConfigPath)

	// Set defaults for creating a new ClientConfig.
	loadingRules := &clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeConfigPath}
	configOverrides := &clientcmd.ConfigOverrides{}

	// Set the context if it was provided.
	if kubeContext != "" {
		configOverrides = &clientcmd.ConfigOverrides{CurrentContext: kubeContext}
	}

	// Create the ClientConfig.
	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides).ClientConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Create a Clientset with the provided values.
	k8sClientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	return k8sClientSet, err
}

// GetNamespace returns the current namespace from a Kubeconfig
func GetNamespace(kubeConfigPath, kubeContext string) (string, error) {
	var err error
	var context string

	// Discover the kubeconfig if an explicit path wasn't provided
	if kubeConfigPath == "" {
		kubeConfigPath, err = findKubeConfig()
		if err != nil {
			return "", err
		}
	}

	kubeConfig, err := clientcmd.LoadFromFile(kubeConfigPath)
	if err != nil {
		return "", err
	}

	// Use the context if it was provided, otherwise use the current context
	if kubeContext != "" {
		context = kubeContext
	} else {
		context = kubeConfig.CurrentContext
	}

	namespace := kubeConfig.Contexts[context].Namespace
	if namespace == "" {
		return "", errors.New(fmt.Sprintf("No namespace set (context: %s, kubeconfig: %s)", context, kubeConfigPath))
	}

	return namespace, nil
}
