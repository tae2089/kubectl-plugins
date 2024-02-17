package kube

import (
	"fmt"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func getK8sClientAndConfig(configFlags *genericclioptions.ConfigFlags) (*kubernetes.Clientset, *rest.Config, error) {

	config, err := configFlags.ToRESTConfig()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read kubeconfig: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create clientset: %w", err)
	}
	return clientset, config, nil
}
