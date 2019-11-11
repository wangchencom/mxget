package xiami_test

import (
	"os"
	"testing"

	"github.com/winterssy/mxget/pkg/provider/xiami"
)

var client *xiami.API

func setup() {
	client = xiami.New(nil)
}

func TestMain(m *testing.M) {
	setup()
	os.Exit(m.Run())
}

func TestAPI_SearchSongs(t *testing.T) {
	result, err := client.SearchSongs("五月天")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)
}

func TestAPI_GetSong(t *testing.T) {
	song, err := client.GetSong("xMPr7Lbbb28")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(song)
}

func TestAPI_GetSongLyric(t *testing.T) {
	lyric, err := client.GetSongLyric("xMPr7Lbbb28")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(lyric)
}

func TestAPI_GetArtist(t *testing.T) {
	artist, err := client.GetArtist("3110")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(artist)
}

func TestAPI_GetAlbum(t *testing.T) {
	album, err := client.GetAlbum("nmTM4c70144")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(album)
}

func TestAPI_GetPlaylist(t *testing.T) {
	playlist, err := client.GetPlaylist("8007523")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(playlist)
}
