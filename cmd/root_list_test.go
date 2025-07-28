package cmd

import (
	"testing"
	"time"

	"github.com/carlosarraes/subs-cli/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestDisplaySubtitleList(t *testing.T) {
	t.Parallel()

	t.Run("displays subtitle list correctly", func(t *testing.T) {
		t.Parallel()

		cli := &CLI{DryRun: true}
		uploadDate := time.Date(2023, 1, 15, 10, 30, 0, 0, time.UTC)

		subtitles := []*models.Subtitle{
			{
				ID:          "1",
				Language:    "en",
				ReleaseName: "The.Office.S03E07.720p.BluRay.x264-GROUP",
				FileName:    "The.Office.S03E07.srt",
				Uploader:    "TestUploader",
				Rating:      8.5,
				Downloads:   1500,
				UploadDate:  uploadDate,
			},
			{
				ID:          "2",
				Language:    "pt-BR",
				ReleaseName: "The.Office.S03E07.720p.BluRay.x264-ANOTHER",
				FileName:    "The.Office.S03E07.pt-BR.srt",
				Uploader:    "AnotherUploader",
				Rating:      7.2,
				Downloads:   850,
				UploadDate:  uploadDate,
			},
		}

		assert.NotPanics(t, func() {
			cli.displaySubtitleList(subtitles)
		})
	})

	t.Run("handles long release names", func(t *testing.T) {
		t.Parallel()

		cli := &CLI{DryRun: false}
		subtitles := []*models.Subtitle{
			{
				ID:          "1",
				Language:    "en",
				ReleaseName: "This.Is.A.Very.Long.Release.Name.That.Should.Be.Truncated.Because.It.Exceeds.The.Display.Limit",
				Uploader:    "VeryLongUploaderNameThatShouldAlsoBeTruncated",
				Rating:      9.1,
				Downloads:   25000,
			},
		}

		assert.NotPanics(t, func() {
			cli.displaySubtitleList(subtitles)
		})
	})
}

func TestCreateSearchParams(t *testing.T) {
	t.Parallel()

	cli := &CLI{}

	t.Run("creates movie search params", func(t *testing.T) {
		t.Parallel()

		mediaInfo := &models.MediaInfo{
			Title: "Inception",
			Year:  "2010",
			Type:  "movie",
		}

		params := cli.createSearchParams(mediaInfo)

		assert.Equal(t, "Inception", params.Query)
		assert.Equal(t, "movie", params.Type)
		assert.Equal(t, 2010, params.Year)
		assert.Equal(t, 0, params.Season)
		assert.Equal(t, 0, params.Episode)
	})

	t.Run("creates TV episode search params", func(t *testing.T) {
		t.Parallel()

		mediaInfo := &models.MediaInfo{
			Title:   "The Office",
			Year:    "2005",
			Season:  3,
			Episode: 7,
			Type:    "episode",
		}

		params := cli.createSearchParams(mediaInfo)

		assert.Equal(t, "The Office", params.Query)
		assert.Equal(t, "episode", params.Type)
		assert.Equal(t, 2005, params.Year)
		assert.Equal(t, 3, params.Season)
		assert.Equal(t, 7, params.Episode)
	})

	t.Run("handles missing year", func(t *testing.T) {
		t.Parallel()

		mediaInfo := &models.MediaInfo{
			Title: "Some Show",
			Type:  "movie",
		}

		params := cli.createSearchParams(mediaInfo)

		assert.Equal(t, "Some Show", params.Query)
		assert.Equal(t, "movie", params.Type)
		assert.Equal(t, 0, params.Year)
	})

	t.Run("handles invalid year", func(t *testing.T) {
		t.Parallel()

		mediaInfo := &models.MediaInfo{
			Title: "Some Show",
			Year:  "invalid",
			Type:  "movie",
		}

		params := cli.createSearchParams(mediaInfo)

		assert.Equal(t, "Some Show", params.Query)
		assert.Equal(t, "movie", params.Type)
		assert.Equal(t, 0, params.Year)
	})
}

func TestTruncateString(t *testing.T) {
	t.Parallel()

	cli := &CLI{}

	tests := []struct {
		name     string
		input    string
		maxLen   int
		expected string
	}{
		{
			name:     "short string unchanged",
			input:    "short",
			maxLen:   10,
			expected: "short",
		},
		{
			name:     "exact length unchanged",
			input:    "exactly10c",
			maxLen:   10,
			expected: "exactly10c",
		},
		{
			name:     "long string truncated",
			input:    "this is a very long string that needs truncation",
			maxLen:   10,
			expected: "this is...",
		},
		{
			name:     "minimum truncation",
			input:    "abcdef",
			maxLen:   5,
			expected: "ab...",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := cli.truncateString(tt.input, tt.maxLen)
			assert.Equal(t, tt.expected, result)
		})
	}
}
