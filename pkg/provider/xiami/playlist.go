package xiami

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

func GetPlaylist(playlistId string) (*provider.Playlist, error) {
	return std.GetPlaylist(playlistId)
}

func (a *API) GetPlaylist(playlistId string) (*provider.Playlist, error) {
	resp, err := a.GetPlaylistDetailRaw(playlistId, 1, SongRequestLimit)
	if err != nil {
		return nil, err
	}

	n, _ := strconv.Atoi(resp.Data.Data.CollectDetail.SongCount)
	if n == 0 {
		return nil, errors.New("get playlist: no data")
	}

	_songs := resp.Data.Data.CollectDetail.Songs
	if n > SongRequestLimit {
		allSongs := resp.Data.Data.CollectDetail.AllSongs
		queue := make(chan []*Song)
		wg := new(sync.WaitGroup)
		for i := SongRequestLimit; i < n; i += SongRequestLimit {
			songIds := allSongs[i:utils.Min(i+SongRequestLimit, n)]
			wg.Add(1)
			go func() {
				resp, err := a.GetSongsRaw(songIds...)
				if err != nil {
					wg.Done()
					return
				}
				queue <- resp.Data.Data.Songs
			}()
		}

		go func() {
			for s := range queue {
				_songs = append(_songs, s...)
				wg.Done()
			}
		}()
		wg.Wait()
	}

	a.patchSongLyric(_songs...)
	songs := resolve(_songs...)
	return &provider.Playlist{
		Name:   strings.TrimSpace(resp.Data.Data.CollectDetail.CollectName),
		PicURL: resp.Data.Data.CollectDetail.CollectLogo,
		Count:  n,
		Songs:  songs,
	}, nil
}

func GetPlaylistDetailRaw(playlistId string, page int, pageSize int) (*PlaylistDetailResponse, error) {
	return std.GetPlaylistDetailRaw(playlistId, page, pageSize)
}

// 获取歌单详情，包含歌单信息跟歌曲
func (a *API) GetPlaylistDetailRaw(playlistId string, page int, pageSize int) (*PlaylistDetailResponse, error) {
	token, err := a.getToken(APIGetPlaylistDetail)
	if err != nil {
		return nil, err
	}

	model := map[string]interface{}{
		"listId": playlistId,
		"pagingVO": map[string]int{
			"page":     page,
			"pageSize": pageSize,
		},
	}
	params := sreq.Params(signPayload(token, model))
	resp := new(PlaylistDetailResponse)
	err = a.Request(sreq.MethodGet, APIGetPlaylistDetail, sreq.WithQuery(params)).JSON(resp)
	if err != nil {
		return nil, err
	}
	err = resp.check()
	if err != nil {
		return nil, fmt.Errorf("get playlist detail: %w", err)
	}

	return resp, nil
}

func GetPlaylistSongsRaw(playlistId string, page int, pageSize int) (*PlaylistSongsResponse, error) {
	return std.GetPlaylistSongsRaw(playlistId, page, pageSize)
}

// 获取歌单歌曲，不包含歌单信息
func (a *API) GetPlaylistSongsRaw(playlistId string, page int, pageSize int) (*PlaylistSongsResponse, error) {
	token, err := a.getToken(APIGetPlaylistSongs)
	if err != nil {
		return nil, err
	}

	model := map[string]interface{}{
		"listId": playlistId,
		"pagingVO": map[string]int{
			"page":     page,
			"pageSize": pageSize,
		},
	}
	params := sreq.Params(signPayload(token, model))
	resp := new(PlaylistSongsResponse)
	err = a.Request(sreq.MethodGet, APIGetPlaylistSongs, sreq.WithQuery(params)).JSON(resp)
	if err != nil {
		return nil, err
	}
	err = resp.check()
	if err != nil {
		return nil, fmt.Errorf("get playlist songs: %w", err)
	}

	return resp, nil
}
