package list

import (
	"deploy-cli/versions"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:     "list",
	Short:   "Get the dir_name that can be rolled back from the server",
	Example: `deploy-cli list`,
	RunE: func(cmd *cobra.Command, args []string) error {
		v, err := versions.NewVersionsWithVersion()
		if err != nil {
			return err
		}
		return v.List()
	},
}
