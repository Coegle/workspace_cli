package cmd

import (
	"github.com/coegle/workspace_cli/core"
	"github.com/coegle/workspace_cli/pkg/updater"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "ws",
	Short:   "A CLI tool for managing your workspace",
	Long:    `ws is a tool to help you easily manage git worktrees and workspaces.`,
	Version: core.Version,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()

	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	cobra.OnFinalize(updater.CheckAndNotice)

	rootCmd.PersistentFlags().StringVar(&cfgFile, core.PersistentFlagConfig, "", fmt.Sprintf("config file (default is $HOME/%s/%s.%s", core.ConfigPath, core.ConfigFile, core.ConfigType))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(filepath.Join(home, core.ConfigPath))
		viper.SetConfigType(core.ConfigType)
		viper.SetConfigName(core.ConfigFile)
	}

	viper.SetEnvPrefix(core.EnvPrefix)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Error reading config file:", err)
	}
}
