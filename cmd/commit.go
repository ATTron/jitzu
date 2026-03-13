package cmd

import (
	"fmt"

	"github.com/ATTron/jitzu/conv"
	"github.com/ATTron/jitzu/jj"
	"github.com/ATTron/jitzu/prompt"
	"github.com/spf13/cobra"
)

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Interactively commit with a conventional message",
	RunE: func(cmd *cobra.Command, args []string) error {
		r, err := prompt.Run(cfg)
		if err != nil {
			return err
		}

		msg := conv.Build(r.Type, r.Scope, r.Subject, r.Body, r.Breaking, r.Refs)

		if err := jj.Commit(msg); err != nil {
			return fmt.Errorf("jj commit failed: %w", err)
		}

		fmt.Println("Commit created successfully.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)
}
