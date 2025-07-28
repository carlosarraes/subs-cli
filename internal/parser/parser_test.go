package parser

import (
	"testing"

	"github.com/carlosarraes/subs-cli/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParser_Parse(t *testing.T) {
	t.Parallel()

	parser := New()

	tests := []struct {
		name     string
		filename string
		want     *models.MediaInfo
		wantErr  bool
	}{
		{
			name:     "TV with year SxxExx format",
			filename: "Dark.Matter.2024.S01E01.1080p.x265-ELiTE.mkv",
			want: &models.MediaInfo{
				Title:   "Dark Matter",
				Year:    "2024",
				Season:  1,
				Episode: 1,
				Quality: "1080p",
				Source:  "ELiTE",
				Codec:   "x265",
				Type:    "episode",
			},
		},
		{
			name:     "TV with year complex title",
			filename: "The.Walking.Dead.2010.S11E24.720p.BluRay.x264-GROUP.mkv",
			want: &models.MediaInfo{
				Title:   "The Walking Dead",
				Year:    "2010",
				Season:  11,
				Episode: 24,
				Quality: "720p",
				Source:  "BluRay.GROUP",
				Codec:   "x264",
				Type:    "episode",
			},
		},

		{
			name:     "TV without year SxxExx format",
			filename: "The.Office.S03E07.720p.BluRay.x264.mkv",
			want: &models.MediaInfo{
				Title:   "The Office",
				Season:  3,
				Episode: 7,
				Quality: "720p",
				Source:  "BluRay",
				Codec:   "x264",
				Type:    "episode",
			},
		},
		{
			name:     "TV without quality info",
			filename: "Friends.S10E18.mkv",
			want: &models.MediaInfo{
				Title:   "Friends",
				Season:  10,
				Episode: 18,
				Type:    "episode",
			},
		},

		{
			name:     "TV alternative xXx format",
			filename: "Series.Name.1x01.720p.WEB-DL.mkv",
			want: &models.MediaInfo{
				Title:   "Series Name",
				Season:  1,
				Episode: 1,
				Quality: "720p",
				Source:  "WEB-DL",
				Type:    "episode",
			},
		},
		{
			name:     "TV alternative with year xXx format",
			filename: "Breaking.Bad.2008.5x16.1080p.HDTV.x264-ASAP.mkv",
			want: &models.MediaInfo{
				Title:   "Breaking Bad",
				Year:    "2008",
				Season:  5,
				Episode: 16,
				Quality: "1080p",
				Source:  "HDTV.ASAP",
				Codec:   "x264",
				Type:    "episode",
			},
		},

		{
			name:     "TV 3-digit episode format",
			filename: "Series.Name.101.720p.x264.mkv",
			want: &models.MediaInfo{
				Title:   "Series Name",
				Season:  1,
				Episode: 1,
				Quality: "720p",
				Codec:   "x264",
				Type:    "episode",
			},
		},
		{
			name:     "TV 3-digit complex episode",
			filename: "Game.of.Thrones.315.1080p.HDTV.x265-DIMENSION.mkv",
			want: &models.MediaInfo{
				Title:   "Game of Thrones",
				Season:  3,
				Episode: 15,
				Quality: "1080p",
				Source:  "HDTV.DIMENSION",
				Codec:   "x265",
				Type:    "episode",
			},
		},

		{
			name:     "Movie with quality",
			filename: "Inception.2010.1080p.BluRay.x264-SPARKS.mkv",
			want: &models.MediaInfo{
				Title:   "Inception",
				Year:    "2010",
				Quality: "1080p",
				Source:  "BluRay.SPARKS",
				Codec:   "x264",
				Type:    "movie",
			},
		},
		{
			name:     "Movie complex title",
			filename: "The.Dark.Knight.Rises.2012.720p.WEB-DL.x264-YTS.mp4",
			want: &models.MediaInfo{
				Title:   "The Dark Knight Rises",
				Year:    "2012",
				Quality: "720p",
				Source:  "WEB-DL.YTS",
				Codec:   "x264",
				Type:    "movie",
			},
		},
		{
			name:     "Movie without quality",
			filename: "Pulp.Fiction.1994.BluRay.x264-GROUP.mp4",
			want: &models.MediaInfo{
				Title:  "Pulp Fiction",
				Year:   "1994",
				Source: "BluRay.GROUP",
				Codec:  "x264",
				Type:   "movie",
			},
		},

		{
			name:     "Filename with spaces TV",
			filename: "Dark Matter 2024 S01E01 1080p x265-ELiTE.mkv",
			want: &models.MediaInfo{
				Title:   "Dark Matter",
				Year:    "2024",
				Season:  1,
				Episode: 1,
				Quality: "1080p",
				Source:  "ELiTE",
				Codec:   "x265",
				Type:    "episode",
			},
		},
		{
			name:     "Filename with spaces movie",
			filename: "The Matrix 1999 1080p BluRay x264-GROUP.mp4",
			want: &models.MediaInfo{
				Title:   "The Matrix",
				Year:    "1999",
				Quality: "1080p",
				Source:  "BluRay.GROUP",
				Codec:   "x264",
				Type:    "movie",
			},
		},

		{
			name:     "HEVC codec",
			filename: "Series.Name.S01E01.2160p.UHD.BluRay.HEVC-GROUP.mkv",
			want: &models.MediaInfo{
				Title:   "Series Name",
				Season:  1,
				Episode: 1,
				Quality: "2160p",
				Source:  "UHD.BluRay.GROUP",
				Codec:   "HEVC",
				Type:    "episode",
			},
		},
		{
			name:     "AV1 codec",
			filename: "Movie.Name.2023.1080p.WEB-DL.AV1-ENCODER.mkv",
			want: &models.MediaInfo{
				Title:   "Movie Name",
				Year:    "2023",
				Quality: "1080p",
				Source:  "WEB-DL.ENCODER",
				Codec:   "AV1",
				Type:    "movie",
			},
		},

		{
			name:     "Full path handling",
			filename: "/path/to/media/The.Office.S03E07.720p.BluRay.x264.mkv",
			want: &models.MediaInfo{
				Title:   "The Office",
				Season:  3,
				Episode: 7,
				Quality: "720p",
				Source:  "BluRay",
				Codec:   "x264",
				Type:    "episode",
			},
		},

		{
			name:     "Invalid filename format",
			filename: "invalid_filename_format.mkv",
			wantErr:  true,
		},
		{
			name:     "No extension",
			filename: "Movie.Name.2023.1080p.BluRay.x264",
			want: &models.MediaInfo{
				Title:   "Movie Name",
				Year:    "2023",
				Quality: "1080p",
				Source:  "BluRay",
				Codec:   "x264",
				Type:    "movie",
			},
		},
		{
			name:     "Empty filename",
			filename: "",
			wantErr:  true,
		},
		{
			name:     "Just extension",
			filename: ".mkv",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := parser.Parse(tt.filename)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "unable to parse filename")
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)

			assert.Equal(t, tt.want.Title, got.Title, "Title mismatch")
			assert.Equal(t, tt.want.Year, got.Year, "Year mismatch")
			assert.Equal(t, tt.want.Season, got.Season, "Season mismatch")
			assert.Equal(t, tt.want.Episode, got.Episode, "Episode mismatch")
			assert.Equal(t, tt.want.Quality, got.Quality, "Quality mismatch")
			assert.Equal(t, tt.want.Source, got.Source, "Source mismatch")
			assert.Equal(t, tt.want.Codec, got.Codec, "Codec mismatch")
			assert.Equal(t, tt.want.Type, got.Type, "Type mismatch")
		})
	}
}

