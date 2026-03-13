package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/ATTron/jitzu/config"
	"github.com/ATTron/jitzu/conv"
	"github.com/ATTron/jitzu/jj"
	"github.com/ATTron/jitzu/prompt"
	"github.com/spf13/cobra"
)

var cfg config.Config

var rootRevision string

var rootCmd = &cobra.Command{
	Use:           "jitzu",
	Short:         "Commitizen for Jujutsu — interactive conventional commits for jj",
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Name() == "init" {
			return nil
		}
		if _, err := exec.LookPath("jj"); err != nil {
			return fmt.Errorf("jj is not installed or not in PATH")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		r, err := prompt.Run(cfg)
		if err != nil {
			return err
		}

		msg := conv.Build(r.Type, r.Scope, r.Subject, r.Body, r.Breaking, r.Refs)

		if err := jj.Describe(msg, rootRevision); err != nil {
			return fmt.Errorf("jj describe failed: %w", err)
		}

		fmt.Println("Commit described successfully.")
		return nil
	},
}

func init() {
	rootCmd.Flags().StringVarP(&rootRevision, "revision", "r", "", "revision to describe (default: working copy)")
}

func Execute() {
	cfg = config.Load()

	if err := rootCmd.Execute(); err != nil {
		if errors.Is(err, prompt.ErrAborted) {
			os.Exit(0)
		}
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
