package netease

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
	id, err := strconv.Atoi(artistId)
	if err != nil {
		return nil, err
	}

	resp, err := a.GetArtistRaw(id)
	if err != nil {
		return nil, err
	}

	n := len(resp.HotSongs)
	if n == 0 {
		return nil, errors.New("get artist: no data")
	}

	a.patchSongURL(SongDefaultBR, resp.HotSongs...)
	a.patchSongLyric(resp.HotSongs...)
	songs := a.resolve(resp.HotSongs...)
	return &provider.Artist{
		Name:   strings.TrimSpace(resp.Artist.Name),
		PicURL: resp.Artist.PicURL,
		Count:  n,
		Songs:  songs,
	}, nil
}

func GetArtistRaw(id int) (*ArtistResponse, error) {
	return std.GetArtistRaw(id)
}

// 获取歌手
func (a *API) GetArtistRaw(id int) (*ArtistResponse, error) {
	resp := new(ArtistResponse)
	err := a.Request(sreq.MethodPost, fmt.Sprintf(GetArtistAPI, id),
		sreq.WithForm(weapi(struct{}{})),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != 200 {
		return nil, fmt.Errorf("get artist: %s", resp.Msg)
	}

	return resp, nil
}
