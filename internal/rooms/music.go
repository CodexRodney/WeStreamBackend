package rooms

import (
	"mime/multipart"
)

type Music struct {
	title          string
	durationInSecs int64
	artist         string
	musicFile      multipart.File
}

func SetMusic(
	title string,
	durationInSecs int64,
	artist string,
	musicFile multipart.File) Music {
	m := Music{}
	m.title = title
	m.durationInSecs = durationInSecs
	m.artist = artist
	m.musicFile = musicFile
	return m
}

func (m *Music) CloseMusicFile() {
	m.musicFile.Close()
}

func (m *Music) GetMusicFile() multipart.File {
	return m.musicFile
}
