package qq

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/winterssy/mxget/pkg/provider"
	"github.com/winterssy/sreq"
)

func GetArtist(singerMid string) (*provider.Artist, error) {
	return std.GetArtist(singerMid)
}

func (a *API) GetArtist(singerMid string) (*provider.Artist, error) {
	resp, err := a.GetArtistRaw(singerMid, 0, 50)
	if err != nil {
		return nil, err
	}

	n := len(resp.Data.List)
	if n == 0 {
		return nil, errors.New("get artist: no data")
	}

	_songs := make([]*Song, 0, n)
	for _, i := range resp.Data.List {
		_songs = append(_songs, i.MusicData)
	}

	a.patchSongURLV1(_songs...)
	a.patchSongLyric(_songs...)
	songs := resolve(_songs...)
	return &provider.Artist{
		Name:   strings.TrimSpace(resp.Data.SingerName),
		PicURL: fmt.Sprintf(ArtistPicURL, resp.Data.SingerMid),
		Count:  n,
		Songs:  songs,
	}, nil
}

func GetArtistRaw(singerMid string, page int, pageSize int) (*ArtistResponse, error) {
	return std.GetArtistRaw(singerMid, page, pageSize)
}

// 获取歌手
func (a *API) GetArtistRaw(singerMid string, page int, pageSize int) (*ArtistResponse, error) {
	params := sreq.Params{
		"singermid": singerMid,
		"begin":     strconv.Itoa(page),
		"num":       strconv.Itoa(pageSize),
	}

	resp := new(ArtistResponse)
	err := a.Request(sreq.MethodGet, APIGetArtist,
		sreq.WithQuery(params),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf("get artist: %d", resp.Code)
	}

	return resp, nil
}
