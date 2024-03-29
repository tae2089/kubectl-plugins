package util

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/tae2089/kubectl-custom-cli/pkg/kube"
)

const (
	// Do not use any sorting, this is the default and acts as a value used
	// to catch other arguments that are passed in which are unsupported.
	SortFieldDefault = "none"

	// Sort by termination timestamp in ascending order.
	SortFieldTerminationTime = "time"
)

var (

// When using the namespace provided by the `--namespace/-n` flag or current context.
// This represents: Pod, Container and Termination Time
// singleNamespaceFormatting = "%s\t%s\t%s\n"
// When using the `all-namespaces` flag, we must show which namespace the pod was in, this becomes an extra column.
// This represents: Namespace, Pod, Container and Termination Time
// allNamespacesFormatting = "%s\t%s\t%s\t%s\n"
)

var t = tabwriter.NewWriter(os.Stdout, 10, 1, 5, ' ', 0)

func PrintPodsInfoInAllNameSpace(pods kube.Pods, noHeaders bool, filterType string) (err error) {
	var allNamespacesFormatting string = "%s\t%s\t%s\t%s\n"
	if !noHeaders {
		if filterType == "restart" {
			allNamespacesFormatting = "%s\t%s\t%s\t%s\t%s\n"
			_, err = fmt.Fprintf(t, allNamespacesFormatting, "NAMESPACE", "POD", "TERMINATION TIME", "RESTARTS", "EXIT CODE")
		} else {
			_, err = fmt.Fprintf(t, allNamespacesFormatting, "POD", "TERMINATION TIME", "RESTARTS")
		}
		if err != nil {
			return err
		}
	}
	for _, p := range pods {
		if filterType == "restart" {
			_, err = fmt.Fprintf(t, allNamespacesFormatting, p.Pod.Namespace, p.Pod.Name, p.TerminatedTime, p.RestartCount, p.ExitCode)
		} else {
			_, err = fmt.Fprintf(t, allNamespacesFormatting, p.Pod.Name, p.TerminatedTime, p.RestartCount)
		}
		if err != nil {
			return err
		}
	}
	t.Flush()
	return nil
}

func PrintPodsInfoInSingleNameSpace(pods kube.Pods, noHeaders bool, filterType string) (err error) {
	var singleNamespaceFormatting string = "%s\t%s\t%s\n"
	if !noHeaders {
		if filterType == "restart" {
			singleNamespaceFormatting = "%s\t%s\t%s\t%s\n"
			_, err = fmt.Fprintf(t, singleNamespaceFormatting, "POD", "TERMINATION TIME", "RESTARTS", "EXIT CODE")
		} else {
			_, err = fmt.Fprintf(t, singleNamespaceFormatting, "POD", "TERMINATION TIME", "RESTARTS")
		}
		if err != nil {
			return err
		}
	}
	for _, p := range pods {
		if filterType == "restart" {
			_, err = fmt.Fprintf(t, singleNamespaceFormatting, p.Pod.Name, p.TerminatedTime, p.RestartCount, p.ExitCode)
		} else {
			_, err = fmt.Fprintf(t, singleNamespaceFormatting, p.Pod.Name, p.TerminatedTime, p.RestartCount)
		}

		if err != nil {
			return err
		}
	}
	t.Flush()
	return nil
}
