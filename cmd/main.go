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

func getToken(conf healthplanet.HealthPlanetConfig) *oauth2.Token {

	if exists(tokenFileName) {
		f, err := os.Open(tokenFileName)
		defer f.Close()
		if err == nil {
			if buf, err := ioutil.ReadAll(f); err == nil {
				var p oauth2.Token
				if err := json.Unmarshal(buf, &p); err == nil {
					return &p
				} else {
					fmt.Print(err)
				}
			}
		}
	}

	url := conf.AuthCodeURL("state")
	fmt.Println(url)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	authCode := scanner.Text()

	token, err := conf.Exchange(context.Background(), authCode)
	if err != nil {
		return nil
	}

	wf, err := os.Create(tokenFileName)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer wf.Close()
	encoder := json.NewEncoder(wf)
	if err := encoder.Encode(token); err != nil {
		log.Fatal(err)
	}
	return token
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
	token := getToken(conf)

	client, err := conf.GetClient(ctx, token)
	if err != nil {
		fmt.Println(err)
		return
	}

	opt := healthplanet.HealthPlanetOption{}
	resp, err := client.Get(healthplanet.Innerscan, &opt)
	if err != nil {
		fmt.Println("-----------")
		fmt.Println(err)
		fmt.Println(resp)
	} else {
		fmt.Println(resp)
	}
}
