package cmd

import (
	"fmt"

	"github.com/ATTron/jitzu/changelog"
	"github.com/spf13/cobra"
)

var changelogRevset string

var changelogCmd = &cobra.Command{
	Use:   "changelog",
	Short: "Generate a markdown changelog from jj history",
	RunE: func(cmd *cobra.Command, args []string) error {
		out, err := changelog.Generate(changelogRevset)
		if err != nil {
			return err
		}

		fmt.Print(out)
		return nil
	},
}

func init() {
	changelogCmd.Flags().StringVarP(&changelogRevset, "revision", "r", "", "revset to generate changelog from")
	rootCmd.AddCommand(changelogCmd)
}
