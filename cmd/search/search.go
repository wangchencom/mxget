package search

import (
	"fmt"
	"strings"

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
	Short: "Search songs from the Internet",
}

func Run(cmd *cobra.Command, args []string) {
	platformId := settings.Cfg.MusicPlatform
	if from != "" {
		pid := settings.GetPlatformId(from)
		if pid == 0 {
			easylog.Fatalf("Unexpected music platform: %q", from)
		}
		platformId = pid
	}

	client := settings.GetClient(platformId)
	fmt.Printf("Search %q from [%s]...\n\n", keyword, settings.GetPlatformDesc(platformId))
	result, err := client.SearchSongs(keyword)
	if err != nil {
		easylog.Fatal(err)
	}

	var sb strings.Builder
	for i, s := range result.Songs {
		fmt.Fprintf(&sb, "[%02d] %s - %s - %s\n", i+1, s.Name, s.Artist, s.Id)
	}
	fmt.Println(sb.String())

	if from != "" {
		fmt.Printf("Command: mxget song --from %s --id [id]\n", from)
	} else {
		fmt.Println("Command: mxget song --id [id]")
	}
}

func init() {
	CmdSearch.Flags().StringVarP(&keyword, "keyword", "k", "", "search keyword")
	CmdSearch.MarkFlagRequired("keyword")
	CmdSearch.Flags().StringVar(&from, "from", "", "music platform")
	CmdSearch.Run = Run
}
