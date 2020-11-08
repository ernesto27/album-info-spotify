package main

import (
	"album-info-spotify/spotify"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/godbus/dbus"
)

type Album struct {
	Artist      string `json:"strArtist"`
	Name        string `json:"strAlbum"`
	Description string `json:"strDescriptionEN"`
	ReleaseYear string `json:"intYearReleased"`
	Style       string `json:"strStyle"`
	Review      string `json:"strReview"`
	Label       string `json:"strLabel"`
}

type ResponseAlbum struct {
	Album []Album
}

var nameBand string
var albumBand string

func main() {
	if len(os.Args) < 3 {
		conn := getConn()
		var meta *spotify.SpotifyMetadata
		meta, err := spotify.GetMetadataSpotify(conn)
		if err != nil {
			log.Fatalf("failed getting metadata, err: %s", err.Error())
		}
		nameBand = strings.ReplaceAll(meta.Artist[0], " ", "%20")
		albumBand = strings.ReplaceAll(meta.Album, " ", "%20")

		// fmt.Println("You have to pass bandName and albumName arguments")
		// fmt.Println("$ go run main.go megadeth 'rust in peace' ")
		// os.Exit(1)
	} else {
		// Search by band , album enter by user
		nameBand = os.Args[1]
		albumBand = os.Args[2]
		nameBand = strings.ReplaceAll(nameBand, " ", "%20")
		albumBand = strings.ReplaceAll(albumBand, " ", "%20")
	}

	album := new(ResponseAlbum)
	getJson("https://www.theaudiodb.com/api/v1/json/1/searchalbum.php?s="+nameBand+"&a="+albumBand, album)
	if len(album.Album) == 0 {
		fmt.Println("Album not found :(")
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println(album.Album[0].Artist + " - " + album.Album[0].Name)
	fmt.Println(album.Album[0].ReleaseYear + " - " + album.Album[0].Style)
	if album.Album[0].Label != "" {
		fmt.Println("Label: " + album.Album[0].Label)
	}
	fmt.Println()
	fmt.Println("----------------------------------------------------------------------------------------------------------------------------------------")
	fmt.Println()

	fmt.Println("DESCRIPTION:")
	fmt.Println(string(album.Album[0].Description))

	fmt.Println()
	fmt.Println("----------------------------------------------------------------------------------------------------------------------------------------")
	fmt.Println()

	if album.Album[0].Review != "" {
		fmt.Println("REVIEW:")
		fmt.Println(string(album.Album[0].Review))
	}

	fmt.Println()
	fmt.Println("----------------------------------------------------------------------------------------------------------------------------------------")

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

func getConn() *dbus.Conn {
	conn, err := dbus.SessionBus()
	if err != nil {
		log.Fatal(err)
	}
	return conn
}
