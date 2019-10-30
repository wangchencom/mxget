package netease

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/winterssy/mxget/pkg/provider"
	"github.com/winterssy/sreq"
)

func GetSong(songId string) (*provider.Song, error) {
	return std.GetSong(songId)
}

func (a *API) GetSong(songId string) (*provider.Song, error) {
	id, err := strconv.Atoi(songId)
	if err != nil {
		return nil, err
	}

	resp, err := a.GetSongRaw(id)
	if err != nil {
		return nil, err
	}
	if len(resp.Songs) == 0 {
		return nil, errors.New("get song: no data")
	}

	_song := resp.Songs[0]
	a.patchSongURL(SongDefaultBR, _song)
	a.patchSongLyric(_song)
	songs := a.resolve(_song)
	return songs[0], nil
}

func GetSongRaw(ids ...int) (*SongResponse, error) {
	return std.GetSongRaw(ids...)
}

// 获取歌曲信息
func (a *API) GetSongRaw(ids ...int) (*SongResponse, error) {
	n := len(ids)
	if n > SongRequestLimit {
		ids = ids[:SongRequestLimit]
		n = SongRequestLimit
	}

	c := make([]map[string]int, 0, n)
	for _, id := range ids {
		c = append(c, map[string]int{"id": id})
	}
	enc, _ := json.Marshal(c)
	data := map[string]string{
		"c": string(enc),
	}

	resp := new(SongResponse)
	err := a.Request(sreq.MethodPost, GetSongAPI,
		sreq.WithForm(weapi(data)),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != 200 {
		return nil, fmt.Errorf("get song: %s", resp.Msg)
	}

	return resp, nil
}

func GetSongURL(id int, br int) (string, error) {
	return std.GetSongURL(id, br)
}

func (a *API) GetSongURL(id int, br int) (string, error) {
	resp, err := a.GetSongURLRaw(br, id)
	if err != nil {
		return "", err
	}
	if len(resp.Data) == 0 {
		return "", errors.New("get song url: no data")
	}
	if resp.Data[0].Code != 200 {
		return "", errors.New("get song url: copyright protection")
	}

	return resp.Data[0].URL, nil
}

func GetSongURLRaw(br int, ids ...int) (*SongURLResponse, error) {
	return std.GetSongURLRaw(br, ids...)
}

// 获取歌曲播放地址，br: 比特率，128/192/320/999
func (a *API) GetSongURLRaw(br int, ids ...int) (*SongURLResponse, error) {
	var _br int
	switch br {
	case 128, 192, 320:
		_br = br
	default:
		_br = 999
	}

	enc, _ := json.Marshal(ids)
	data := map[string]interface{}{
		"br":  _br * 1000,
		"ids": string(enc),
	}

	resp := new(SongURLResponse)
	err := a.Request(sreq.MethodPost, GetSongURLAPI,
		sreq.WithForm(weapi(data)),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != 200 {
		return nil, fmt.Errorf("get song url: %s", resp.Msg)
	}

	return resp, nil
}

func GetSongLyric(id int) (string, error) {
	return std.GetSongLyric(id)
}

func (a *API) GetSongLyric(id int) (string, error) {
	resp, err := a.GetSongLyricRaw(id)
	if err != nil {
		return "", err
	}
	return resp.Lrc.Lyric, nil
}

func GetSongLyricRaw(id int) (*SongLyricResponse, error) {
	return std.GetSongLyricRaw(id)
}

// 获取歌词
func (a *API) GetSongLyricRaw(id int) (*SongLyricResponse, error) {
	data := map[string]interface{}{
		"method": "POST",
		"url":    "https://music.163.com/api/song/lyric?lv=-1&kv=-1&tv=-1",
		"params": map[string]int{
			"id": id,
		},
	}
	resp := new(SongLyricResponse)
	err := a.Request(sreq.MethodPost, LinuxAPI,
		sreq.WithForm(linuxapi(data)),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != 200 {
		return nil, fmt.Errorf("get song lyric: %s", resp.Msg)
	}

	return resp, nil
}
