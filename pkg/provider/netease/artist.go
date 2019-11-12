package netease

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/winterssy/mxget/pkg/provider"
	"github.com/winterssy/sreq"
)

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
	songs := resolve(resp.HotSongs...)
	return &provider.Artist{
		Id:     strconv.Itoa(resp.Artist.Id),
		Name:   strings.TrimSpace(resp.Artist.Name),
		PicURL: resp.Artist.PicURL,
		Count:  n,
		Songs:  songs,
	}, nil
}

// 获取歌手
func (a *API) GetArtistRaw(id int) (*ArtistResponse, error) {
	resp := new(ArtistResponse)
	err := a.Request(sreq.MethodPost, fmt.Sprintf(APIGetArtist, id),
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
