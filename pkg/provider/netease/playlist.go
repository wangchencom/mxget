package netease

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/winterssy/mxget/pkg/provider"
	"github.com/winterssy/mxget/pkg/utils"
	"github.com/winterssy/sreq"
)

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

	tracks := resp.Playlist.Tracks
	if n > SongRequestLimit {
		extra := n - SongRequestLimit
		trackIds := make([]int, 0, extra)
		for i := SongRequestLimit; i < n; i++ {
			trackIds = append(trackIds, resp.Playlist.TrackIds[i].Id)
		}

		queue := make(chan []*Song)
		wg := new(sync.WaitGroup)
		for i := 0; i < extra; i += SongRequestLimit {
			songIds := trackIds[i:utils.Min(i+SongRequestLimit, extra)]
			wg.Add(1)
			go func() {
				resp, err := a.GetSongsRaw(songIds...)
				if err != nil {
					wg.Done()
					return
				}
				queue <- resp.Songs
			}()
		}

		go func() {
			for s := range queue {
				resp.Playlist.Tracks = append(tracks, s...)
				wg.Done()
			}
		}()
		wg.Wait()
	}

	a.patchSongURL(SongDefaultBR, tracks...)
	a.patchSongLyric(tracks...)
	songs := resolve(tracks...)
	return &provider.Playlist{
		Id:     strconv.Itoa(resp.Playlist.Id),
		Name:   strings.TrimSpace(resp.Playlist.Name),
		PicURL: resp.Playlist.PicURL,
		Count:  n,
		Songs:  songs,
	}, nil
}

// 获取歌单
func (a *API) GetPlaylistRaw(id int) (*PlaylistResponse, error) {
	data := map[string]int{
		"id": id,
		"n":  100000,
	}

	resp := new(PlaylistResponse)
	err := a.Request(sreq.MethodPost, APIGetPlaylist,
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
