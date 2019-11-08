package kugou

import (
	"crypto/md5"
	"errors"
	"fmt"
	"math/rand"

	"github.com/winterssy/mxget/pkg/provider"
	"github.com/winterssy/sreq"
)

func GetSong(hash string) (*provider.Song, error) {
	return std.GetSong(hash)
}

func (a *API) GetSong(hash string) (*provider.Song, error) {
	resp, err := a.GetSongRaw(hash)
	if err != nil {
		return nil, err
	}

	a.patchAlbumInfo(&resp.Song)
	a.patchSongLyric(&resp.Song)
	songs := resolve(&resp.Song)
	return songs[0], nil
}

func GetSongRaw(hash string) (*SongResponse, error) {
	return std.GetSongRaw(hash)
}

// 获取歌曲信息
func (a *API) GetSongRaw(hash string) (*SongResponse, error) {
	params := sreq.Params{
		"hash": hash,
	}

	resp := new(SongResponse)
	err := a.Request(sreq.MethodGet, APIGetSong,
		sreq.WithQuery(params),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.ErrCode != 0 {
		return nil, fmt.Errorf("get song: %s", resp.Error)
	}

	return resp, nil
}

func GetSongURL(hash string) (string, error) {
	return std.GetSongURL(hash)
}

func (a *API) GetSongURL(hash string) (string, error) {
	resp, err := a.GetSongURLRaw(hash)
	if err != nil {
		return "", err
	}
	if len(resp.URL) == 0 {
		return "", errors.New("get song url: no data")
	}

	return resp.URL[rand.Intn(len(resp.URL))], nil
}

func GetSongURLRaw(hash string) (*SongURLResponse, error) {
	return std.GetSongURLRaw(hash)
}

// 获取歌曲播放地址
func (a *API) GetSongURLRaw(hash string) (*SongURLResponse, error) {
	data := []byte(hash + "kgcloudv2")
	key := fmt.Sprintf("%x", md5.Sum(data))

	params := sreq.Params{
		"hash": hash,
		"key":  key,
	}

	resp := new(SongURLResponse)
	err := a.Request(sreq.MethodGet, APIGetSongURL,
		sreq.WithQuery(params),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.Status != 1 {
		if resp.Status == 2 {
			err = errors.New("get song url: copyright protection")
		} else {
			err = fmt.Errorf("get song url: %d", resp.Status)
		}
		return nil, err
	}

	return resp, nil
}

func GetSongLyric(hash string) (string, error) {
	return std.GetSongLyric(hash)
}

// 获取歌词
func (a *API) GetSongLyric(hash string) (string, error) {
	params := sreq.Params{
		"hash": hash,
	}
	return a.Request(sreq.MethodGet, APIGetSongLyric,
		sreq.WithQuery(params),
	).Text()
}
