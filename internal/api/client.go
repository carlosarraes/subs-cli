package api

import (
	"context"

	"github.com/carlosarraes/subs-cli/pkg/models"
)

type Client interface {
	Search(ctx context.Context, params *models.SearchParams) ([]*models.Subtitle, error)
	Download(ctx context.Context, subtitle *models.Subtitle) ([]byte, error)
	Authenticate(ctx context.Context) error
}

type Config struct {
	APIKey    string
	UserAgent string
	BaseURL   string
	Username  string
	Password  string
}