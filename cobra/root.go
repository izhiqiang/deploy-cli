package cobra

import (
	"deploy-cli/cobra/list"
	"deploy-cli/cobra/ping"
	"deploy-cli/cobra/rollback"
	"deploy-cli/cobra/run"
	"deploy-cli/cobra/ssl"
	"deploy-cli/env"
	"deploy-cli/versions"
	"github.com/spf13/cobra"
	"os"
	"path"
	"path/filepath"
	"runtime"
)

var originalWorkingDir = ""
var continueVersionCommands = []string{
	ping.Cmd.Use,
	ssl.Cmd.Use,
}
var rootCmd = &cobra.Command{
	Args:                  cobra.ArbitraryArgs,
	Version:               env.Version,
	DisableAutoGenTag:     true,
	DisableFlagsInUseLine: true,
	SilenceErrors:         true,
	Use:                   "deploy-cli",
	Short:                 "Deploy your code to the server using the CLI command",
	Long: `From development to production, a robust and easy-to-use developer tool
that makes adoption quick and easy for building and deploying  native applications.
`,

	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true
		workDirFlag := cmd.Flags().Lookup("working_dir")
		if workDirFlag != nil && workDirFlag.Value.String() != "" {
			workDir := workDirFlag.Value.String()
			if originalWorkingDir != "" {
				if err = os.Chdir(originalWorkingDir); err != nil {
					return
				}
			}
			if !path.IsAbs(workDir) {
				if workDir, err = filepath.Abs(workDir); err != nil {
					return
				}
			}
			if err = os.Chdir(workDir); err != nil {
				return
			}
			if originalWorkingDir == "" {
				if originalWorkingDir, err = os.Getwd(); err != nil {
					return
				}
			}
			env.Set("PWD", workDir)
		}
		//解决在mac电脑运行在压缩的时候出现`._`
		if runtime.GOOS == "darwin" {
			env.Set("COPYFILE_DISABLE", "1")
		}
		hostGroup := cmd.Flags().Lookup("host_group")
		env.Set(env.HOST_GROUP, hostGroup.Value.String())
		hosts := cmd.Flags().Lookup("hosts")
		env.Set(env.DEPLOY_HOSTS, hosts.Value.String())
		for _, command := range continueVersionCommands {
			if cmd.Use == command {
				return nil
			}
		}
		err = versions.AddLookupPath()
		if err != nil {
			return err
		}
		return
	},
}

func init() {
	rootCmd.PersistentFlags().StringP("working_dir", "W", "", "Changes the working directory for the console")
	rootCmd.PersistentFlags().StringP("host_group", "G", "", "Operation host group name")
	rootCmd.PersistentFlags().StringP("hosts", "H", env.Get(env.DEPLOY_HOSTS), "SSH link hosts, for example: ssh://username:passwd@host:port")
	rootCmd.AddCommand(ping.Cmd)
	rootCmd.AddCommand(run.Cmd)
	rootCmd.AddCommand(list.Cmd)
	rootCmd.AddCommand(rollback.Cmd)
	rootCmd.AddCommand(ssl.Cmd)
}

func Execute() error {
	return rootCmd.Execute()
}
