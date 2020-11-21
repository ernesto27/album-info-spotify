package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

var apiURL string = "https://www.theaudiodb.com/api/v1/json/1/"
var myClient = &http.Client{Timeout: 10 * time.Second}

type Album struct {
	Artist      string `json:"strArtist"`
	Name        string `json:"strAlbum"`
	Description string `json:"strDescriptionEN"`
	ReleaseYear string `json:"intYearReleased"`
	Style       string `json:"strStyle"`
	Review      string `json:"strReview"`
	Label       string `json:"strLabel"`
	ThumbFront  string `json:"strAlbumThumb"`
	ThumbBack   string `json:"strAlbumThumbBack"`
	ThumbCD     string `json:"strAlbumCDart"`
}

type ResponseAlbum struct {
	Album []Album
}

type Track struct {
	Description string `json:"strDescriptionEN"`
	YoutubeURL  string `json:"strMusicVid"`
	Thumb       string `json:"strTrackThumb"`
}

type ResponseTrack struct {
	Track []Track `json:"track"`
}

type Band struct {
	Name       string `json:"strArtist"`
	Biograhpy  string `json:"strBiographyEN"`
	FormedYear string `json:"intFormedYear"`
	Country    string `json:"strCountry"`
}

type ResponseBand struct {
	Band []Band `json:"artists"`
}

func GetAlbumInfo(nameBand string, albumBand string, wg *sync.WaitGroup, albumChannel chan []string) {
	defer wg.Done()

	album := new(ResponseAlbum)
	cleanInfo := cleanStrings(nameBand, albumBand, "")
	nameBand = cleanInfo[0]
	albumBand = cleanInfo[1]

	fmt.Println(apiURL + "searchalbum.php?s=" + nameBand + "&a=" + albumBand)
	getJson(apiURL+"searchalbum.php?s="+nameBand+"&a="+albumBand, album)

	resp := []string{
		album.Album[0].ThumbFront,
		album.Album[0].ThumbBack,
		album.Album[0].ThumbCD,
		album.Album[0].Description,
		album.Album[0].Review,
		album.Album[0].Artist,
		album.Album[0].Name,
		album.Album[0].ReleaseYear,
		album.Album[0].Style,
		album.Album[0].Label,
	}

	albumChannel <- resp
}

func GetTrackInfo(nameBand string, trackName string, wg *sync.WaitGroup, trackChannel chan []string) {
	defer wg.Done()
	var track = new(ResponseTrack)
	cleanInfo := cleanStrings(nameBand, "", trackName)
	nameBand = cleanInfo[0]
	trackName = cleanInfo[2]

	fmt.Println(apiURL + "searchtrack.php?s=" + nameBand + "&t=" + trackName)
	getJson(apiURL+"searchtrack.php?s="+nameBand+"&t="+trackName, track)

	if len(track.Track) == 0 {
		trackChannel <- []string{}
		return
	}

	// Get id url youtube
	var idVideo string
	if track.Track[0].YoutubeURL != "" {
		idVideo = strings.Split(track.Track[0].YoutubeURL, "v=")[1]
	}

	resp := []string{
		track.Track[0].Description,
		idVideo,
		track.Track[0].Thumb,
		track.Track[0].YoutubeURL,
	}
	trackChannel <- resp
}

func GetBandInfo(nameBand string, wg *sync.WaitGroup, bandChannel chan []string) {
	defer wg.Done()
	var band = new(ResponseBand)
	cleanInfo := cleanStrings(nameBand, "", "")
	nameBand = cleanInfo[0]
	fmt.Println(apiURL + "search.php?s=" + nameBand)
	getJson(apiURL+"search.php?s="+nameBand, band)
	// if len(band.Band) == 0 {
	// 	return []string{}, errors.New("Band not found")
	// }
	resp := []string{
		band.Band[0].Biograhpy,
		band.Band[0].FormedYear,
		band.Band[0].Country,
	}

	bandChannel <- resp
}

func cleanStrings(nameBand string, albumBand string, trackName string) []string {
	trackName = strings.ReplaceAll(trackName, " ", "%20")
	nameBand = strings.ReplaceAll(nameBand, " ", "%20")
	albumBand = strings.ReplaceAll(albumBand, " ", "%20")

	m := regexp.MustCompile("\\(.*\\)")
	albumBand = m.ReplaceAllString(albumBand, "")
	albumBand = strings.Trim(albumBand, "%20")

	resp := []string{
		nameBand,
		albumBand,
		trackName,
	}
	return resp
}

func getJson(url string, target interface{}) error {
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}
