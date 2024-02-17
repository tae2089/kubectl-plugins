package kube

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func GetNamespace(configFlags *genericclioptions.ConfigFlags, all bool) (string, error) {

	if all {
		return metav1.NamespaceAll, nil
	}

	if *configFlags.Namespace != "" {
		ns := configFlags.Namespace
		return *ns, nil
	}

	// Retrieve the current namespace from the raw kubeconfig struct
	// when all namespaces are not requested and no namespace is given.
	currentNamespace, _, err := configFlags.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		return "", fmt.Errorf("failed to during creating raw kubeconfig: %w", err)
	}
	return currentNamespace, nil
}
