package migu

import (
	"errors"
	"fmt"
	"strings"

	"github.com/winterssy/mxget/pkg/provider"
	"github.com/winterssy/sreq"
)

/*
	注意！
	GetSongIdRaw, GetSongPicRaw, GetSongLyricRaw 为网页版API
	这些API限流，并发请求经常503，不适用于批量获取
*/

func (a *API) GetSongId(copyrightId string) (string, error) {
	resp, err := a.GetSongIdRaw(copyrightId)
	if err != nil {
		return "", err
	}
	if len(resp.Items) == 0 || resp.Items[0].SongId == "" {
		return "", errors.New("get song id: no data")
	}

	return resp.Items[0].SongId, nil
}

// 根据版权id获取歌曲id
func (a *API) GetSongIdRaw(copyrightId string) (*SongIdResponse, error) {
	params := sreq.Params{
		"copyrightId": copyrightId,
	}

	resp := new(SongIdResponse)
	err := a.Request(sreq.MethodGet, APIGetSongId,
		sreq.WithQuery(params),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.ReturnCode != "000000" {
		return nil, fmt.Errorf("get song id: %s", resp.Msg)
	}

	return resp, nil
}

func (a *API) GetSong(id string) (*provider.Song, error) {
	/*
		先判断id是版权id还是歌曲id，以减少1次API请求
		测试表现版权id的长度是11位，以6开头并且可能包含字符，歌曲id为纯数字，长度不定
		不确定是否会误判，待反馈
	*/
	var songId string
	var err error
	if len(id) > 10 && strings.HasPrefix(id, "6") {
		songId, err = a.GetSongId(id)
		if err != nil {
			return nil, err
		}
	} else {
		songId = id
	}

	resp, err := a.GetSongRaw(songId)
	if err != nil {
		return nil, err
	}
	if len(resp.Resource) == 0 {
		return nil, errors.New("get song: no data")
	}

	_song := resp.Resource[0]

	// 单曲请求可调用网页版API获取歌词，不会出现乱码现象
	lyric, err := a.GetSongLyric(_song.CopyrightId)
	if err == nil {
		_song.Lyric = lyric
	}
	songs := resolve(_song)
	return songs[0], nil
}

// 获取歌曲详情
func (a *API) GetSongRaw(songId string) (*SongResponse, error) {
	params := sreq.Params{
		"songId": songId,
	}

	resp := new(SongResponse)
	err := a.Request(sreq.MethodGet, APIGetSong,
		sreq.WithQuery(params),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != "000000" {
		return nil, fmt.Errorf("get song: %s", resp.Info)
	}

	return resp, nil
}

func (a *API) GetSongURL(contentId, resourceType string) (string, error) {
	resp, err := a.GetSongURLRaw(contentId, resourceType)
	if err != nil {
		return "", err
	}

	return resp.Data.URL, nil
}

// 获取歌曲播放地址
func (a *API) GetSongURLRaw(contentId, resourceType string) (*SongURLResponse, error) {
	params := sreq.Params{
		"contentId":             contentId,
		"lowerQualityContentId": contentId,
		"resourceType":          resourceType,
	}

	resp := new(SongURLResponse)
	err := a.Request(sreq.MethodGet, APIGetSongURL,
		sreq.WithQuery(params),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != "000000" {
		return nil, fmt.Errorf("get song url: %s", resp.Info)
	}

	return resp, nil
}

func (a *API) GetSongPic(songId string) (string, error) {
	resp, err := a.GetSongPicRaw(songId)
	if err != nil {
		return "", err
	}
	return resp.LargePic, nil
}

// 获取歌曲专辑封面
func (a *API) GetSongPicRaw(songId string) (*SongPicResponse, error) {
	params := sreq.Params{
		"songId": songId,
	}

	resp := new(SongPicResponse)
	err := a.Request(sreq.MethodGet, APIGetSongPic,
		sreq.WithQuery(params),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.ReturnCode != "000000" {
		return nil, fmt.Errorf("get song pic: %s", resp.Msg)
	}

	return resp, nil
}

func (a *API) GetSongLyric(copyrightId string) (string, error) {
	resp, err := a.GetSongLyricRaw(copyrightId)
	if err != nil {
		return "", err
	}
	return resp.Lyric, nil
}

// 获取歌词
func (a *API) GetSongLyricRaw(copyrightId string) (*SongLyricResponse, error) {
	params := sreq.Params{
		"copyrightId": copyrightId,
	}

	resp := new(SongLyricResponse)
	err := a.Request(sreq.MethodGet, APIGetSongLyric,
		sreq.WithQuery(params),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.ReturnCode != "000000" {
		return nil, fmt.Errorf("get song lyric: %s", resp.Msg)
	}

	return resp, nil
}
