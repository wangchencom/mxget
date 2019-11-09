package qq

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"

	"github.com/winterssy/mxget/pkg/provider"
	"github.com/winterssy/sreq"
)

const (
	// 跟vkey配合获取歌曲下载地址，可为任意字符串
	Guid = "0"
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
	a.patchSongURLV1(_song)
	a.patchSongLyric(_song)
	songs := resolve(_song)
	return songs[0], nil
}

func GetSongRaw(songMid string) (*SongResponse, error) {
	return std.GetSongRaw(songMid)
}

// 获取歌曲详情
func (a *API) GetSongRaw(songMid string) (*SongResponse, error) {
	params := sreq.Params{
		"songmid": songMid,
	}

	resp := new(SongResponse)
	err := a.Request(sreq.MethodGet, APIGetSong,
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

func GetSongURLV1(songMid string, mediaMid string) (string, error) {
	return std.GetSongURLV1(songMid, mediaMid)
}

func (a *API) GetSongURLV1(songMid string, mediaMid string) (string, error) {
	resp, err := a.GetSongURLV1Raw(songMid, mediaMid)
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

	return fmt.Sprintf(SongURL, item.FileName, item.Vkey), nil
}

func GetSongURLV1Raw(songMid string, mediaMid string) (*SongURLResponseV1, error) {
	return std.GetSongURLV1Raw(songMid, mediaMid)
}

// 获取歌曲播放地址
func (a *API) GetSongURLV1Raw(songMid string, mediaMid string) (*SongURLResponseV1, error) {
	params := sreq.Params{
		"songmid":  songMid,
		"filename": "M500" + mediaMid + ".mp3",
	}

	resp := new(SongURLResponseV1)
	err := a.Request(sreq.MethodGet, APIGetSongURLV1,
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

func GetSongURLV2(songMid string) (string, error) {
	return std.GetSongURLV2(songMid)
}

func (a *API) GetSongURLV2(songMid string) (string, error) {
	resp, err := a.GetSongsURLV2Raw(songMid)
	if err != nil {
		return "", err
	}
	if len(resp.Req0.Data.MidURLInfo) == 0 {
		return "", errors.New("get song url: no data")
	}
	if len(resp.Req0.Data.Sip) == 0 {
		return "", errors.New("get song url: no sip")
	}

	// 随机获取一个sip
	sip := resp.Req0.Data.Sip[rand.Intn(len(resp.Req0.Data.Sip))]
	urlInfo := resp.Req0.Data.MidURLInfo[0]
	if urlInfo.PURL == "" {
		return "", errors.New("get song url: copyright protection")
	}
	return sip + urlInfo.PURL, nil
}

func GetSongsURLV2Raw(songMids ...string) (*SongURLResponseV2, error) {
	return std.GetSongsURLV2Raw(songMids...)
}

// 批量获取歌曲播放地址
func (a *API) GetSongsURLV2Raw(songMids ...string) (*SongURLResponseV2, error) {
	if len(songMids) > SongURLRequestLimit {
		songMids = songMids[:SongURLRequestLimit]
	}

	param := map[string]interface{}{
		"guid":      Guid,
		"loginflag": 1,
		"songmid":   songMids,
		"uin":       "0",
		"platform":  "20",
	}
	req0 := map[string]interface{}{
		"module": "vkey.GetVkeyServer",
		"method": "CgiGetVkey",
		"param":  param,
	}
	data := map[string]interface{}{
		"req0": req0,
	}

	enc, _ := json.Marshal(data)
	params := sreq.Params{
		"data": string(enc),
	}
	resp := new(SongURLResponseV2)
	err := a.Request(sreq.MethodGet, APIGetSongURLV2,
		sreq.WithQuery(params),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf("get song url: %d", resp.Code)
	}

	return resp, nil
}

func GetSongLyric(songMid string) (string, error) {
	return std.GetSongLyric(songMid)
}

func (a *API) GetSongLyric(songMid string) (string, error) {
	resp, err := a.GetSongLyricRaw(songMid)
	if err != nil {
		return "", err
	}

	// lyric, err := base64.StdEncoding.DecodeString(resp.Lyric)
	// if err != nil {
	// 	return "", err
	// }

	return resp.Lyric, nil
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
	err := a.Request(sreq.MethodGet, APIGetSongLyric,
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
