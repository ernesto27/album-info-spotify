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

func main() {
	conn := getConn()
	var meta *spotify.SpotifyMetadata
	meta, err := spotify.GetMetadataSpotify(conn)
	if err != nil {
		fmt.Println("Seems that you don't have the spotify app desktop installed  or is not open :(")
		log.Fatalf("failed getting metadata, err: %s", err.Error())
	}
	fmt.Println(meta)

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
				<div class="col-md-4">
					<img src="` + albumInfo[0] + `" class="img-fluid img-thumbnail" alt="Responsive image">
				</div>
				<div class="col-md-4">
					<img src="` + albumInfo[1] + `" class="img-fluid img-thumbnail" alt="Responsive image">
				</div>
				<div class="col-md-4">
					<img src="` + albumInfo[2] + `" class="img-fluid img-thumbnail" alt="Responsive image">
				</div>
			</div>

			<br>
			<p class="font-weight-bold">Description album: </p>
			<p class="text-justify">` + albumInfo[3] + `</p>

			<br>
			<p class="font-weight-bold">Review album: </p>
			<p class="text-justify">` + albumInfo[4] + `</p>

			<hr />
			<!-- TRACK INFO -->
			<br />
			<h3>Track info</h3>
			<p class="text-justify">` + trackInfo[0] + `</p>

			<div class="row">
				<div class="col-md-4">
					<img src="` + trackInfo[2] + `" class="img-fluid img-thumbnail" alt="Responsive image">
				</div>
			</div>
			<br />
			<iframe width="560" height="315" src="https://www.youtube.com/embed/` + trackInfo[1] + `" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>

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
