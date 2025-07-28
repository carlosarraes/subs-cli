package api

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/carlosarraes/subs-cli/pkg/models"
)

const (
	DefaultBaseURL   = "https://api.opensubtitles.com/api/v1"
	DefaultUserAgent = "subs-cli/1.0"
)

type OpenSubtitlesClient struct {
	client *resty.Client
	config *Config
	token  string
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	User struct {
		AllowedDownloads int    `json:"allowed_downloads"`
		Level            string `json:"level"`
		UserID           int    `json:"user_id"`
		ExtInstalled     bool   `json:"ext_installed"`
		VIP              struct {
			ID             int    `json:"id"`
			Name           string `json:"name"`
			PointsSpent    int    `json:"points_spent"`
			RemainingDowns int    `json:"remaining_downloads"`
		} `json:"vip"`
	} `json:"user"`
	Token  string `json:"token"`
	Status int    `json:"status"`
}

type SearchResponse struct {
	TotalPages int `json:"total_pages"`
	TotalCount int `json:"total_count"`
	PerPage    int `json:"per_page"`
	Page       int `json:"page"`
	Data       []struct {
		ID         string `json:"id"`
		Type       string `json:"type"`
		Attributes struct {
			SubtitleID   string    `json:"subtitle_id"`
			Language     string    `json:"language"`
			DownloadCount int      `json:"download_count"`
			NewDownloadCount int  `json:"new_download_count"`
			HearingImpaired bool  `json:"hearing_impaired"`
			HD               bool  `json:"hd"`
			FPS              float64 `json:"fps"`
			Votes            int   `json:"votes"`
			Ratings          float64 `json:"ratings"`
			FromTrusted      bool  `json:"from_trusted"`
			ForeignPartsOnly bool  `json:"foreign_parts_only"`
			AITranslated     bool  `json:"ai_translated"`
			MachineTranslated bool `json:"machine_translated"`
			UploadDate       string `json:"upload_date"`
			Release          string `json:"release"`
			Comments         string `json:"comments"`
			LegacySubtitleID int   `json:"legacy_subtitle_id"`
			Uploader         struct {
				UploaderID int    `json:"uploader_id"`
				Name       string `json:"name"`
				Rank       string `json:"rank"`
			} `json:"uploader"`
			FeatureDetails struct {
				FeatureID   int    `json:"feature_id"`
				FeatureType string `json:"feature_type"`
				Year        int    `json:"year"`
				Title       string `json:"title"`
				MovieName   string `json:"movie_name"`
				IMDBID      int    `json:"imdb_id"`
				TMDBID      int    `json:"tmdb_id"`
			} `json:"feature_details"`
			URL       string `json:"url"`
			RelatedLinks []struct {
				Label string `json:"label"`
				URL   string `json:"url"`
				ImgURL string `json:"img_url"`
			} `json:"related_links"`
			Files []struct {
				FileID int `json:"file_id"`
				CDID   int `json:"cd_number"`
				FileName string `json:"file_name"`
			} `json:"files"`
		} `json:"attributes"`
	} `json:"data"`
}

type DownloadRequest struct {
	FileID int `json:"file_id"`
}

type DownloadResponse struct {
	Link       string `json:"link"`
	FileName   string `json:"file_name"`
	Requests   int    `json:"requests"`
	Remaining  int    `json:"remaining"`
	Message    string `json:"message"`
	ResetTime  string `json:"reset_time"`
	ResetTimeUTC string `json:"reset_time_utc"`
}

func NewOpenSubtitlesClient(config *Config) *OpenSubtitlesClient {
	if config.BaseURL == "" {
		config.BaseURL = DefaultBaseURL
	}
	if config.UserAgent == "" {
		config.UserAgent = DefaultUserAgent
	}

	client := resty.New()
	client.SetBaseURL(config.BaseURL)
	client.SetHeader("User-Agent", config.UserAgent)
	if config.APIKey != "" {
		client.SetHeader("Api-Key", config.APIKey)
	}
	client.SetTimeout(30 * time.Second)

	return &OpenSubtitlesClient{
		client: client,
		config: config,
	}
}

func (c *OpenSubtitlesClient) Authenticate(ctx context.Context) error {
	if c.config.Username == "" || c.config.Password == "" {
		return fmt.Errorf("username and password are required for authentication")
	}

	loginReq := LoginRequest{
		Username: c.config.Username,
		Password: c.config.Password,
	}

	var loginResp LoginResponse
	resp, err := c.client.R().
		SetContext(ctx).
		SetBody(loginReq).
		SetResult(&loginResp).
		Post("/login")

	if err != nil {
		return fmt.Errorf("authentication request failed: %w", err)
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("authentication failed with status %d: %s", resp.StatusCode(), resp.String())
	}

	if loginResp.Status != 200 {
		return fmt.Errorf("authentication failed: invalid credentials")
	}

	c.token = loginResp.Token
	c.client.SetAuthToken(c.token)

	return nil
}

