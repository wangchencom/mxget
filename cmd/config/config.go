package config

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/winterssy/easylog"
	"github.com/winterssy/mxget/internal/settings"
)

var (
	cwd  string
	from string
)

var CmdSet = &cobra.Command{
	Use:   "config",
	Short: "Specify the default behavior of mxget.",
}

func Run(cmd *cobra.Command, args []string) {
	if cwd == "" && from == "" {
		fmt.Printf("download dir: %s\n", settings.Cfg.DownloadDir)
		fmt.Printf("music platform: %d\n", settings.Cfg.MusicPlatform)
		return
	}

	if cwd != "" {
		if err := os.MkdirAll(cwd, 0755); err != nil {
			easylog.Fatalf("Failed to make download dir: %q: %v", cwd, err)
		}
		settings.Cfg.DownloadDir = cwd
	}
	if from != "" {
		platform := settings.Platform(from)
		if !settings.VerifyPlatform(platform) {
			easylog.Fatalf("Unexpected music platform: %q", from)
		}
		settings.Cfg.MusicPlatform = platform
	}

	_ = settings.Cfg.Save()
}

func init() {
	CmdSet.Flags().StringVar(&cwd, "cwd", "", "specify the default download directory")
	CmdSet.Flags().StringVar(&from, "from", "", "specify the default music platform")
	CmdSet.Run = Run
}
