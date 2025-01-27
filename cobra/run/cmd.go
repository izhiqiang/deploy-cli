package run

import (
	"deploy-cli/versions"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:     "run",
	Short:   "Publish configuration information to the server through `.deploy.yml`",
	Example: `deploy-cli run`,
	RunE: func(cmd *cobra.Command, args []string) error {
		v, err := versions.NewVersionsWithVersion()
		if err != nil {
			return err
		}
		return v.Run()
	},
}
