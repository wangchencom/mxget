package baidu

import (
	"context"
	"fmt"
	"strings"

	"github.com/winterssy/mxget/pkg/api"
	"github.com/winterssy/sreq"
)

func (a *API) GetSong(ctx context.Context, songId string) (*api.SongResponse, error) {
	resp, err := a.GetSongRaw(ctx, songId)
	if err != nil {
		return nil, err
	}

	resp.SongInfo.URL = songURL(resp.SongURL.URL)
	a.patchSongsLyric(ctx, &resp.SongInfo)
	songs := resolve(&resp.SongInfo)
	return songs[0], nil
}

// 获取歌曲
func (a *API) GetSongRaw(ctx context.Context, songId string) (*SongResponse, error) {
	resp := new(SongResponse)
	err := a.Request(sreq.MethodGet, APIGetSong, sreq.WithQuery(aesCBCEncrypt(songId))).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.ErrorCode != 22000 {
		return nil, fmt.Errorf("get song: %v", resp.errorMessage())
	}

	return resp, nil
}

// 批量获取歌曲，遗留接口，不推荐使用
func (a *API) GetSongsRaw(ctx context.Context, songIds ...string) (*SongsResponse, error) {
	params := sreq.Params{
		"songIds": strings.Join(songIds, ","),
	}
	resp := new(SongsResponse)
	err := a.Request(sreq.MethodGet, APIGetSongs,
		sreq.WithQuery(params),
		sreq.WithContext(ctx),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.ErrorCode != 22000 {
		return nil, fmt.Errorf("get songs: %d", resp.ErrorCode)
	}

	return resp, nil
}

func (a *API) GetSongLyric(ctx context.Context, songId string) (string, error) {
	resp, err := a.GetSongLyricRaw(ctx, songId)
	if err != nil {
		return "", err
	}

	return resp.LrcContent, nil
}

// 获取歌词
func (a *API) GetSongLyricRaw(ctx context.Context, songId string) (*SongLyricResponse, error) {
	params := sreq.Params{
		"songid": songId,
	}

	resp := new(SongLyricResponse)
	err := a.Request(sreq.MethodGet, APIGetSongLyric,
		sreq.WithQuery(params),
		sreq.WithContext(ctx),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.ErrorCode != 0 && resp.ErrorCode != 22000 {
		return nil, fmt.Errorf("get lyric: %v", resp.errorMessage())
	}

	return resp, nil
}
