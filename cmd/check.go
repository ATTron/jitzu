package cmd

import (
	"fmt"

	"github.com/ATTron/jitzu/conv"
	"github.com/ATTron/jitzu/jj"
	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check [REV]",
	Short: "Validate commit message(s) against conventional format",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		rev := "@"
		if len(args) > 0 {
			rev = args[0]
		}

		out, err := jj.Log(`description.first_line() ++ "\n"`, rev)
		if err != nil {
			return fmt.Errorf("reading commit message: %w", err)
		}

		problems := conv.Validate(out, cfg)
		if len(problems) == 0 {
			fmt.Println("Commit message is valid.")
			return nil
		}

		fmt.Println("Commit message problems:")
		for _, p := range problems {
			fmt.Printf("  - %s\n", p)
		}
		return fmt.Errorf("validation failed with %d problem(s)", len(problems))
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
