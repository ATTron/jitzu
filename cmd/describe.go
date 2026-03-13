package cmd

import (
	"fmt"

	"github.com/ATTron/jitzu/conv"
	"github.com/ATTron/jitzu/jj"
	"github.com/ATTron/jitzu/prompt"
	"github.com/spf13/cobra"
)

var describeRevision string

var describeCmd = &cobra.Command{
	Use:   "describe",
	Short: "Interactively describe a commit",
	RunE: func(cmd *cobra.Command, args []string) error {
		r, err := prompt.Run(cfg)
		if err != nil {
			return err
		}

		msg := conv.Build(r.Type, r.Scope, r.Subject, r.Body, r.Breaking, r.Refs)

		if err := jj.Describe(msg, describeRevision); err != nil {
			return fmt.Errorf("jj describe failed: %w", err)
		}

		fmt.Println("Commit described successfully.")
		return nil
	},
}

func init() {
	describeCmd.Flags().StringVarP(&describeRevision, "revision", "r", "", "revision to describe (default: working copy)")
	rootCmd.AddCommand(describeCmd)
}
