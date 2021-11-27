package cli

import (
	"errors"
	"github.com/shibumi/kubectl-htpasswd/internal"
	"github.com/shibumi/kubectl-htpasswd/pkg/htpasswd"
	"github.com/spf13/cobra"
	"strings"
)

type genOptions struct {
	DryRun    bool
	Algorithm string
	Cost      int
	Namespace string
	Format    string
	Key       string
}

func (o *genOptions) AddFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&o.DryRun, "dry-run", "", false, "print the k8s secret to stdout without creating it on the cluster")
	cmd.Flags().StringVarP(&o.Algorithm, "algorithm", "a", "bcrypt", "select the hash algorithm. Can be one out of ['bcrypt']")
	cmd.Flags().IntVarP(&o.Cost, "cost", "c", 10, "select the hash algorithm cost. Must be between 4 and 31")
	cmd.Flags().StringVarP(&o.Namespace, "namespace", "n", "", "select the target namespace for the k8s secret")
	cmd.Flags().StringVarP(&o.Format, "output", "o", "", "output format. Can be one of ['json','yaml']")
	cmd.Flags().StringVarP(&o.Key, "key", "k", "auth", "key in the kubernetes secret data object")
}

func Create() *cobra.Command {
	o := &genOptions{}
	cmd := &cobra.Command{
		Use:   "create [secretName] [user=password]...",
		Short: "create a htpasswd secret in kubernetes",
		Long:  "create a htpasswd secret in kubernetes",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New("foo")
			}
			var secretName string
			var entries []string
			for i, arg := range args {
				if i == 0 {
					secretName = arg
					continue
				}
				pair := strings.SplitN(arg, "=", 2) // split in two, this allows passwords with "="
				res, err := htpasswd.BuildEntry(pair[0], pair[1], o.Algorithm, o.Cost)
				if err != nil {
					return err
				}
				entries = append(entries, res)
			}
			data := strings.Join(entries, "\n")
			client, err := internal.NewClient(o.DryRun, o.Namespace, secretName, o.Format, o.Key, []byte(data))
			if err != nil {
				return err
			}
			err = client.Create()
			if err != nil {
				return err
			}
			return nil
		},
	}
	o.AddFlags(cmd)
	return cmd
}
