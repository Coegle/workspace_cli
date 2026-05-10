package cmd

import (
	"github.com/coegle/workspace_cli/core"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configCmd = &cobra.Command{
	Use:   "config [key] [value]",
	Short: "Manage workspace CLI configurations",
	Long: `Print all current configuration settings, or set new ones.

Available keys and their corresponding Environment Variables:
  base_repos      (Env: WS_BASE_REPOS)     The directory where all base service repositories are stored.
  base_ws         (Env: WS_BASE_WS)        The directory where feature workspaces will be created.
  replace_slash   (Env: WS_REPLACE_SLASH)  (true/false) Replace '/' with '_' in branch names when creating workspace directories.

Examples:
  ws config
  ws config base_repos /path/to/repos
  ws config replace_slash false
  WS_REPLACE_SLASH=false ws config`,
	Args: cobra.MaximumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := core.GetConfig()
		if err != nil {
			return err
		}

		if len(args) == 0 {
			// Print all config
			fmt.Printf("Current Configuration:\n")
			fmt.Printf("  base_repos: %s\n", cfg.GetReposPath())
			fmt.Printf("  base_ws:    %s\n", cfg.GetWorkspacePath())
			fmt.Printf("  replace_slash: %v\n", cfg.ReplaceSlashInDir)
			fmt.Printf("\nConfig file loaded from: %s\n", viper.ConfigFileUsed())
			fmt.Printf("Note: You can override these using environment variables (e.g., WS_REPLACE_SLASH=false)\n")
			return nil
		}

		if len(args) == 1 {
			// Get specific config
			key := args[0]
			switch key {
			case "base_repos":
				fmt.Println(cfg.GetReposPath())
			case "base_ws":
				fmt.Println(cfg.GetWorkspacePath())
			case "replace_slash":
				fmt.Println(cfg.ReplaceSlashInDir)
			default:
				return fmt.Errorf("unknown config key: %s", key)
			}
			return nil
		}

		// Set specific config
		key := args[0]
		value := args[1]

		switch key {
		case "base_repos", "base_ws":
			absPath, err := filepath.Abs(value)
			if err != nil {
				return fmt.Errorf("failed to resolve absolute path: %w", err)
			}
			value = absPath

			// Create directory if it doesn't exist
			if err := os.MkdirAll(value, 0755); err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Warning: failed to create directory %s: %v\n", value, err)
			}
		case "replace_slash":
			if value != "true" && value != "false" {
				return fmt.Errorf("value for replace_slash must be 'true' or 'false'")
			}
		default:
			return fmt.Errorf("unknown config key: %s", key)
		}

		// Use viper to update and save
		viper.Set(key, value)

		// Ensure config directory exists before writing
		configDir := filepath.Dir(viper.ConfigFileUsed())
		if configDir == "" || configDir == "." {
			home, _ := os.UserHomeDir()
			configDir = filepath.Join(home, core.ConfigPath)
			_ = os.MkdirAll(configDir, 0755)
			viper.SetConfigFile(filepath.Join(configDir, core.ConfigFile+"."+core.ConfigType))
		} else {
			_ = os.MkdirAll(configDir, 0755)
		}

		if err := viper.WriteConfig(); err != nil {
			// If file doesn't exist, use SafeWriteConfig
			if err := viper.SafeWriteConfig(); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}
		}

		fmt.Printf("Successfully set %s = %s\n", key, value)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
