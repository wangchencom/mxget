package netease

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/winterssy/mxget/pkg/provider"
	"github.com/winterssy/sreq"
)

func GetPlaylist(playlistId string) (*provider.Playlist, error) {
	return std.GetPlaylist(playlistId)
}

func (a *API) GetPlaylist(playlistId string) (*provider.Playlist, error) {
	id, err := strconv.Atoi(playlistId)
	if err != nil {
		return nil, err
	}

	resp, err := a.GetPlaylistRaw(id)
	if err != nil {
		return nil, err
	}

	n := resp.Playlist.Total
	if n == 0 {
		return nil, errors.New("get playlist: no data")
	}

	ids := make([]int, 0, n)
	for _, t := range resp.Playlist.TrackIds {
		ids = append(ids, t.Id)
	}

	if n > SongRequestLimit {
		queue := make(chan []*Song)
		wg := new(sync.WaitGroup)
		for i := SongRequestLimit; i < n; i += SongRequestLimit {
			j := i + SongRequestLimit
			if j > n {
				j = n
			}

			ids := make([]int, 0, j-i)
			for k := i; k < j; k++ {
				ids = append(ids, resp.Playlist.TrackIds[k].Id)
			}

			wg.Add(1)
			go func() {
				resp, _ := a.GetSongRaw(ids...)
				queue <- resp.Songs
			}()
		}

		go func() {
			for s := range queue {
				if len(s) != 0 {
					resp.Playlist.Tracks = append(resp.Playlist.Tracks, s...)
				}
				wg.Done()
			}
		}()
		wg.Wait()
	}

	a.patchSongURL(SongDefaultBR, resp.Playlist.Tracks...)
	a.patchSongLyric(resp.Playlist.Tracks...)
	songs := a.resolve(resp.Playlist.Tracks...)
	return &provider.Playlist{
		Name:   strings.TrimSpace(resp.Playlist.Name),
		PicURL: resp.Playlist.PicURL,
		Count:  n,
		Songs:  songs,
	}, nil
}

func GetPlaylistRaw(id int) (*PlaylistResponse, error) {
	return std.GetPlaylistRaw(id)
}

// 获取歌单
func (a *API) GetPlaylistRaw(id int) (*PlaylistResponse, error) {
	data := map[string]int{
		"id": id,
		"n":  100000,
	}

	resp := new(PlaylistResponse)
	err := a.Request(sreq.MethodPost, GetPlaylistAPI,
		sreq.WithForm(weapi(data)),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != 200 {
		return nil, fmt.Errorf("get playlist: %s", resp.Msg)
	}

	return resp, nil
}
