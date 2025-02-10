package rollback

import (
	"deploy-cli/versions"
	"github.com/spf13/cobra"
)

var (
	name string
)

func init() {
	Cmd.Flags().StringVarP(&name, "name", "N", "", "Directory Name")
}

var Cmd = &cobra.Command{
	Use:     "rollback",
	Short:   "rollback the code to the specified version",
	Example: `deploy-cli rollback --name test-20250106152443`,
	RunE: func(cmd *cobra.Command, args []string) error {
		v, err := versions.NewVersionsWithVersion()
		if err != nil {
			return err
		}
		return v.Rollback(name)
	},
}
