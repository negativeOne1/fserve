package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/martin.kluth1/fserve/cmd/client"
	"gitlab.com/martin.kluth1/fserve/cmd/server"
)

var rootCmd = &cobra.Command{}

func init() {
	rootCmd.AddCommand(printConfigCmd)
	rootCmd.AddCommand(server.RunCmd)
	rootCmd.AddCommand(client.ClientCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
