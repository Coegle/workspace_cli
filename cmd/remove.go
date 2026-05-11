package cmd

import (
	"github.com/coegle/workspace_cli/core"
	"github.com/coegle/workspace_cli/handler"

	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:     "remove <feat_branch_name> <service_repo_name...>",
	Aliases: []string{"rm"},
	Short:   "Remove a service worktree from the feature workspace",
	Args:    cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		featBranchName := args[0]
		serviceRepoNames := args[1:]

		cfg, err := core.GetConfig()
		if err != nil {
			return err
		}

		input := core.RemoveInput{
			BaseInput: core.BaseInput{
				FeatBranchName: featBranchName,
				Cfg:            cfg,
			},
			ServiceRepoNames: serviceRepoNames,
		}
		handler.Remove(input)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
