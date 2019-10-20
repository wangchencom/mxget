package cli

import (
	"fmt"
	"github.com/bogem/id3v2"
	"github.com/winterssy/easylog"
	"github.com/winterssy/mxget/internal/settings"
	"github.com/winterssy/mxget/pkg/concurrency"
	"github.com/winterssy/mxget/pkg/provider"
	"github.com/winterssy/mxget/pkg/utils"
	"github.com/winterssy/sreq"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
)

func ConcurrentDownload(client provider.API, savePath string, songs ...*provider.Song) {
	savePath = filepath.Join(settings.Cfg.DownloadDir, utils.TrimInvalidFilePathChars(savePath))
	if err := os.MkdirAll(savePath, 0755); err != nil {
		easylog.Fatal(err)
	}

	var limit int
	switch {
	case settings.Limit < 1:
		limit = runtime.NumCPU()
	case settings.Limit > 32:
		limit = 32
	default:
		limit = settings.Limit
	}

	c := concurrency.New(limit)
	for _, s := range songs {
		c.Add(1)
		go func(s *provider.Song) {
			defer c.Done()
			songInfo := fmt.Sprintf("%s - %s", s.Artist, s.Name)
			if !s.Playable {
				easylog.Errorf("Download [%s] failed: song unavailable", songInfo)
				return
			}

			filePath := filepath.Join(savePath, utils.TrimInvalidFilePathChars(songInfo))
			easylog.Infof("Start download [%s]", songInfo)
			mp3FilePath := filePath + ".mp3"
			if !settings.Force {
				_, err := os.Stat(mp3FilePath)
				if err == nil {
					easylog.Infof("Song already downloaded: [%s]", songInfo)
					return
				}
			}

			err := client.Request(sreq.MethodGet, s.URL).Save(mp3FilePath)
			if err != nil {
				easylog.Errorf("Download [%s] failed: %v", songInfo, err)
				_ = os.Remove(mp3FilePath)
				return
			}
			easylog.Infof("Download [%s] complete", songInfo)

			if settings.Tag {
				easylog.Infof("Update music metadata: [%s]", songInfo)
				writeTag(client, mp3FilePath, s)
			}

			if settings.Lyric && s.Lyric != "" {
				easylog.Infof("Save lyric: [%s]", songInfo)
				lrcFilePath := filePath + ".lrc"
				saveLyric(lrcFilePath, s.Lyric)
			}
		}(s)
	}
	c.Wait()
}

func saveLyric(filePath string, lyric string) {
	err := ioutil.WriteFile(filePath, []byte(lyric), 0644)
	if err != nil {
		_ = os.Remove(filePath)
	}
}

func writeTag(client provider.API, filePath string, song *provider.Song) {
	tag, err := id3v2.Open(filePath, id3v2.Options{Parse: true})
	if err != nil {
		return
	}
	defer tag.Close()

	tag.SetDefaultEncoding(id3v2.EncodingUTF8)
	tag.SetTitle(song.Name)
	tag.SetArtist(song.Artist)
	tag.SetAlbum(song.Album)

	if song.Lyric != "" {
		uslt := id3v2.UnsynchronisedLyricsFrame{
			Encoding:          id3v2.EncodingUTF8,
			Language:          "eng",
			ContentDescriptor: song.Name,
			Lyrics:            song.Lyric,
		}
		tag.AddUnsynchronisedLyricsFrame(uslt)
	}

	if song.PicURL != "" {
		pic, err := client.Request(sreq.MethodGet, song.PicURL).Raw()
		if err == nil {
			picFrame := id3v2.PictureFrame{
				Encoding:    id3v2.EncodingUTF8,
				MimeType:    "image/jpeg",
				PictureType: id3v2.PTFrontCover,
				Description: "Front cover",
				Picture:     pic,
			}
			tag.AddAttachedPicture(picFrame)
		}
	}

	_ = tag.Save()
}
