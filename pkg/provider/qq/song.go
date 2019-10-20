package qq

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/winterssy/mxget/pkg/provider"
	"github.com/winterssy/sreq"
)

func GetSong(songMid string) (*provider.Song, error) {
	return std.GetSong(songMid)
}

func (a *API) GetSong(songMid string) (*provider.Song, error) {
	resp, err := a.GetSongRaw(songMid)
	if err != nil {
		return nil, err
	}
	if len(resp.Data) == 0 {
		return nil, errors.New("get song: no data")
	}

	_song := resp.Data[0]
	a.patchSongInfo(_song)
	a.patchSongURL(_song)
	a.patchSongLyric(_song)
	return &provider.Song{
		Name:     strings.TrimSpace(_song.Title),
		Artist:   _song.Artist,
		Album:    strings.TrimSpace(_song.Album.Name),
		PicURL:   _song.PicURL,
		Lyric:    _song.Lyric,
		Playable: _song.URL != "",
		URL:      _song.URL,
	}, nil
}

func GetSongRaw(songMid string) (*SongResponse, error) {
	return std.GetSongRaw(songMid)
}

// 获取歌曲信息
func (a *API) GetSongRaw(songMid string) (*SongResponse, error) {
	params := sreq.Params{
		"songmid": songMid,
	}

	resp := new(SongResponse)
	err := a.Request(sreq.MethodGet, GetSongAPI,
		sreq.WithQuery(params),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf("get song: %d", resp.Code)
	}

	return resp, nil
}

func GetSongURL(songMid string) (string, error) {
	return std.GetSongURL(songMid)
}

func (a *API) GetSongURL(songMid string) (string, error) {
	resp, err := a.GetSongURLRaw(songMid)
	if err != nil {
		return "", err
	}
	if len(resp.Data.Items) == 0 {
		return "", errors.New("get song url: no data")
	}

	item := resp.Data.Items[0]
	if item.SubCode != 0 {
		return "", fmt.Errorf("get song url: %d", item.SubCode)
	}

	return fmt.Sprintf(SongPlayURL, item.FileName, item.Vkey), nil
}

func GetSongURLRaw(songMid string) (*SongURLResponse, error) {
	return std.GetSongURLRaw(songMid)
}

// 获取歌曲播放地址
func (a *API) GetSongURLRaw(songMid string) (*SongURLResponse, error) {
	params := sreq.Params{
		"songmid":  songMid,
		"filename": "M500" + songMid + ".mp3",
	}

	resp := new(SongURLResponse)
	err := a.Request(sreq.MethodGet, GetSongURLAPI,
		sreq.WithQuery(params),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf("get song url: %s", resp.ErrInfo)
	}

	return resp, nil
}

// func GetSongURL(songMid string) (string, error) {
// 	return std.GetSongURL(songMid)
// }
//
// func (a *API) GetSongURL(songMid string) (string, error) {
// 	resp, err := a.GetSongURLRaw(songMid)
// 	if err != nil {
// 		return "", err
// 	}
// 	if len(resp.Req0.Data.MidURLInfo) == 0 {
// 		return "", errors.New("get song url: no data")
// 	}
// 	if len(resp.Req0.Data.Sip) == 0 {
// 		return "", errors.New("get song url: no sip")
// 	}
//
// 	// 随机获取一个sip
// 	sip := resp.Req0.Data.Sip[rand.Intn(len(resp.Req0.Data.Sip))]
// 	urlInfo := resp.Req0.Data.MidURLInfo[0]
// 	if urlInfo.PURL == "" {
// 		return "", errors.New("get song url: copyright protection")
// 	}
// 	return sip + urlInfo.PURL, nil
// }
//
// func GetSongURLRaw(songMids ...string) (*SongURLResponse, error) {
// 	return std.GetSongURLRaw(songMids...)
// }
//
// // 获取歌曲播放地址
// func (a *API) GetSongURLRaw(songMids ...string) (*SongURLResponse, error) {
// 	if len(songMids) > SongURLRequestLimit {
// 		songMids = songMids[:SongURLRequestLimit]
// 	}
//
// 	param := map[string]interface{}{
// 		"guid":      "7332953645",
// 		"loginflag": 1,
// 		"songmid":   songMids,
// 		"uin":       "0",
// 		"platform":  "20",
// 	}
// 	req0 := map[string]interface{}{
// 		"module": "vkey.GetVkeyServer",
// 		"method": "CgiGetVkey",
// 		"param":  param,
// 	}
// 	data := map[string]interface{}{
// 		"req0": req0,
// 	}
//
// 	enc, _ := json.Marshal(data)
// 	params := sreq.Params{
// 		"data": string(enc),
// 	}
// 	resp := new(SongURLResponse)
// 	err := a.request(sreq.MethodGet, GetSongURLAPI,
// 		sreq.WithQuery(params),
// 	).JSON(resp)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if resp.Code != 0 {
// 		return nil, fmt.Errorf("get song url: %d", resp.Code)
// 	}
//
// 	return resp, nil
// }

func GetSongLyric(songMid string) (string, error) {
	return std.GetSongLyric(songMid)
}

func (a *API) GetSongLyric(songMid string) (string, error) {
	resp, err := a.GetSongLyricRaw(songMid)
	if err != nil {
		return "", err
	}

	lyric, err := base64.StdEncoding.DecodeString(resp.Lyric)
	if err != nil {
		return "", err
	}

	return string(lyric), nil
}

func GetSongLyricRaw(songMid string) (*SongLyricResponse, error) {
	return std.GetSongLyricRaw(songMid)
}

// 获取歌词
func (a *API) GetSongLyricRaw(songMid string) (*SongLyricResponse, error) {
	params := sreq.Params{
		"songmid": songMid,
	}

	resp := new(SongLyricResponse)
	err := a.Request(sreq.MethodGet, GetSongLyricAPI,
		sreq.WithQuery(params),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf("get song lyric: %d", resp.Code)
	}

	return resp, nil
}
