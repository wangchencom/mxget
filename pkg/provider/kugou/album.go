package kugou

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/winterssy/mxget/pkg/provider"
	"github.com/winterssy/sreq"
)

func (a *API) GetAlbum(albumId string) (*provider.Album, error) {
	albumInfo, err := a.GetAlbumInfoRaw(albumId)
	if err != nil {
		return nil, err
	}

	albumSongs, err := a.GetAlbumSongsRaw(albumId, 1, -1)
	if err != nil {
		return nil, err
	}

	n := len(albumSongs.Data.Info)
	if n == 0 {
		return nil, errors.New("get album songs: no data")
	}

	a.patchSongInfo(albumSongs.Data.Info...)
	a.patchAlbumInfo(albumSongs.Data.Info...)
	a.patchSongLyric(albumSongs.Data.Info...)
	songs := resolve(albumSongs.Data.Info...)
	return &provider.Album{
		Name:   strings.TrimSpace(albumInfo.Data.AlbumName),
		PicURL: strings.ReplaceAll(albumInfo.Data.ImgURL, "{size}", "480"),
		Count:  n,
		Songs:  songs,
	}, nil
}

// 获取专辑信息
func (a *API) GetAlbumInfoRaw(albumId string) (*AlbumInfoResponse, error) {
	params := sreq.Params{
		"albumid": albumId,
	}

	resp := new(AlbumInfoResponse)
	err := a.Request(sreq.MethodGet, APIGetAlbumInfo,
		sreq.WithQuery(params),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.ErrCode != 0 {
		return nil, fmt.Errorf("get album info: %s", resp.Error)
	}

	return resp, nil
}

// 获取专辑歌曲，page: 页码；pageSize: 每页数量，-1获取全部
func (a *API) GetAlbumSongsRaw(albumId string, page int, pageSize int) (*AlbumSongsResponse, error) {
	params := sreq.Params{
		"albumid":  albumId,
		"page":     strconv.Itoa(page),
		"pagesize": strconv.Itoa(pageSize),
	}

	resp := new(AlbumSongsResponse)
	err := a.Request(sreq.MethodGet, APIGetAlbumSongs,
		sreq.WithQuery(params),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.ErrCode != 0 {
		return nil, fmt.Errorf("get album songs: %s", resp.Error)
	}

	return resp, nil
}
