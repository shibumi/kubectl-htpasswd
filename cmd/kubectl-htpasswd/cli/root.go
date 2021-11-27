package cli

import "github.com/spf13/cobra"

type RootOptions struct {
	Verbose bool
}

func (o *RootOptions) AddFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVarP(&o.Verbose, "verbose", "v", false, "log debug output")
}

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "kubectl-htpasswd",
		Short:             "kubectl plugin for generating/managing htpasswd secrets in k8s",
		Long:              "kubectl plugin for generating/managing htpasswd secrets in kubernetes",
		SilenceUsage:      true,
		DisableAutoGenTag: true,
	}

	// add root flags
	ro := &RootOptions{}
	ro.AddFlags(cmd)

	// add subcommands
	cmd.AddCommand(versionCmd)
	cmd.AddCommand(Create())

	return cmd
}
