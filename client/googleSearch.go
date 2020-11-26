package client

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	googlesearch "github.com/rocketlaunchr/google-search"
)

type Image struct {
	Link string `json:"link"`
}

type Items struct {
	Images []Image `json:"items"`
}

func GetImagesBand(artistName string, year string) *Items {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	req, err := http.NewRequest("GET", "https://www.googleapis.com/customsearch/v1", nil)
	if err != nil {
		log.Print(err)
	}

	q := req.URL.Query()
	q.Add("key", os.Getenv("GOOGLE_SEARCH_APIKEY"))
	q.Add("cx", os.Getenv("GOOGLE_SEARCH_CX"))
	q.Add("q", artistName+" band "+year)
	q.Add("searchType", "image")

	req.URL.RawQuery = q.Encode()

	var myClient = &http.Client{Timeout: 10 * time.Second}
	r, err := myClient.Get(req.URL.String())
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	var items = new(Items)
	json.NewDecoder(r.Body).Decode(items)

	return items

}

func GetWikipediaLink(nameBand string, albumBand string) ([]googlesearch.Result, error) {

	ctx := context.Background()
	return googlesearch.Search(ctx, nameBand+" "+albumBand+" wikipedia")

}
