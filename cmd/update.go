package cmd

import (
	"github.com/coegle/workspace_cli/handler"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update ws to the latest version",
	RunE: func(cmd *cobra.Command, args []string) error {
		return handler.Update()
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