func TestParser_ValidationErrors(t *testing.T) {
	t.Parallel()

	parser := New()

	tests := []struct {
		name     string
		filename string
		errorMsg string
	}{
		{
			name:     "Invalid season number - zero",
			filename: "Series.Name.S00E01.720p.x264.mkv",
			errorMsg: "unable to parse filename",
		},
		{
			name:     "Invalid episode number - zero",
			filename: "Series.Name.S01E00.720p.x264.mkv",
			errorMsg: "unable to parse filename",
		},
		{
			name:     "Season too high",
			filename: "Series.Name.S100E01.720p.x264.mkv",
			errorMsg: "unable to parse filename",
		},
		{
			name:     "Episode too high",
			filename: "Series.Name.S01E1000.720p.x264.mkv",
			errorMsg: "unable to parse filename",
		},
		{
			name:     "Invalid year - too old",
			filename: "Movie.Name.1800.1080p.BluRay.x264.mkv",
			errorMsg: "unable to parse filename",
		},
		{
			name:     "Invalid year - future",
			filename: "Movie.Name.2050.1080p.BluRay.x264.mkv",
			errorMsg: "unable to parse filename",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := parser.Parse(tt.filename)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.errorMsg)
		})
	}
}

func TestCleanFilename(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		filename string
		want     string
	}{
		{
			name:     "Convert spaces to dots",
			filename: "Movie Name 2023 1080p BluRay x264.mkv",
			want:     "Movie.Name.2023.1080p.BluRay.x264.mkv",
		},
		{
			name:     "Remove path",
			filename: "/path/to/Movie.Name.2023.1080p.mkv",
			want:     "Movie.Name.2023.1080p.mkv",
		},
		{
			name:     "Multiple consecutive dots",
			filename: "Movie..Name...2023.mkv",
			want:     "Movie.Name.2023.mkv",
		},
		{
			name:     "Mixed spaces and dots",
			filename: "Movie Name.2023 1080p.mkv",
			want:     "Movie.Name.2023.1080p.mkv",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := cleanFilename(tt.filename)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCleanTitle(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		title string
		want  string
	}{
		{
			name:  "Convert dots to spaces",
			title: "The.Dark.Knight",
			want:  "The Dark Knight",
		},
		{
			name:  "Remove extra spaces",
			title: "The  Walking   Dead",
			want:  "The Walking Dead",
		},
		{
			name:  "Mixed dots and spaces",
			title: "Game.of Thrones",
			want:  "Game of Thrones",
		},
		{
			name:  "Trim spaces",
			title: "  Breaking Bad  ",
			want:  "Breaking Bad",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := cleanTitle(tt.title)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExtractSourceAndCodec(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		combined string
		wantSrc  string
		wantCode string
	}{
		{
			name:     "BluRay with x264",
			combined: "BluRay.x264-GROUP",
			wantSrc:  "BluRay.GROUP",
			wantCode: "x264",
		},
		{
			name:     "WEB-DL with HEVC",
			combined: "WEB-DL.HEVC.ENCODER",
			wantSrc:  "WEB-DL.ENCODER",
			wantCode: "HEVC",
		},
		{
			name:     "Multiple codecs, take first",
			combined: "BluRay.x264.x265-GROUP",
			wantSrc:  "BluRay.GROUP",
			wantCode: "x264",
		},
		{
			name:     "No codec",
			combined: "BluRay.HDTV-GROUP",
			wantSrc:  "BluRay.HDTV-GROUP",
			wantCode: "",
		},
		{
			name:     "Only codec",
			combined: "x265",
			wantSrc:  "",
			wantCode: "x265",
		},
		{
			name:     "Empty string",
			combined: "",
			wantSrc:  "",
			wantCode: "",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotSrc, gotCode := extractSourceAndCodec(tt.combined)
			assert.Equal(t, tt.wantSrc, gotSrc, "Source mismatch")
			assert.Equal(t, tt.wantCode, gotCode, "Codec mismatch")
		})
	}
}

func TestIsCodec(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		s    string
		want bool
	}{
		{"x264", "x264", true},
		{"x265", "x265", true},
		{"HEVC", "hevc", true},
		{"H264", "h264", true},
		{"AV1", "av1", true},
		{"Not codec", "bluray", false},
		{"Not codec", "group", false},
		{"Empty", "", false},
		{"Partial match", "some-x264-release", true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := isCodec(tt.s)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMediaInfo_Methods(t *testing.T) {
	t.Parallel()

	t.Run("IsEpisode", func(t *testing.T) {
		t.Parallel()

		episode := &models.MediaInfo{Type: "episode"}
		movie := &models.MediaInfo{Type: "movie"}

		assert.True(t, episode.IsEpisode())
		assert.False(t, movie.IsEpisode())
	})

	t.Run("IsMovie", func(t *testing.T) {
		t.Parallel()

		episode := &models.MediaInfo{Type: "episode"}
		movie := &models.MediaInfo{Type: "movie"}

		assert.False(t, episode.IsMovie())
		assert.True(t, movie.IsMovie())
	})

	t.Run("HasSeasonEpisode", func(t *testing.T) {
		t.Parallel()

		valid := &models.MediaInfo{Season: 1, Episode: 1}
		noSeason := &models.MediaInfo{Season: 0, Episode: 1}
		noEpisode := &models.MediaInfo{Season: 1, Episode: 0}

		assert.True(t, valid.HasSeasonEpisode())
		assert.False(t, noSeason.HasSeasonEpisode())
		assert.False(t, noEpisode.HasSeasonEpisode())
	})

	t.Run("GetDisplayTitle", func(t *testing.T) {
		t.Parallel()

		withYear := &models.MediaInfo{Title: "Inception", Year: "2010"}
		withoutYear := &models.MediaInfo{Title: "Inception"}

		assert.Equal(t, "Inception (2010)", withYear.GetDisplayTitle())
		assert.Equal(t, "Inception", withoutYear.GetDisplayTitle())
	})
}
