package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	// Path to the config directory, contains named git configs
	cfgDir        string
	globalProfile bool
	removeProfile bool
	saveProfile   bool
	rootCmd       = &cobra.Command{
		Use:   "gctx",
		Short: "Git context switcher",
		Long: `gctx is a command line tool that helps you 
				switch git context with pre-defined profiles
				includes the user name and email`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			store := NewProfileStore(cfgDir)
			if len(args) == 0 {
				if saveProfile || removeProfile {
					return fmt.Errorf("%s", ProfileNameRequiredMessage())
				}
				return printProfiles(cmd.OutOrStdout(), store, globalProfile)
			}

			profile := args[0]
			switch {
			case saveProfile && removeProfile:
				return fmt.Errorf("%s", ConflictingActionMessage())
			case saveProfile:
				if err := saveCurrentGitProfile(store, profile, globalProfile); err != nil {
					return err
				}
				fmt.Fprintln(cmd.OutOrStdout(), SavedContextMessage(profile))
			case removeProfile:
				if err := store.Remove(profile); err != nil {
					return err
				}
				fmt.Fprintln(cmd.OutOrStdout(), RemovedContextMessage(profile))
			default:
				if err := applyGitProfile(store, profile, globalProfile); err != nil {
					return err
				}
				fmt.Fprintln(cmd.OutOrStdout(), SwitchedContextMessage(profile))
			}
			return nil
		},
	}
)

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgDir, "config", "", "config directory (default is $HOME/.config/gctx)")
	rootCmd.Flags().BoolVarP(&globalProfile, "global", "g", false, "use the global git config")
	rootCmd.Flags().BoolVarP(&removeProfile, "remove", "r", false, "remove the named profile")
	rootCmd.Flags().BoolVarP(&saveProfile, "save", "s", false, "save the current git config as the named profile")
}

func initConfig() {
	if cfgDir == "" {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		cfgDir = filepath.Join(home, ".config", "gctx")
	}

	// Creating a folder to store profiles in =
	err := os.MkdirAll(cfgDir, 0o700)
	cobra.CheckErr(err)
}

func Execute() error {
	return rootCmd.Execute()
}

func printProfiles(out io.Writer, store *ProfileStore, global bool) error {
	current, err := getCurrentGtxProfile(global)
	if err != nil {
		fmt.Fprintln(out, NoActiveContextMessage)
	} else {
		fmt.Fprintln(out, CurrentContextMessage(current))
	}

	profiles := store.List()
	if len(profiles) > 0 {
		fmt.Fprintln(out, strings.Join(profiles, " "))
	}
	return nil
}
