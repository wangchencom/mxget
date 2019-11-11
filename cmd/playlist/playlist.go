package playlist

import (
	"github.com/spf13/cobra"
	"github.com/winterssy/easylog"
	"github.com/winterssy/mxget/internal/cli"
	"github.com/winterssy/mxget/internal/settings"
)

var (
	id   string
	from string
)

var CmdPlaylist = &cobra.Command{
	Use:   "playlist",
	Short: "Fetch and download playlist songs via its id",
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
	easylog.Infof("Fetch playlist %s from %s", id, settings.GetSite(platformId))
	playlist, err := client.GetPlaylist(id)
	if err != nil {
		easylog.Fatal(err)
	}

	cli.ConcurrentDownload(client, playlist.Name, playlist.Songs...)
}

func init() {
	CmdPlaylist.Flags().StringVar(&id, "id", "", "playlist id")
	CmdPlaylist.MarkFlagRequired("id")
	CmdPlaylist.Flags().StringVar(&from, "from", "", "music platform")
	CmdPlaylist.Flags().IntVar(&settings.Limit, "limit", 0, "concurrent download limit")
	CmdPlaylist.Flags().BoolVar(&settings.Tag, "tag", false, "update music metadata")
	CmdPlaylist.Flags().BoolVar(&settings.Lyric, "lyric", false, "download lyric")
	CmdPlaylist.Flags().BoolVar(&settings.Force, "force", false, "overwrite already downloaded music")
	CmdPlaylist.Run = Run
}
