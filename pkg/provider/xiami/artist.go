package xiami

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

	artistSongs, err := a.GetArtistSongsRaw(artistId, 1, 50)
	if err != nil {
		return nil, err
	}

	_songs := artistSongs.Data.Data.Songs
	n := len(_songs)
	if n == 0 {
		return nil, errors.New("get artist songs: no data")
	}

	a.patchSongLyric(_songs...)
	songs := resolve(_songs...)
	return &provider.Artist{
		Name:   strings.TrimSpace(artistInfo.Data.Data.ArtistDetailVO.ArtistName),
		PicURL: artistInfo.Data.Data.ArtistDetailVO.ArtistLogo,
		Count:  n,
		Songs:  songs,
	}, nil
}

func GetArtistInfoRaw(artistId string) (*ArtistInfoResponse, error) {
	return std.GetArtistInfoRaw(artistId)
}

// 获取歌手信息
func (a *API) GetArtistInfoRaw(artistId string) (*ArtistInfoResponse, error) {
	token, err := a.getToken(APIGetArtistInfo)
	if err != nil {
		return nil, err
	}

	model := make(map[string]string)
	_, err = strconv.Atoi(artistId)
	if err != nil {
		model["artistStringId"] = artistId
	} else {
		model["artistId"] = artistId
	}
	params := sreq.Params(signPayload(token, model))
	resp := new(ArtistInfoResponse)
	err = a.Request(sreq.MethodGet, APIGetArtistInfo, sreq.WithQuery(params)).JSON(resp)
	if err != nil {
		return nil, err
	}
	err = resp.check()
	if err != nil {
		return nil, fmt.Errorf("get artist info: %w", err)
	}

	return resp, nil
}

func GetArtistSongsRaw(artistId string, page int, pageSize int) (*ArtistSongsResponse, error) {
	return std.GetArtistSongsRaw(artistId, page, pageSize)
}

// 获取歌手歌曲
func (a *API) GetArtistSongsRaw(artistId string, page int, pageSize int) (*ArtistSongsResponse, error) {
	token, err := a.getToken(APIGetArtistSongs)
	if err != nil {
		return nil, err
	}

	model := map[string]interface{}{
		"pagingVO": map[string]int{
			"page":     page,
			"pageSize": pageSize,
		},
	}
	_, err = strconv.Atoi(artistId)
	if err != nil {
		model["artistStringId"] = artistId
	} else {
		model["artistId"] = artistId
	}
	params := sreq.Params(signPayload(token, model))
	resp := new(ArtistSongsResponse)
	err = a.Request(sreq.MethodGet, APIGetArtistSongs, sreq.WithQuery(params)).JSON(resp)
	if err != nil {
		return nil, err
	}
	err = resp.check()
	if err != nil {
		return nil, fmt.Errorf("get artist songs: %w", err)
	}

	return resp, nil
}
