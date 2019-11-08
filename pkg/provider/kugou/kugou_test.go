package kugou

import "testing"

func TestAPI_SearchSongs(t *testing.T) {
	result, err := SearchSongs("五月天")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)
}

func TestAPI_GetSong(t *testing.T) {
	song, err := GetSong("1571941D82D63AD614E35EAD9DB6A6A2")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(song)
}

func TestAPI_GetSongURL(t *testing.T) {
	url, err := GetSongURL("1571941D82D63AD614E35EAD9DB6A6A2")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(url)
}

func TestAPI_GetSongLyric(t *testing.T) {
	lyric, err := GetSongLyric("1571941D82D63AD614E35EAD9DB6A6A2")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(lyric)
}

func TestAPI_GetArtist(t *testing.T) {
	artist, err := GetArtist("8965")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(artist)
}

func TestAPI_GetAlbum(t *testing.T) {
	album, err := GetAlbum("976965")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(album)
}

func TestAPI_GetPlaylist(t *testing.T) {
	playlist, err := GetPlaylist("610433")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(playlist)
}
