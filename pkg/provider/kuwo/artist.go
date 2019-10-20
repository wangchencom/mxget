package kuwo

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/winterssy/mxget/pkg/provider"
	"github.com/winterssy/sreq"
)

func GetArtist(artistId string) (*provider.Artist, error) {
	return std.GetArtist(artistId)
}

func (a *API) GetArtist(artistId string) (*provider.Artist, error) {
	artistInfo, err := a.GetArtistInfoRaw(artistId)
	if err != nil {
		return nil, err
	}

	resp, err := a.GetArtistSongRaw(artistId, 1, 50)
	if err != nil {
		return nil, err
	}

	n := len(resp.Data.List)
	if n == 0 {
		return nil, errors.New("get artist song: no data")
	}

	a.patchSongURL(SongDefaultBR, resp.Data.List...)
	a.patchSongLyric(resp.Data.List...)
	songs := a.resolve(resp.Data.List)
	return &provider.Artist{
		Name:   strings.TrimSpace(artistInfo.Data.Name),
		PicURL: artistInfo.Data.Pic300,
		Count:  n,
		Songs:  songs,
	}, nil
}

func GetArtistInfoRaw(artistId string) (*ArtistInfoResponse, error) {
	return std.GetArtistInfoRaw(artistId)
}

// 获取歌手信息
func (a *API) GetArtistInfoRaw(artistId string) (*ArtistInfoResponse, error) {
	params := sreq.Params{
		"artistid": artistId,
	}

	resp := new(ArtistInfoResponse)
	err := a.Request(sreq.MethodGet, GetArtistInfoAPI,
		sreq.WithQuery(params),
	).JSON(resp)
	if err != nil {
		return nil, err
	}

	if resp.Code != 200 {
		return nil, fmt.Errorf("get artist info: %s", resp.Msg)
	}

	return resp, nil
}

func GetArtistSongRaw(artistId string, page int, pageSize int) (*ArtistSongResponse, error) {
	return std.GetArtistSongRaw(artistId, page, pageSize)
}

// 获取歌手歌曲，page: 页码； pageSize: 每页数量，如果要获取全部请设置为较大的值
func (a *API) GetArtistSongRaw(artistId string, page int, pageSize int) (*ArtistSongResponse, error) {
	params := sreq.Params{
		"artistid": artistId,
		"pn":       strconv.Itoa(page),
		"rn":       strconv.Itoa(pageSize),
	}

	resp := new(ArtistSongResponse)
	err := a.Request(sreq.MethodGet, GetArtistSongAPI,
		sreq.WithQuery(params),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != 200 {
		return nil, fmt.Errorf("get artist song: %s", resp.Msg)
	}

	return resp, nil
}
