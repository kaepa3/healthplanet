package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kaepa3/healthplanet"
)

func loadEnv() (string, string) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_SECRET")
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

	url := conf.AuthCodeURL("state")
	fmt.Println(url)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	authCode := scanner.Text()

	client, err := conf.GetClient(ctx, authCode)
	if err != nil {
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
