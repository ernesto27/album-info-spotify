package main

import (
	"album-info-spotify/client"
	"album-info-spotify/spotify"
	"fmt"
	"log"
	"net/url"

	"github.com/godbus/dbus"
	"github.com/zserge/lorca"
)

func renderImage(imageURL string) string {
	if imageURL == "" {
		return ""
	}
	return `<div class="col-md-4">
		<img src="` + imageURL + `" class="img-fluid img-thumbnail" alt="">
	</div>`
}

func renderReview(text string) string {
	if text == "" {
		return ""
	}
	return `<p class="font-weight-bold">Review album: </p>
			<p class="text-justify">` + text + `</p> `
}

func renderTrackInfo(text string, videoURL string, thumb string) string {
	var resp string = ""
	if text != "" {
		resp += `
		<h3>Track info</h3>
		<p class="text-justify">` + text + `</p>`
	}

	if thumb != "" {
		resp += `
		<div class="row">
			<div class="col-md-4">
				<img src="` + thumb + `" class="img-fluid img-thumbnail" alt="Responsive image">
			</div>
		</div>`
	}

	if videoURL != "" {
		resp += `<br />
				<iframe width="560" height="315" src="https://www.youtube.com/embed/` + videoURL + `" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>`
	}

	return resp

}

func main() {
	conn := getConn()
	var meta *spotify.SpotifyMetadata
	meta, err := spotify.GetMetadataSpotify(conn)
	if err != nil {
		fmt.Println("Seems that you don't have the spotify app desktop installed  or is not open :(")
		log.Fatalf("failed getting metadata, err: %s", err.Error())
	}

	albumInfo, err := client.GetAlbumInfo(meta.Artist[0], meta.Album)
	if err != nil {
		panic(err)
	}

	trackInfo, err := client.GetTrackInfo(meta.Artist[0], meta.Title)
	if err != nil {
		panic(err)
	}

	bandInfo, err := client.GetBandInfo(meta.Artist[0])
	if err != nil {
		panic(err)
	}

	// Create UI with data URI
	var htmlBody string = `
	<!doctype html>
	<html lang="en">
	<head>
		<!-- Required meta tags -->
		<meta charset="utf-8">
		<meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">

		<!-- Bootstrap CSS -->
		<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bootstrap@4.5.3/dist/css/bootstrap.min.css" integrity="sha384-TX8t27EcRE3e/ihU7zmQxVncDAy5uIKz4rEkgIXeMed4M0jlfIDPvg6uqKI2xXr2" crossorigin="anonymous">

		<title>` + meta.Artist[0] + `  ` + meta.Album + `</title>
	</head>
	<body>
		<div class="container mt-5" >
			<h2>SPOTIFY - ALBUM BAND INFO</h2>

			<p class="font-weight-bold">` + meta.Artist[0] + `  ` + meta.Album + `</p>
			<div class="row">
				` + renderImage(albumInfo[0]) + renderImage(albumInfo[1]) + renderImage(albumInfo[2]) + `
			</div>

			<br>
			<p class="font-weight-bold">Description album: </p>
			<p class="text-justify">` + albumInfo[3] + `</p>

			<br>
			` + renderReview(albumInfo[4]) + `

			<hr />
			<!-- TRACK INFO -->
			<br />
			` + renderTrackInfo(trackInfo[0], trackInfo[1], trackInfo[2]) + `
	
			<hr />

			<!-- BAND INFO -->
			<p class="font-weight-bold">BIO: </p>
			<p class="text-justify">` + bandInfo[0] + `</p>

			<p>FORM YEAR: ` + bandInfo[1] + ` ` + bandInfo[2] + `</p>

		</div>


	</body>
	</html> 	
	`
	ui, _ := lorca.New("data:text/html,"+url.PathEscape(htmlBody), "", 900, 600)
	defer ui.Close()
	// Create a GoLang function callable from JS
	ui.Bind("hello", func() string { return "World!" })
	// Call above `hello` function then log to the JS console
	ui.Eval("hello().then( (x) => { console.log(x) })")
	// Wait until UI window is closed
	<-ui.Done()
}

func getConn() *dbus.Conn {
	conn, err := dbus.SessionBus()
	if err != nil {
		log.Fatal(err)
	}
	return conn
}
