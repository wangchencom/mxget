package migu

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/winterssy/mxget/pkg/provider"
	"github.com/winterssy/sreq"
)

func (a *API) GetArtist(singerId string) (*provider.Artist, error) {
	artistInfo, err := a.GetArtistInfoRaw(singerId)
	if err != nil {
		return nil, err
	}
	if len(artistInfo.Resource) == 0 {
		return nil, errors.New("get artist info: no data")
	}

	artistSongs, err := a.GetArtistSongsRaw(singerId, 1, 50)
	if err != nil {
		return nil, err
	}
	if len(artistSongs.Data.ContentItemList) == 0 ||
		len(artistSongs.Data.ContentItemList[0].ItemList) == 0 {
		return nil, errors.New("get artist songs: no data")
	}

	itemList := artistSongs.Data.ContentItemList[0].ItemList
	n := len(itemList)
	_songs := make([]*Song, 0, n/2)
	for i := 0; i < n; i += 2 {
		_songs = append(_songs, &itemList[i].Song)
	}

	a.patchSongLyric(_songs...)
	songs := resolve(_songs...)
	return &provider.Artist{
		Id:     artistInfo.Resource[0].SingerId,
		Name:   strings.TrimSpace(artistInfo.Resource[0].Singer),
		PicURL: picURL(artistInfo.Resource[0].Imgs),
		Count:  len(songs),
		Songs:  songs,
	}, nil
}

// 获取歌手信息
func (a *API) GetArtistInfoRaw(singerId string) (*ArtistInfoResponse, error) {
	params := sreq.Params{
		"resourceId": singerId,
	}

	resp := new(ArtistInfoResponse)
	err := a.Request(sreq.MethodGet, APIGetArtistInfo,
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

// 获取歌手歌曲，page: 页码；pageSize: 每页数量
func (a *API) GetArtistSongsRaw(singerId string, page int, pageSize int) (*ArtistSongsResponse, error) {
	params := sreq.Params{
		"singerId": singerId,
		"pageNo":   strconv.Itoa(page),
		"pageSize": strconv.Itoa(pageSize),
	}

	resp := new(ArtistSongsResponse)
	err := a.Request(sreq.MethodGet, APIGetArtistSongs,
		sreq.WithQuery(params),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != "000000" {
		return nil, fmt.Errorf("get artist songs: %s", resp.Info)
	}

	return resp, nil
}
