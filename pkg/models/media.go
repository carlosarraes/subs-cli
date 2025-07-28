package models

import "time"

type MediaInfo struct {
	Title    string `json:"title"`
	Year     string `json:"year,omitempty"`
	Season   int    `json:"season,omitempty"`
	Episode  int    `json:"episode,omitempty"`
	Quality  string `json:"quality,omitempty"`
	Source   string `json:"source,omitempty"`
	Codec    string `json:"codec,omitempty"`
	Language string `json:"language,omitempty"`
	Type     string `json:"type"`
}

type SearchParams struct {
	Query     string `json:"query"`
	Language  string `json:"language"`
	Season    int    `json:"season,omitempty"`
	Episode   int    `json:"episode,omitempty"`
	Year      int    `json:"year,omitempty"`
	Type      string `json:"type"`
	MovieHash string `json:"movie_hash,omitempty"`
}

type Subtitle struct {
	ID          string    `json:"id"`
	Language    string    `json:"language"`
	ReleaseName string    `json:"release_name"`
	FileName    string    `json:"file_name"`
	FileID      string    `json:"file_id"`
	Uploader    string    `json:"uploader"`
	Rating      float64   `json:"rating"`
	Downloads   int       `json:"download_count"`
	UploadDate  time.Time `json:"upload_date"`
	MovieHash   string    `json:"movie_hash"`
	FPS         float64   `json:"fps"`
	Duration    int       `json:"duration"`
	SubFormat   string    `json:"sub_format"`
}

func (m *MediaInfo) IsEpisode() bool {
	return m.Type == "episode"
}

func (m *MediaInfo) IsMovie() bool {
	return m.Type == "movie"
}

func (m *MediaInfo) HasSeasonEpisode() bool {
	return m.Season > 0 && m.Episode > 0
}

func (m *MediaInfo) GetDisplayTitle() string {
	if m.Year != "" {
		return m.Title + " (" + m.Year + ")"
	}
	return m.Title
}
