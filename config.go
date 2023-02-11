package healthplanet

import (
	"context"

	"golang.org/x/oauth2"
)

type HealthPlanetInit struct {
	ClientID     string
	ClientSecret string
	Scopes       []string
	RedirectURL  string
}

type HealthPlanetConfig struct {
	config oauth2.Config
}

func NewConfig(c *HealthPlanetInit) HealthPlanetConfig {
	return HealthPlanetConfig{
		config: oauth2.Config{
			ClientID:     c.ClientID,
			ClientSecret: c.ClientSecret,
			RedirectURL:  c.RedirectURL,
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://www.healthplanet.jp/oauth/auth",
				TokenURL: "https://www.healthplanet.jp/oauth/token",
			},
			Scopes: c.Scopes,
		},
	}
}
func (c *HealthPlanetConfig) AuthCodeURL(state string) string {
	return c.config.AuthCodeURL(state)
}
func (c *HealthPlanetConfig) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return c.config.Exchange(ctx, code)
}

func (c *HealthPlanetConfig) GetClient(ctx context.Context, token *oauth2.Token) (*HealthPlanetClient, error) {
	client := c.config.Client(ctx, token)
	return &HealthPlanetClient{client: client, token: token}, nil
}
