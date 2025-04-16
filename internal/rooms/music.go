package rooms

import (
	"mime/multipart"
)

type Music struct {
	title          string
	durationInSecs int64
	artist         string
	addedBy        string
	musicFile      multipart.File
}

func SetMusic(
	title string,
	durationInSecs int64,
	artist string,
	addedBy string,
	musicFile multipart.File) *Music {
	m := Music{}
	m.title = title
	m.durationInSecs = durationInSecs
	m.artist = artist
	m.musicFile = musicFile
	m.addedBy = addedBy
	return &m
}

func (m *Music) CloseMusicFile() {
	m.musicFile.Close()
}

func (m *Music) GetMusicFile() multipart.File {
	return m.musicFile
}

func (m *Music) extractMusicInfo() map[string]interface{} {
	return map[string]interface{}{
		"title":          m.title,
		"durationInSecs": int(m.durationInSecs),
		"artist":         m.artist,
		"addedBy":        m.addedBy,
	}
}

func ExtractMusicInfo(musics []*Music) []map[string]interface{} {
	musicsInfo := make([]map[string]interface{}, 0, len(musics))

	for _, m := range musics {
		info := map[string]interface{}{
			"title":          m.title,
			"durationInSecs": int(m.durationInSecs),
			"artist":         m.artist,
			"addedBy":        m.addedBy,
		}
		musicsInfo = append(musicsInfo, info)
	}

	return musicsInfo
}
