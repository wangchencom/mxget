package netease_test

import (
	"os"
	"testing"

	"github.com/winterssy/mxget/pkg/provider/netease"
)

var client *netease.API

func setup() {
	client = netease.New(nil)
}

func TestMain(m *testing.M) {
	setup()
	os.Exit(m.Run())
}

func TestAPI_SearchSongs(t *testing.T) {
	result, err := client.SearchSongs("Alan Walker")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)
}

func TestAPI_GetSong(t *testing.T) {
	song, err := client.GetSong("444269135")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(song)
}

func TestAPI_GetSongURL(t *testing.T) {
	url, err := client.GetSongURL(444269135, 320)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(url)
}

func TestAPI_GetSongLyric(t *testing.T) {
	lyric, err := client.GetSongLyric(444269135)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(lyric)
}

func TestAPI_GetArtist(t *testing.T) {
	artist, err := client.GetArtist("1045123")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(artist)
}

func TestAPI_GetAlbum(t *testing.T) {
	album, err := client.GetAlbum("35023284")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(album)
}

func TestAPI_GetPlaylist(t *testing.T) {
	playlist, err := client.GetPlaylist("156934569")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(playlist)
}
