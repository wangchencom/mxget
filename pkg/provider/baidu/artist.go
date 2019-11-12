package baidu

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/winterssy/mxget/pkg/provider"
	"github.com/winterssy/sreq"
)

func (a *API) GetArtist(tingUid string) (*provider.Artist, error) {
	resp, err := a.GetArtistRaw(tingUid, 0, 50)
	if err != nil {
		return nil, err
	}

	n := len(resp.SongList)
	if n == 0 {
		return nil, errors.New("get artist: no data")
	}

	a.patchSongURL(resp.SongList...)
	a.patchSongLyric(resp.SongList...)
	songs := resolve(resp.SongList...)
	return &provider.Artist{
		Id:     resp.ArtistInfo.TingUid,
		Name:   strings.TrimSpace(resp.ArtistInfo.Name),
		PicURL: strings.SplitN(resp.ArtistInfo.AvatarBig, "@", 2)[0],
		Count:  n,
		Songs:  songs,
	}, nil
}

// 获取歌手
func (a *API) GetArtistRaw(tingUid string, offset int, limits int) (*ArtistResponse, error) {
	params := sreq.Params{
		"tinguid": tingUid,
		"offset":  strconv.Itoa(offset),
		"limits":  strconv.Itoa(limits),
	}

	resp := new(ArtistResponse)
	err := a.Request(sreq.MethodGet, APIGetArtist,
		sreq.WithQuery(params),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.ErrorCode != 22000 {
		return nil, fmt.Errorf("get artist: %v", resp.errorMessage())
	}

	return resp, nil
}
