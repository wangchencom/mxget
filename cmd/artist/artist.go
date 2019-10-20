package artist

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

var CmdArtist = &cobra.Command{
	Use:   "artist",
	Short: "Fetch and download artist hot songs via its id.",
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
	easylog.Infof("Fetch artist %s from %s", id, settings.Site(platform))
	artist, err := client.GetArtist(id)
	if err != nil {
		easylog.Fatal(err)
	}

	cli.ConcurrentDownload(client, artist.Name, artist.Songs...)
}

func init() {
	CmdArtist.Flags().StringVar(&id, "id", "", "artist id")
	CmdArtist.MarkFlagRequired("id")
	CmdArtist.Flags().StringVar(&from, "from", "", "music platform")
	CmdArtist.Flags().IntVar(&settings.Limit, "limit", 0, "concurrent download limit")
	CmdArtist.Flags().BoolVar(&settings.Tag, "tag", false, "update music metadata")
	CmdArtist.Flags().BoolVar(&settings.Lyric, "lyric", false, "download lyric")
	CmdArtist.Flags().BoolVar(&settings.Force, "force", false, "overwrite already downloaded music")
	CmdArtist.Run = Run
}
