package kdesc

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tae2089/kubectl-custom-cli/pkg/kube"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var (

	// KubernetesConfigFlags provides the generic flags which are available to
	// regular `kubectl` commands, such as `--context` and `--namespace`.
	KubernetesConfigFlags *genericclioptions.ConfigFlags

	// Provides the `--no-headers` flag, this removes them from being printed to stdout.
	noHeaders bool

	// Provides the `--all-namespaces` or `-A` flag which iterates over all namespaces
	// and adds an extra 'NAMESPACE' header to the output.
	allNamespaces bool

	// Provides the `--version` or `-v` flag, displaying build/version information.
	showVersion bool

	// Provides the `--sort-field` flag, allowing sorting by field.
	// Only 'time' is supported currently.
	sortField string

	// Formatting for table output, similar to other kubectl commands.
	t       = tabwriter.NewWriter(os.Stdout, 10, 1, 5, ' ', 0)
	VERSION = "VERSION"
)

func CreateRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "kubectl desc",
		Short:         "kubectl desc",
		Long:          "kubectl desc",
		SilenceErrors: true,
		SilenceUsage:  true,
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlags(cmd.Flags())
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			var nameSpace string
			var err error
			if showVersion {
				fmt.Printf("%s", VERSION)
				return nil
			}

			nameSpace, err = kube.GetNamespace(KubernetesConfigFlags, allNamespaces)
			if err != nil {
				return err
			}
			fmt.Println(nameSpace, "kubectl desc")
			return nil
		},
	}
	cmd.Flags().StringVar(&sortField, "sort-field", "none", "Sort by particular field. (Only 'time' is supported currently)")
	cmd.Flags().BoolVar(&noHeaders, "no-headers", false, "Don't print headers")
	cmd.Flags().BoolVarP(&allNamespaces, "all-namespaces", "A", false, "Show OOMKilled containers across all namespaces")
	cmd.Flags().BoolVarP(&showVersion, "version", "v", false, "Display version and build information")
	KubernetesConfigFlags = genericclioptions.NewConfigFlags(true)
	KubernetesConfigFlags.AddFlags(cmd.Flags())
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	return cmd
}
