package cmd

import (
	"coegle/workspace_cli/core"
	"coegle/workspace_cli/handler"

	"github.com/spf13/cobra"
)

var syncBaseBranch string

var syncCmd = &cobra.Command{
	Use:   "sync <service_repo_name...>",
	Short: "Sync the base branch of service repositories with origin",
	Long: `Sync updates the specified base branch (default: master) of the given service repositories
by fetching from origin. It uses fast-forward pull if the branch is currently checked out,
otherwise it fetches the remote branch directly into the local branch pointer.`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		serviceRepoNames := args

		cfg, err := core.GetConfig()
		if err != nil {
			return err
		}

		input := core.SyncInput{
			ServiceRepoNames: serviceRepoNames,
			BaseBranchName:   syncBaseBranch,
			Cfg:              cfg,
		}
		handler.Sync(input)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
	syncCmd.Flags().StringVarP(&syncBaseBranch, "base_branch", "b", "", "The base branch to sync (e.g., master, main). If not provided, it detects the default branch.")
}
