//go:build cleaner

package cleaner

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tae2089/kubectl-custom-cli/pkg/kube"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	meta "k8s.io/client-go/metadata"
)

var (

	// KubernetesConfigFlags provides the generic flags which are available to
	// regular `kubectl` commands, such as `--context` and `--namespace`.
	KubernetesConfigFlags *genericclioptions.ConfigFlags

	// Provides the `--all-namespaces` or `-A` flag which iterates over all namespaces
	// and adds an extra 'NAMESPACE' header to the output.
	allNameSpaces bool
	dryRun        bool
	VERSION       = "VERSION"
)

func CreateRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "kubectl cleaner",
		Short:         "delete pods with completed status",
		Long:          "delete pods with completed status",
		SilenceErrors: true,
		SilenceUsage:  true,
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlags(cmd.Flags())
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			var successMsg string = "Pods with completed status deleted"
			var err error
			var nameSpace string
			restConfig, err := KubernetesConfigFlags.ToRESTConfig()
			if err != nil {
				return err
			}
			client, err := meta.NewForConfig(restConfig)
			nameSpace, err = kube.GetNamespace(KubernetesConfigFlags, allNameSpaces)
			if err != nil {
				return err
			}
			if err = kube.CleanPodsWithStatus(client, nameSpace, kube.WithDryRun(dryRun), kube.WithPodStatus(kube.POD_STATUS_SUCCEEDED)); err != nil {
				fmt.Println(err)
				return err
			}
			if dryRun {
				successMsg += " (dry run)"
			}
			fmt.Println(successMsg)
			return nil
		},
	}
	cmd.Flags().BoolVarP(&allNameSpaces, "all-namespaces", "A", false, "delete containers that status.phase is succeeded across all namespaces")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Dry run")
	KubernetesConfigFlags = genericclioptions.NewConfigFlags(true)
	KubernetesConfigFlags.AddFlags(cmd.Flags())
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	return cmd
}
