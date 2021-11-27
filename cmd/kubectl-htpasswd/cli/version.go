package cli

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	commit = "none"
	date   = "unknown"
	tag    = "dev"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display the version",
	Long:  "Display version of kubectl-htpasswd",
	Run:   version,
}

func version(cmd *cobra.Command, args []string) {
	fmt.Println("commit : ", commit)
	fmt.Println("date   : ", date)
	fmt.Println("version: ", tag)
}
