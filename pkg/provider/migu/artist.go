package migu

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
	if len(artistInfo.Resource) == 0 {
		return nil, errors.New("get artist info: no data")
	}

	artistSong, err := a.GetArtistSongRaw(singerId, 1, 50)
	if err != nil {
		return nil, err
	}
	if len(artistSong.Data.ContentItemList) == 0 ||
		len(artistSong.Data.ContentItemList[0].ItemList) == 0 {
		return nil, errors.New("get artist song: no data")
	}

	itemList := artistSong.Data.ContentItemList[0].ItemList
	n := len(itemList)
	_songs := make([]*Song, 0, n/2)
	for i := 0; i < n; i += 2 {
		_songs = append(_songs, &itemList[i].Song)
	}

	a.patchSongURL(SongDefaultBR, _songs...)
	a.patchSongLyric(_songs...)
	songs := a.resolve(_songs...)
	return &provider.Artist{
		Name:   strings.TrimSpace(artistInfo.Resource[0].Singer),
		PicURL: a.picURL(artistInfo.Resource[0].Imgs),
		Count:  len(songs),
		Songs:  songs,
	}, nil
}

func GetArtistInfoRaw(singerId string) (*ArtistInfoResponse, error) {
	return std.GetArtistInfoRaw(singerId)
}

// 获取歌手信息
func (a *API) GetArtistInfoRaw(singerId string) (*ArtistInfoResponse, error) {
	params := sreq.Params{
		"resourceId": singerId,
	}

	resp := new(ArtistInfoResponse)
	err := a.Request(sreq.MethodGet, GetArtistInfoAPI,
		sreq.WithQuery(params),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != "000000" {
		return nil, fmt.Errorf("get artist info: %s", resp.Info)
	}

	return resp, nil
}

func GetArtistSongRaw(singerId string, page int, pageSize int) (*ArtistSongResponse, error) {
	return std.GetArtistSongRaw(singerId, page, pageSize)
}

// 获取歌手歌曲，page: 页码；pageSize: 每页数量
func (a *API) GetArtistSongRaw(singerId string, page int, pageSize int) (*ArtistSongResponse, error) {
	params := sreq.Params{
		"singerId": singerId,
		"pageNo":   strconv.Itoa(page),
		"pageSize": strconv.Itoa(pageSize),
	}

	resp := new(ArtistSongResponse)
	err := a.Request(sreq.MethodGet, GetArtistSongAPI,
		sreq.WithQuery(params),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != "000000" {
		return nil, fmt.Errorf("get artist song: %s", resp.Info)
	}

	return resp, nil
}
