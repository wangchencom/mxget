package album

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

var CmdAlbum = &cobra.Command{
	Use:   "album",
	Short: "Fetch and download album songs via its id.",
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
	easylog.Infof("Fetch album %s from %s", id, settings.Site(platform))
	album, err := client.GetAlbum(id)
	if err != nil {
		easylog.Fatal(err)
	}

	cli.ConcurrentDownload(client, album.Name, album.Songs...)
}

func init() {
	CmdAlbum.Flags().StringVar(&id, "id", "", "album id")
	CmdAlbum.MarkFlagRequired("id")
	CmdAlbum.Flags().StringVar(&from, "from", "", "music platform")
	CmdAlbum.Flags().IntVar(&settings.Limit, "limit", 0, "concurrent download limit")
	CmdAlbum.Flags().BoolVar(&settings.Tag, "tag", false, "update music metadata")
	CmdAlbum.Flags().BoolVar(&settings.Lyric, "lyric", false, "download lyric")
	CmdAlbum.Flags().BoolVar(&settings.Force, "force", false, "overwrite already downloaded music")
	CmdAlbum.Run = Run
}
