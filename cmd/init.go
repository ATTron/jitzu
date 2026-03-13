package cmd

import (
	"fmt"
	"os"

	"github.com/ATTron/jitzu/jj"
	"github.com/spf13/cobra"
)

var installAlias bool

const defaultConfig = `# jitzu configuration
# See https://github.com/ATTron/jitzu for documentation

# Uncomment to customize commit types:
# [[types]]
# name = "feat"
# description = "A new feature"

# Uncomment to restrict scopes:
# scopes = ["api", "ui", "core"]

# scope_required = false
# body_required = false
# subject_max_len = 72
# body_max_len = 0
`

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a default .jitzu.toml and print alias setup instructions",
	RunE: func(cmd *cobra.Command, args []string) error {
		const filename = ".jitzu.toml"

		if _, err := os.Stat(filename); err == nil {
			return fmt.Errorf("%s already exists", filename)
		}

		if err := os.WriteFile(filename, []byte(defaultConfig), 0644); err != nil {
			return fmt.Errorf("writing %s: %w", filename, err)
		}

		fmt.Printf("Created %s\n\n", filename)

		if installAlias {
			_, err := jj.Run("config", "set", "--user", "aliases.z", `["util", "exec", "--", "jitzu", "describe"]`)
			if err != nil {
				return fmt.Errorf("setting jj alias: %w", err)
			}
			fmt.Println("Alias installed! You can now use: jj z")
		} else {
			fmt.Println("To set up the jj alias, add this to ~/.jjconfig.toml:")
			fmt.Println()
			fmt.Println(`[aliases]`)
			fmt.Println(`z = ["util", "exec", "--", "jitzu", "describe"]`)
			fmt.Println()
			fmt.Println("Or run: jitzu init --install-alias")
		}

		return nil
	},
}

func init() {
	initCmd.Flags().BoolVar(&installAlias, "install-alias", false, "automatically install the jj alias")
	rootCmd.AddCommand(initCmd)
}
