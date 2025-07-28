package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/carlosarraes/subs-cli/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewOpenSubtitlesClient(t *testing.T) {
	t.Parallel()

	t.Run("with default values", func(t *testing.T) {
		t.Parallel()

		config := &Config{
			Username: "test",
			Password: "pass",
		}

		client := NewOpenSubtitlesClient(config)

		require.NotNil(t, client)
		assert.Equal(t, DefaultBaseURL, client.config.BaseURL)
		assert.Equal(t, DefaultUserAgent, client.config.UserAgent)
	})

	t.Run("with custom values", func(t *testing.T) {
		t.Parallel()

		config := &Config{
			BaseURL:   "https://custom.api.com",
			UserAgent: "custom-agent/1.0",
			APIKey:    "test-key",
			Username:  "test",
			Password:  "pass",
		}

		client := NewOpenSubtitlesClient(config)

		require.NotNil(t, client)
		assert.Equal(t, "https://custom.api.com", client.config.BaseURL)
		assert.Equal(t, "custom-agent/1.0", client.config.UserAgent)
		assert.Equal(t, "test-key", client.config.APIKey)
	})
}

func TestOpenSubtitlesClient_Authenticate(t *testing.T) {
	t.Parallel()

	t.Run("successful authentication", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/login", r.URL.Path)
			assert.Equal(t, "POST", r.Method)

			var loginReq LoginRequest
			err := json.NewDecoder(r.Body).Decode(&loginReq)
			require.NoError(t, err)

			assert.Equal(t, "testuser", loginReq.Username)
			assert.Equal(t, "testpass", loginReq.Password)

			response := LoginResponse{
				Token:  "test-token-123",
				Status: 200,
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		config := &Config{
			BaseURL:  server.URL,
			Username: "testuser",
			Password: "testpass",
		}

		client := NewOpenSubtitlesClient(config)
		err := client.Authenticate(context.Background())

		require.NoError(t, err)
		assert.Equal(t, "test-token-123", client.token)
	})

	t.Run("missing credentials", func(t *testing.T) {
		t.Parallel()

		config := &Config{}
		client := NewOpenSubtitlesClient(config)

		err := client.Authenticate(context.Background())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "username and password are required")
	})

	t.Run("authentication failed", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Invalid credentials"))
		}))
		defer server.Close()

		config := &Config{
			BaseURL:  server.URL,
			Username: "wrong",
			Password: "wrong",
		}

		client := NewOpenSubtitlesClient(config)
		err := client.Authenticate(context.Background())

		require.Error(t, err)
		assert.Contains(t, err.Error(), "authentication failed with status 401")
	})
}

