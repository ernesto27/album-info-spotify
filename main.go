package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

type Album struct {
	Description string `json:"strDescriptionEN"`
	ReleaseYear string `json:"intYearReleased"`
	Style       string `json:"strStyle"`
}

type ResponseAlbum struct {
	Album []Album
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("You have to pass bandName and albumName arguments")
		fmt.Println("$ go run main.go megadeth 'rust in peace' ")
		os.Exit(1)
	}

	nameBand := strings.ReplaceAll(os.Args[1], " ", "%20")
	albumBand := strings.ReplaceAll(os.Args[2], " ", "%20")

	album := new(ResponseAlbum)
	getJson("https://www.theaudiodb.com/api/v1/json/1/searchalbum.php?s="+nameBand+"&a="+albumBand, album)
	if len(album.Album) == 0 {
		fmt.Println("Album not found :(")
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println(os.Args[1] + " - " + os.Args[2])
	fmt.Println(album.Album[0].ReleaseYear + " - " + album.Album[0].Style)
	fmt.Println()
	fmt.Println("----------------------------------------------------------------------------------------------------------------------------------------")
	fmt.Println("----------------------------------------------------------------------------------------------------------------------------------------")
	fmt.Println("----------------------------------------------------------------------------------------------------------------------------------------")
	fmt.Println()

	fmt.Println(string(album.Album[0].Description))

	fmt.Println()
	fmt.Println("----------------------------------------------------------------------------------------------------------------------------------------")
	fmt.Println("----------------------------------------------------------------------------------------------------------------------------------------")
	fmt.Println("----------------------------------------------------------------------------------------------------------------------------------------")
	fmt.Println()

}

var myClient = &http.Client{Timeout: 10 * time.Second}

func getJson(url string, target interface{}) error {
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}
