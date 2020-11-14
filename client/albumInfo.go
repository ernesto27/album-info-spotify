package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
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

func GetAlbumInfo(nameBand string, albumBand string) ([]string, error) {
	album := new(ResponseAlbum)
	cleanInfo := cleanStrings(nameBand, albumBand, "")
	nameBand = cleanInfo[0]
	albumBand = cleanInfo[1]

	fmt.Println(apiURL + "searchalbum.php?s=" + nameBand + "&a=" + albumBand)
	getJson(apiURL+"searchalbum.php?s="+nameBand+"&a="+albumBand, album)
	if len(album.Album) == 0 {
		return []string{}, errors.New("Album not found")
	}
	resp := []string{
		album.Album[0].ThumbFront,
		album.Album[0].ThumbBack,
		album.Album[0].ThumbCD,
		album.Album[0].Description,
		album.Album[0].Review,
	}
	return resp, nil
}

func GetTrackInfo(nameBand string, trackName string) ([]string, error) {
	var track = new(ResponseTrack)
	cleanInfo := cleanStrings(nameBand, "", trackName)
	nameBand = cleanInfo[0]
	trackName = cleanInfo[2]

	fmt.Println(apiURL + "searchtrack.php?s=" + nameBand + "&t=" + trackName)
	getJson(apiURL+"searchtrack.php?s="+nameBand+"&t="+trackName, track)

	if len(track.Track) == 0 {
		return []string{}, errors.New("Track not found :(")
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
	}

	return resp, nil
}

func GetBandInfo(nameBand string) ([]string, error) {
	var band = new(ResponseBand)
	cleanInfo := cleanStrings(nameBand, "", "")
	nameBand = cleanInfo[0]
	fmt.Println(apiURL + "search.php?s=" + nameBand)
	getJson(apiURL+"search.php?s="+nameBand, band)
	if len(band.Band) == 0 {
		return []string{}, errors.New("Band not found")
	}
	resp := []string{
		band.Band[0].Biograhpy,
		band.Band[0].FormedYear,
		band.Band[0].Country,
	}
	return resp, nil
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
