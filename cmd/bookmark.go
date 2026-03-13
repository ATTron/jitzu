package cmd

import (
	"fmt"
	"strings"

	"github.com/ATTron/jitzu/jj"
	"github.com/ATTron/jitzu/prompt"
	"github.com/spf13/cobra"
)

var trunkNames = map[string]bool{
	"main":   true,
	"master": true,
	"trunk":  true,
}

type bookmarkInfo struct {
	Revision  string // "@", "@-", "@--", etc.
	Bookmarks []string
}

var bookmarkCmd = &cobra.Command{
	Use:   "bookmark",
	Short: "Intelligently create or move a bookmark for the current change",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, err := detectBookmarkContext()
		if err != nil {
			return fmt.Errorf("detecting bookmark context: %w", err)
		}

		return handleBookmark(ctx)
	},
}

func init() {
	rootCmd.AddCommand(bookmarkCmd)
}

// detectBookmarkContext walks @, @-, @-- looking for bookmarks and
// classifies what situation the user is in.
func detectBookmarkContext() ([]bookmarkInfo, error) {
	revs := []string{"@", "@-", "@--"}
	var infos []bookmarkInfo

	for _, rev := range revs {
		out, err := jj.Log(`bookmarks ++ "\n"`, rev)
		if err != nil {
			// @-- might not exist in a shallow repo
			break
		}
		var bookmarks []string
		for _, b := range strings.Fields(out) {
			// jj appends * to local bookmarks, strip it
			b = strings.TrimRight(b, "*")
			if b != "" {
				bookmarks = append(bookmarks, b)
			}
		}
		infos = append(infos, bookmarkInfo{Revision: rev, Bookmarks: bookmarks})
	}

	return infos, nil
}

func handleBookmark(infos []bookmarkInfo) error {
	// Case 1: Current revision already has a bookmark
	if len(infos) > 0 && len(infos[0].Bookmarks) > 0 {
		fmt.Printf("Current revision already has bookmark(s): %s\n", strings.Join(infos[0].Bookmarks, ", "))

		action, err := prompt.SelectAction("What would you like to do?", []string{
			"Keep as-is",
			"Create an additional bookmark",
		})
		if err != nil {
			return err
		}
		if action == "Keep as-is" {
			return nil
		}
		return createNewBookmark()
	}

	// Case 2: Check ancestors for a nearby non-trunk bookmark
	for _, info := range infos[1:] {
		for _, b := range info.Bookmarks {
			if trunkNames[b] {
				continue
			}
			// Found a feature bookmark on a recent ancestor
			action, err := prompt.SelectAction(
				fmt.Sprintf("Found bookmark %q on %s.", b, info.Revision),
				[]string{
					fmt.Sprintf("Advance %q to here", b),
					fmt.Sprintf("Set %q to here", b),
					"Create a new bookmark instead",
					"Do nothing",
				},
			)
			if err != nil {
				return err
			}

			switch {
			case strings.HasPrefix(action, "Advance"):
				if err := jj.BookmarkAdvance(b); err != nil {
					return fmt.Errorf("advancing bookmark: %w", err)
				}
				fmt.Printf("Bookmark %q advanced to current revision.\n", b)
				return nil
			case strings.HasPrefix(action, "Set"):
				if err := jj.BookmarkSet(b, "@"); err != nil {
					return fmt.Errorf("setting bookmark: %w", err)
				}
				fmt.Printf("Bookmark %q set to current revision.\n", b)
				return nil
			case strings.HasPrefix(action, "Create"):
				return createNewBookmark()
			default:
				return nil
			}
		}
	}

	// Case 3: Nearest bookmark is trunk or no bookmarks found — new feature
	nearest := "unknown"
	for _, info := range infos[1:] {
		for _, b := range info.Bookmarks {
			if trunkNames[b] {
				nearest = b
			}
		}
	}

	if nearest != "unknown" {
		fmt.Printf("Branching off %s — creating a new bookmark.\n", nearest)
	} else {
		fmt.Println("No nearby bookmarks found — creating a new bookmark.")
	}

	return createNewBookmark()
}

func createNewBookmark() error {
	name, err := prompt.BookmarkName(cfg)
	if err != nil {
		return err
	}

	if err := jj.BookmarkCreate(name); err != nil {
		return fmt.Errorf("creating bookmark: %w", err)
	}

	fmt.Printf("Bookmark %q created on current revision.\n", name)
	return nil
}
