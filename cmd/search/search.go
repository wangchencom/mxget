package search

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/winterssy/easylog"
	"github.com/winterssy/mxget/internal/settings"
)

var (
	keyword string
	from    string
)

var CmdSearch = &cobra.Command{
	Use:   "search",
	Short: "Search song from the Internet.",
}

func Run(cmd *cobra.Command, args []string) {
	platform := settings.Cfg.MusicPlatform
	if from != "" {
		p := settings.Platform(from)
		if !settings.VerifyPlatform(p) {
			easylog.Fatalf("Unexpected music platform: %q", from)
		}
		platform = p
	}

	client := settings.Client(platform)
	result, err := client.SearchSong(keyword)
	if err != nil {
		easylog.Fatal(err)
	}

	for i, s := range result.Songs {
		fmt.Printf("[%02d] %s - %s - %s\n", i+1, s.Name, s.Artist, s.Id)
	}

	if from != "" {
		fmt.Printf(`
Command: 
    mxget song --from %s --id [id]
`, from)
	} else {
		fmt.Print(`
Command: 
    mxget song --id [id]
`)
	}
}

func init() {
	CmdSearch.Flags().StringVarP(&keyword, "keyword", "k", "", "search keyword")
	CmdSearch.MarkFlagRequired("keyword")
	CmdSearch.Flags().StringVar(&from, "from", "", "music platform")
	CmdSearch.Run = Run
}
