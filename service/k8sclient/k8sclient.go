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
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	log "github.com/sirupsen/logrus"
)

func New(kubeContext, kubeConfig string) (*kubernetes.Clientset, error) {
	// Set the kubeconfig to the default location if the path wasn't provided.
	if kubeConfig == "" {
		if home := homedir.HomeDir(); home != "" {
			kubeConfig = filepath.Join(home, ".kube", "config")
		}
	}

	// Set defaults for creating a new ClientConfig.
	loadingRules := &clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeConfig}
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
