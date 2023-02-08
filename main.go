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

func main() {

	ctx := context.Background()
	clientID, clientSecret := loadEnv()
	conf := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  "http://localhost",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://www.healthplanet.jp/oauth/auth",
			TokenURL: "https://www.healthplanet.jp/oauth/token",
		},
		Scopes: []string{
			"innerscan",
		},
	}
	url := conf.AuthCodeURL("state")
	server()
	fmt.Println(url)

	scanner := bufio.NewScanner(os.Stdin) // 標準入力を受け付けるスキャナ
	scanner.Scan()                        // １行分の入力を取得する
	authCode := scanner.Text()
	token, err := conf.Exchange(ctx, authCode)
	if err != nil {
		fmt.Println("----------------------------------------")
		fmt.Println(err)
	}
	client := conf.Client(context.Background(), token)
	resp, err := client.Get("https://www.healthplanet.jp/status/innerscan.json?access_token=" + token.AccessToken)
	if err != nil {
		fmt.Println("-----------")
		fmt.Println(err)
		fmt.Println(resp)
	} else {
		fmt.Println(resp)
	}

}
func server() {
	go func() {
		http.HandleFunc("/start", start)
		http.HandleFunc("/callback", callback)
		err := http.ListenAndServe("localhost:51200", nil)
		if err != nil {
			log.Fatal(err)
		}
	}()
}
func start(w http.ResponseWriter, req *http.Request) {
	fmt.Println(w)
}
func callback(w http.ResponseWriter, req *http.Request) {
	fmt.Println(w)
}