func (c *OpenSubtitlesClient) Search(ctx context.Context, params *models.SearchParams) ([]*models.Subtitle, error) {
	if c.token == "" {
		if err := c.Authenticate(ctx); err != nil {
			return nil, fmt.Errorf("authentication required: %w", err)
		}
	}

	request := c.client.R().SetContext(ctx)
	
	if params.Query != "" {
		request = request.SetQueryParam("query", params.Query)
	}
	
	if params.Language != "" {
		request = request.SetQueryParam("languages", params.Language)
	}
	
	if params.Type != "" {
		request = request.SetQueryParam("type", params.Type)
	}
	
	if params.Year > 0 {
		request = request.SetQueryParam("year", strconv.Itoa(params.Year))
	}
	
	if params.Season > 0 {
		request = request.SetQueryParam("season_number", strconv.Itoa(params.Season))
	}
	
	if params.Episode > 0 {
		request = request.SetQueryParam("episode_number", strconv.Itoa(params.Episode))
	}
	
	if params.MovieHash != "" {
		request = request.SetQueryParam("moviehash", params.MovieHash)
	}

	var searchResp SearchResponse
	resp, err := request.
		SetResult(&searchResp).
		Get("/subtitles")

	if err != nil {
		return nil, fmt.Errorf("search request failed: %w", err)
	}

	if resp.StatusCode() == 401 {
		c.token = ""
		return nil, fmt.Errorf("authentication expired, please retry")
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("search failed with status %d: %s", resp.StatusCode(), resp.String())
	}

	subtitles := make([]*models.Subtitle, 0, len(searchResp.Data))
	for _, item := range searchResp.Data {
		attrs := item.Attributes
		
		uploadDate, _ := time.Parse("2006-01-02T15:04:05", attrs.UploadDate)
		
		var fileName, fileID string
		if len(attrs.Files) > 0 {
			fileName = attrs.Files[0].FileName
			fileID = strconv.Itoa(attrs.Files[0].FileID)
		}
		
		subtitle := &models.Subtitle{
			ID:          item.ID,
			Language:    attrs.Language,
			ReleaseName: attrs.Release,
			FileName:    fileName,
			FileID:      fileID,
			Uploader:    attrs.Uploader.Name,
			Rating:      attrs.Ratings,
			Downloads:   attrs.DownloadCount,
			UploadDate:  uploadDate,
			FPS:         attrs.FPS,
			SubFormat:   "srt",
		}
		
		subtitles = append(subtitles, subtitle)
	}

	return subtitles, nil
}

func (c *OpenSubtitlesClient) Download(ctx context.Context, subtitle *models.Subtitle) ([]byte, error) {
	if c.token == "" {
		if err := c.Authenticate(ctx); err != nil {
			return nil, fmt.Errorf("authentication required: %w", err)
		}
	}

	fileID, err := strconv.Atoi(subtitle.FileID)
	if err != nil {
		return nil, fmt.Errorf("invalid file ID: %s", subtitle.FileID)
	}

	downloadReq := DownloadRequest{
		FileID: fileID,
	}

	var downloadResp DownloadResponse
	resp, err := c.client.R().
		SetContext(ctx).
		SetBody(downloadReq).
		SetResult(&downloadResp).
		Post("/download")

	if err != nil {
		return nil, fmt.Errorf("download request failed: %w", err)
	}

	if resp.StatusCode() == 401 {
		c.token = ""
		return nil, fmt.Errorf("authentication expired, please retry")
	}

	if resp.StatusCode() == 406 {
		return nil, fmt.Errorf("download limit exceeded: %s", downloadResp.Message)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("download failed with status %d: %s", resp.StatusCode(), resp.String())
	}

	if downloadResp.Link == "" {
		return nil, fmt.Errorf("no download link provided")
	}

	fileResp, err := c.client.R().
		SetContext(ctx).
		Get(downloadResp.Link)

	if err != nil {
		return nil, fmt.Errorf("failed to download subtitle file: %w", err)
	}

	if fileResp.StatusCode() != 200 {
		return nil, fmt.Errorf("subtitle file download failed with status %d", fileResp.StatusCode())
	}

	return fileResp.Body(), nil
}
