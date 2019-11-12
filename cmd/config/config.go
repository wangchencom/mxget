package config

import (
	"fmt"
	"os"
	"path/filepath"

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
	Short: "Specify the default behavior of mxget",
}

func Run(cmd *cobra.Command, args []string) {
	if cwd == "" && from == "" {
		fmt.Print(fmt.Sprintf(`
Current settings:
    download dir -> %s
    music platform -> %d [%s]
`, settings.Cfg.DownloadDir, settings.Cfg.MusicPlatform, settings.GetPlatformDesc(settings.Cfg.MusicPlatform)))
		return
	}

	if cwd != "" {
		cwd = filepath.Clean(cwd)
		if err := os.MkdirAll(cwd, 0755); err != nil {
			easylog.Fatalf("Can't make download dir: %v", err)
		}
		settings.Cfg.DownloadDir = cwd
	}
	if from != "" {
		pid := settings.GetPlatformId(from)
		if pid == 0 {
			easylog.Fatalf("Unexpected music platform: %q", from)
		}
		settings.Cfg.MusicPlatform = pid
	}

	_ = settings.Cfg.Save()
}

func init() {
	CmdSet.Flags().StringVar(&cwd, "cwd", "", "specify the default download directory")
	CmdSet.Flags().StringVar(&from, "from", "", "specify the default music platform")
	CmdSet.Run = Run
}
