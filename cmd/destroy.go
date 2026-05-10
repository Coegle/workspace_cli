package cmd

import (
	"coegle/workspace_cli/core"
	"coegle/workspace_cli/handler"

	"github.com/spf13/cobra"
)

const forceFlag = "force"

var destroyCmd = &cobra.Command{
	Use:   "destroy <feat_branch_name>",
	Short: "Destroy an entire feature workspace",
	Long: `Destroy an entire feature workspace by removing all of its worktrees
and deleting the workspace directory.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		featBranchName := args[0]
		force, _ := cmd.Flags().GetBool(forceFlag)

		cfg, err := core.GetConfig()
		if err != nil {
			return err
		}

		input := core.DestroyInput{
			BaseInput: core.BaseInput{
				FeatBranchName: featBranchName,
				Cfg:            cfg,
			},
			Force: force,
		}
		handler.Destroy(input)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(destroyCmd)
	destroyCmd.Flags().BoolP(forceFlag, "f", false, "Force destruction without prompting for confirmation")
}
