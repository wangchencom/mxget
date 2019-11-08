package netease

import "testing"

func TestAPI_SearchSongs(t *testing.T) {
	result, err := SearchSongs("Alan Walker")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)
}

func TestAPI_GetSong(t *testing.T) {
	song, err := GetSong("444269135")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(song)
}

func TestAPI_GetSongURL(t *testing.T) {
	url, err := GetSongURL(444269135, 320)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(url)
}

func TestAPI_GetSongLyric(t *testing.T) {
	lyric, err := GetSongLyric(444269135)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(lyric)
}

func TestAPI_GetArtist(t *testing.T) {
	artist, err := GetArtist("1045123")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(artist)
}

func TestAPI_GetAlbum(t *testing.T) {
	album, err := GetAlbum("35023284")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(album)
}

func TestAPI_GetPlaylist(t *testing.T) {
	playlist, err := GetPlaylist("156934569")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(playlist)
}
