package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kaepa3/healthplanet"
	"golang.org/x/oauth2"
)

func loadEnv() (string, string) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_SECRET")
}

const (
	tokenFileName = ".token"
)

func getTokenFromFile() (*oauth2.Token, error) {
	f, err := os.Open(tokenFileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	buf, err := ioutil.ReadAll(f)
	if err != nil {

		return nil, err
	}
	var p oauth2.Token
	err = json.Unmarshal(buf, &p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}
func getTokenFromWeb(conf healthplanet.HealthPlanetConfig) (*oauth2.Token, error) {
	url := conf.AuthCodeURL("state")
	fmt.Println(url)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	authCode := scanner.Text()

	token, err := conf.Exchange(context.Background(), authCode)
	if err != nil {
		return nil, err
	}

	wf, err := os.Create(tokenFileName)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer wf.Close()
	encoder := json.NewEncoder(wf)
	if err := encoder.Encode(token); err != nil {
		log.Fatal(err)
	}
	return token, nil
}

func getToken(conf healthplanet.HealthPlanetConfig) (*oauth2.Token, error) {

	if exists(tokenFileName) {
		token, err := getTokenFromFile()
		if err == nil {
			return token, nil
		}
	}

	return getTokenFromWeb(conf)
}

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func main() {

	ctx := context.Background()
	clientID, clientSecret := loadEnv()

	conf := healthplanet.NewConfig(&healthplanet.HealthPlanetInit{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  "http://localhost",
		Scopes: []string{
			"innerscan",
		},
	})
	token, err := getToken(conf)
	if err != nil {
		fmt.Println(err)
		return
	}

	client, err := conf.GetClient(ctx, token)
	if err != nil {
		fmt.Println(err)
		return
	}

	opt := healthplanet.HealthPlanetOption{}
	resp, err := client.Get(healthplanet.Innerscan, &opt)
	if err != nil {
		fmt.Println(err)
		fmt.Println(resp)
	} else {
		fmt.Println(resp)
	}
}
