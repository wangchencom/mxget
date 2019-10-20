package song

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

var CmdSong = &cobra.Command{
	Use:   "song",
	Short: "Fetch and download song via its id.",
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
	easylog.Infof("Fetch song %s from %s", id, settings.Site(platform))
	song, err := client.GetSong(id)
	if err != nil {
		easylog.Fatal(err)
	}

	cli.ConcurrentDownload(client, ".", song)
}

func init() {
	CmdSong.Flags().StringVar(&id, "id", "", "song id")
	CmdSong.MarkFlagRequired("id")
	CmdSong.Flags().StringVar(&from, "from", "", "music platform")
	CmdSong.Flags().BoolVar(&settings.Tag, "tag", false, "update music metadata")
	CmdSong.Flags().BoolVar(&settings.Lyric, "lyric", false, "download lyric")
	CmdSong.Flags().BoolVar(&settings.Force, "force", false, "overwrite already downloaded music")
	CmdSong.Run = Run
}
