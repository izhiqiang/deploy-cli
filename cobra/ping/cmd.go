package ping

import (
	"deploy-cli/versions"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:     "ping",
	Short:   "Check server status",
	Example: `deploy-cli ping`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return versions.NewVersions().Ping()
	},
}
