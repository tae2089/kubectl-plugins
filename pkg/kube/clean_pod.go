package kube

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	meta "k8s.io/client-go/metadata"
)

type PodStatus string

var (
	POD_STATUS_PENDING   PodStatus = "Pending"
	POD_STATUS_RUNNING   PodStatus = "Running"
	POD_STATUS_SUCCEEDED PodStatus = "Succeeded"
	POD_STATUS_FAILED    PodStatus = "Failed"
	POD_STATUS_UNKNOWN   PodStatus = "Unknown"
)

type CleanPodsOption struct {
	DryRun    []string
	PodStatus string
}

type CleanPodsOptionFunc func(*CleanPodsOption) error

func CleanPodsWithStatus(client meta.Interface, namespace string, options ...CleanPodsOptionFunc) error {
	d := &CleanPodsOption{}
	for _, opt := range options {
		opt(d)
	}
	gvr := schema.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"}

	output, err := client.Resource(gvr).List(context.Background(), metav1.ListOptions{
		FieldSelector: d.PodStatus,
	})
	if err != nil {
		return err
	}
	for _, item := range output.Items {
		fmt.Println(item.GetName())
	}

	if err := client.Resource(gvr).Namespace(namespace).DeleteCollection(context.Background(), metav1.DeleteOptions{
		DryRun: d.DryRun,
	}, metav1.ListOptions{
		FieldSelector: d.PodStatus,
	}); err != nil {
		return err
	}
	return nil
}

func WithPodStatus(status PodStatus) CleanPodsOptionFunc {
	return func(d *CleanPodsOption) error {
		d.PodStatus = "status.phase=" + string(status)
		return nil
	}
}

func WithDryRun(dryRun bool) CleanPodsOptionFunc {
	return func(d *CleanPodsOption) error {
		if dryRun {
			d.DryRun = []string{"All"}
		} else {
			d.DryRun = nil
		}
		return nil
	}
}
