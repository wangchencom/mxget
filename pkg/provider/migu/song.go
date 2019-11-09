package migu

import (
	"errors"
	"fmt"

	"github.com/winterssy/mxget/pkg/provider"
	"github.com/winterssy/sreq"
)

func GetSongId(copyrightId string) (string, error) {
	return std.GetSongId(copyrightId)
}

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

func GetSongIdRaw(copyrightId string) (*SongIdResponse, error) {
	return std.GetSongIdRaw(copyrightId)
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

func GetSong(copyrightId string) (*provider.Song, error) {
	return std.GetSong(copyrightId)
}

func (a *API) GetSong(copyrightId string) (*provider.Song, error) {
	songId, err := a.GetSongId(copyrightId)
	if err != nil {
		return nil, err
	}

	resp, err := a.GetSongRaw(songId)
	if err != nil {
		return nil, err
	}
	if len(resp.Resource) == 0 {
		return nil, errors.New("get song: no data")
	}

	_song := resp.Resource[0]
	a.patchSongURL(SongDefaultBR, _song)
	lyric, err := a.GetSongLyric(_song.CopyrightId)
	if err == nil {
		_song.Lyric = lyric
	}
	songs := resolve(_song)
	return songs[0], nil
}

func GetSongRaw(songId string) (*SongResponse, error) {
	return std.GetSongRaw(songId)
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

func GetSongURL(contentId string, br int) string {
	return std.GetSongURL(contentId, br)
}

func (a *API) GetSongURL(contentId string, br int) string {
	var _br int
	switch br {
	case 64, 128, 320, 999:
		_br = br
	default:
		_br = 320
	}
	return fmt.Sprintf(SongURL, contentId, "E", codeRate[_br])
}

func GetSongURLRaw(contentId, resourceType string) (*SongURLResponse, error) {
	return std.GetSongURLRaw(contentId, resourceType)
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

func GetSongPic(songId string) (string, error) {
	return std.GetSongPic(songId)
}

func (a *API) GetSongPic(songId string) (string, error) {
	resp, err := a.GetSongPicRaw(songId)
	if err != nil {
		return "", err
	}
	return resp.LargePic, nil
}

func GetSongPicRaw(songId string) (*SongPicResponse, error) {
	return std.GetSongPicRaw(songId)
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

func GetSongLyric(copyrightId string) (string, error) {
	return std.GetSongLyric(copyrightId)
}

func (a *API) GetSongLyric(copyrightId string) (string, error) {
	resp, err := a.GetSongLyricRaw(copyrightId)
	if err != nil {
		return "", err
	}
	return resp.Lyric, nil
}

func GetSongLyricRaw(copyrightId string) (*SongLyricResponse, error) {
	return std.GetSongLyricRaw(copyrightId)
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
