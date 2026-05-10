package cmd

import (
	"github.com/coegle/workspace_cli/core"
	"github.com/coegle/workspace_cli/handler"

	"github.com/spf13/cobra"
)

var openCmd = &cobra.Command{
	Use:   "open <feat_branch_name>",
	Short: "Open an existing workspace in GoLand",
	Long: `Open a workspace corresponding to the given feature branch name in GoLand.
It will construct the workspace path based on your configuration and launch the IDE.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		featBranchName := args[0]

		cfg, err := core.GetConfig()
		if err != nil {
			return err
		}

		input := core.OpenInput{
			BaseInput: core.BaseInput{
				FeatBranchName: featBranchName,
				Cfg:            cfg,
			},
		}
		return handler.Open(input)
	},
}

func init() {
	rootCmd.AddCommand(openCmd)
}