func TestOpenSubtitlesClient_Search(t *testing.T) {
	t.Parallel()

	t.Run("successful search", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/login" {
				response := LoginResponse{
					Token:  "test-token",
					Status: 200,
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
				return
			}

			if r.URL.Path == "/subtitles" {
				assert.Equal(t, "GET", r.Method)
				assert.Equal(t, "The Office", r.URL.Query().Get("query"))
				assert.Equal(t, "en", r.URL.Query().Get("languages"))
				assert.Equal(t, "episode", r.URL.Query().Get("type"))
				assert.Equal(t, "2005", r.URL.Query().Get("year"))
				assert.Equal(t, "3", r.URL.Query().Get("season_number"))
				assert.Equal(t, "7", r.URL.Query().Get("episode_number"))

				mockResponse := map[string]interface{}{
					"total_count": 1,
					"data": []map[string]interface{}{
						{
							"id":   "test-id-123",
							"type": "subtitle",
							"attributes": map[string]interface{}{
								"language":       "en",
								"download_count": 1500,
								"fps":            23.976,
								"ratings":        8.5,
								"upload_date":    "2023-01-15T10:30:00",
								"release":        "The.Office.S03E07.720p.BluRay.x264",
								"uploader": map[string]interface{}{
									"name": "TestUploader",
								},
								"files": []map[string]interface{}{
									{
										"file_id":   12345,
										"file_name": "The.Office.S03E07.srt",
									},
								},
							},
						},
					},
				}

				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(mockResponse)
				return
			}

			w.WriteHeader(http.StatusNotFound)
		}))
		defer server.Close()

		config := &Config{
			BaseURL:  server.URL,
			Username: "test",
			Password: "test",
		}

		client := NewOpenSubtitlesClient(config)
		
		params := &models.SearchParams{
			Query:    "The Office",
			Language: "en",
			Type:     "episode",
			Year:     2005,
			Season:   3,
			Episode:  7,
		}

		subtitles, err := client.Search(context.Background(), params)

		require.NoError(t, err)
		require.Len(t, subtitles, 1)

		subtitle := subtitles[0]
		assert.Equal(t, "test-id-123", subtitle.ID)
		assert.Equal(t, "en", subtitle.Language)
		assert.Equal(t, "The.Office.S03E07.720p.BluRay.x264", subtitle.ReleaseName)
		assert.Equal(t, "The.Office.S03E07.srt", subtitle.FileName)
		assert.Equal(t, "12345", subtitle.FileID)
		assert.Equal(t, "TestUploader", subtitle.Uploader)
		assert.Equal(t, 8.5, subtitle.Rating)
		assert.Equal(t, 1500, subtitle.Downloads)
		assert.Equal(t, 23.976, subtitle.FPS)
		assert.Equal(t, "srt", subtitle.SubFormat)

		expectedDate, _ := time.Parse("2006-01-02T15:04:05", "2023-01-15T10:30:00")
		assert.Equal(t, expectedDate, subtitle.UploadDate)
	})

	t.Run("search with minimal params", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/login" {
				response := LoginResponse{Token: "test-token", Status: 200}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
				return
			}

			if r.URL.Path == "/subtitles" {
				assert.Equal(t, "test movie", r.URL.Query().Get("query"))
				assert.Equal(t, "", r.URL.Query().Get("languages"))
				assert.Equal(t, "", r.URL.Query().Get("type"))

				response := map[string]interface{}{
					"data": []map[string]interface{}{},
				}

				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
				return
			}

			w.WriteHeader(http.StatusNotFound)
		}))
		defer server.Close()

		config := &Config{BaseURL: server.URL, Username: "test", Password: "test"}
		client := NewOpenSubtitlesClient(config)
		
		params := &models.SearchParams{Query: "test movie"}
		subtitles, err := client.Search(context.Background(), params)

		require.NoError(t, err)
		assert.Empty(t, subtitles)
	})

	t.Run("authentication error", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/login" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}))
		defer server.Close()

		config := &Config{BaseURL: server.URL, Username: "wrong", Password: "wrong"}
		client := NewOpenSubtitlesClient(config)
		
		params := &models.SearchParams{Query: "test"}
		_, err := client.Search(context.Background(), params)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "authentication required")
	})
}

func TestOpenSubtitlesClient_Download(t *testing.T) {
	t.Parallel()

	t.Run("successful download", func(t *testing.T) {
		t.Parallel()

		subtitleContent := "1\n00:00:01,000 --> 00:00:05,000\nHello World\n\n"
		var serverURL string

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/login" {
				response := LoginResponse{Token: "test-token", Status: 200}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
				return
			}

			if r.URL.Path == "/download" {
				assert.Equal(t, "POST", r.Method)

				var downloadReq DownloadRequest
				err := json.NewDecoder(r.Body).Decode(&downloadReq)
				require.NoError(t, err)
				assert.Equal(t, 12345, downloadReq.FileID)

				response := DownloadResponse{
					Link:     serverURL + "/subtitle-file",
					FileName: "test.srt",
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
				return
			}

			if r.URL.Path == "/subtitle-file" {
				w.Header().Set("Content-Type", "text/plain")
				w.Write([]byte(subtitleContent))
				return
			}

			w.WriteHeader(http.StatusNotFound)
		}))
		defer server.Close()
		serverURL = server.URL

		config := &Config{BaseURL: server.URL, Username: "test", Password: "test"}
		client := NewOpenSubtitlesClient(config)
		
		subtitle := &models.Subtitle{
			ID:     "test-id",
			FileID: "12345",
		}

		content, err := client.Download(context.Background(), subtitle)

		require.NoError(t, err)
		assert.Equal(t, subtitleContent, string(content))
	})

	t.Run("invalid file ID", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/login" {
				response := LoginResponse{Token: "test-token", Status: 200}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
				return
			}
			w.WriteHeader(http.StatusNotFound)
		}))
		defer server.Close()

		config := &Config{BaseURL: server.URL, Username: "test", Password: "test"}
		client := NewOpenSubtitlesClient(config)
		
		subtitle := &models.Subtitle{FileID: "invalid"}
		_, err := client.Download(context.Background(), subtitle)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid file ID")
	})

	t.Run("download limit exceeded", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/login" {
				response := LoginResponse{Token: "test-token", Status: 200}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
				return
			}

			if r.URL.Path == "/download" {
				w.WriteHeader(http.StatusNotAcceptable)
				response := DownloadResponse{Message: "Daily download limit exceeded"}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
				return
			}
		}))
		defer server.Close()

		config := &Config{BaseURL: server.URL, Username: "test", Password: "test"}
		client := NewOpenSubtitlesClient(config)
		
		subtitle := &models.Subtitle{FileID: "12345"}
		_, err := client.Download(context.Background(), subtitle)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "download limit exceeded")
	})
}