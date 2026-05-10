package cmd

import (
	"coegle/workspace_cli/core"
	"coegle/workspace_cli/handler"
	"fmt"

	"github.com/spf13/cobra"
)

var addBaseBranch string

var addCmd = &cobra.Command{
	Use:   "add <feat_branch_name> <service_repo_name...>",
	Short: "Add service repositories to a feature workspace",
	Long:  "Add service repositories to a feature workspace. If the workspace does not exist, it will be created automatically.",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := core.GetConfig()
		if err != nil {
			return err
		}
		if cfg.GetReposPath() == "" || cfg.GetWorkspacePath() == "" {
			return fmt.Errorf("please run `ws config base_repos <path>` and `ws config base_ws <path>` first")
		}

		input := core.AddInput{
			BaseInput: core.BaseInput{
				FeatBranchName: args[0],
				Cfg:            cfg,
			},
			ServiceRepoNames: args[1:],
			BaseBranchName:   addBaseBranch,
		}
		return handler.Add(input)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringVarP(&addBaseBranch, "base_branch", "b", "", "The base branch to branch off from (e.g., master, main). If not provided, it detects the default branch.")
}
