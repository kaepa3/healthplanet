package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

func loadEnv() (string, string) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_SECRET")
}

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
func (c *HealthPlanetConfig) GetClient(ctx context.Context, code string) (*HealthPlanetClient, error) {
	token, err := c.config.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	client := c.config.Client(context.Background(), token)
	return &HealthPlanetClient{client: client, token: token}, nil
}

type HealthPlanetClient struct {
	client *http.Client
	ctx    context.Context
	token  *oauth2.Token
}

const (
	apiRoot = "https://www.healthplanet.jp/status/innerscan"
)

func (c *HealthPlanetClient) Get() (*http.Response, error) {
	request := apiRoot + c.getReturnType() + "?access_token=" + c.token.AccessToken
	fmt.Println(request)
	return c.client.Get(request)
}

func (c *HealthPlanetClient) getReturnType() string {
	return ".json"
}

func main() {

	ctx := context.Background()
	clientID, clientSecret := loadEnv()

	conf := NewConfig(&HealthPlanetInit{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  "http://localhost",
		Scopes: []string{
			"innerscan",
		},
	})

	url := conf.AuthCodeURL("state")
	fmt.Println(url)

	scanner := bufio.NewScanner(os.Stdin) // 標準入力を受け付けるスキャナ
	scanner.Scan()                        // １行分の入力を取得する
	authCode := scanner.Text()

	client, err := conf.GetClient(ctx, authCode)
	if err != nil {
		return
	}

	resp, err := client.Get()
	if err != nil {
		fmt.Println("-----------")
		fmt.Println(err)
		fmt.Println(resp)
	} else {
		fmt.Println(resp)
	}
}
