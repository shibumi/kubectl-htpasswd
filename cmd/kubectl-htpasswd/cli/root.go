package cli

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var logger *zap.Logger

type RootOptions struct {
	Verbose bool
}

func (o *RootOptions) AddFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVarP(&o.Verbose, "verbose", "v", false, "log debug output")
}

func New() *cobra.Command {
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}

	cmd := &cobra.Command{
		Use:               "kubectl-htpasswd",
		Short:             "kubectl plugin for generating htpasswd secrets in k8s",
		Long:              "kubectl plugin for generating htpasswd secrets in kubernetes",
		SilenceUsage:      true,
		DisableAutoGenTag: true,
		// PersistentPreRun sets a development logger after flag evaluation.
		// This way we are able to choose the debugging level via verbose.
		// Verbose will log on debug level.
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			verbose, err := cmd.Parent().PersistentFlags().GetBool("verbose")
			if err != nil {
				panic(err)
			}
			if verbose {
				logger, err = zap.NewDevelopment()
				if err != nil {
					panic(err)
				}
				logger.Debug("Set log level to debug")
			}
		},
	}

	// add root flags
	ro := &RootOptions{}
	ro.AddFlags(cmd)
	// add subcommands
	cmd.AddCommand(versionCmd)
	cmd.AddCommand(Create())

	return cmd
}
