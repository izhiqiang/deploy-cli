package main

import (
	"deploy-cli/cobra"
	"deploy-cli/logger"
	"os"
)

func main() {

	if err := cobra.Execute(); err != nil {
		logger.Error(err)
		os.Exit(1)
	}
	os.Exit(0)
}
