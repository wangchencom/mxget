package kugou

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/winterssy/mxget/pkg/provider"
	"github.com/winterssy/sreq"
)

func GetArtist(singerId string) (*provider.Artist, error) {
	return std.GetArtist(singerId)
}

func (a *API) GetArtist(singerId string) (*provider.Artist, error) {
	artistInfo, err := a.GetArtistInfoRaw(singerId)
	if err != nil {
		return nil, err
	}

	artistSongs, err := a.GetArtistSongsRaw(singerId, 1, 50)
	if err != nil {
		return nil, err
	}

	n := len(artistSongs.Data.Info)
	if n == 0 {
		return nil, errors.New("get artist songs: no data")
	}

	a.patchSongInfo(artistSongs.Data.Info...)
	a.patchAlbumInfo(artistSongs.Data.Info...)
	a.patchSongLyric(artistSongs.Data.Info...)
	songs := resolve(artistSongs.Data.Info...)
	return &provider.Artist{
		Name:   strings.TrimSpace(artistInfo.Data.SingerName),
		PicURL: strings.ReplaceAll(artistInfo.Data.ImgURL, "{size}", "480"),
		Count:  n,
		Songs:  songs,
	}, nil
}

func GetArtistInfoRaw(singerId string) (*ArtistInfoResponse, error) {
	return std.GetArtistInfoRaw(singerId)
}

// 获取歌手信息
func (a *API) GetArtistInfoRaw(singerId string) (*ArtistInfoResponse, error) {
	params := sreq.Params{
		"singerid": singerId,
	}

	resp := new(ArtistInfoResponse)
	err := a.Request(sreq.MethodGet, APIGetArtistInfo,
		sreq.WithQuery(params),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.ErrCode != 0 {
		return nil, fmt.Errorf("get artist info: %s", resp.Error)
	}

	return resp, nil
}

func GetArtistSongsRaw(singerId string, page int, pageSize int) (*ArtistSongsResponse, error) {
	return std.GetArtistSongsRaw(singerId, page, pageSize)
}

// 获取歌手歌曲，page: 页码；pageSize: 每页数量，-1获取全部
func (a *API) GetArtistSongsRaw(singerId string, page int, pageSize int) (*ArtistSongsResponse, error) {
	params := sreq.Params{
		"singerid": singerId,
		"page":     strconv.Itoa(page),
		"pagesize": strconv.Itoa(pageSize),
	}

	resp := new(ArtistSongsResponse)
	err := a.Request(sreq.MethodGet, APIGetArtistSongs,
		sreq.WithQuery(params),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.ErrCode != 0 {
		return nil, fmt.Errorf("get artist songs: %s", resp.Error)
	}

	return resp, nil
}
