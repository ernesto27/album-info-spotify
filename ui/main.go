package main

import (
	"album-info-spotify/client"
	"album-info-spotify/spotify"
	"fmt"
	"log"
	"net/url"

	"github.com/zserge/lorca"
)

func renderImage(imageURL string) string {
	if imageURL == "" {
		return ""
	}
	return `<div class="col-md-4">
		<img id="album-img-" src="` + imageURL + `" class="img-fluid img-thumbnail" alt="">
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
	meta, err := spotify.GetMetadataSpotify()
	if err != nil {
		fmt.Println("Seems that you don't have the spotify app desktop installed  or is not open :(")
		log.Fatalf("failed getting metadata, err: %s", err.Error())
	}

	artistName := meta.ArtistName[0]

	albumInfo, err := client.GetAlbumInfo(artistName, meta.AlbumName)
	if err != nil {
		panic(err)
	}

	trackInfo, err := client.GetTrackInfo(artistName, meta.TrackName)
	if err != nil {
		panic(err)
	}

	bandInfo, err := client.GetBandInfo(artistName)
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

			<title>` + artistName + `  ` + meta.AlbumName + `</title>
		</head>
		<body>
			<div class="container mt-5" >

				
				<h2>SPOTIFY - ALBUM BAND INFO</h2>
				<button type="button" class="btn btn-primary" id="button">REFRESH DATA</button>
				<br /> <br />
				<p class="font-weight-bold" id="title-artist-album">` + artistName + `  ` + meta.AlbumName + `</p>
				<div class="row" id="wrapper-album-images">
					` + renderImage(albumInfo[0]) + renderImage(albumInfo[1]) + renderImage(albumInfo[2]) + `
				</div>

				<br>
				<p class="font-weight-bold">Description album: </p>
				<p class="text-justify" id="description-album">` + albumInfo[3] + `</p>

				<br>
				<div id="review-album">
				` + renderReview(albumInfo[4]) + `
				</div>
				<hr />
				<!-- TRACK INFO -->
				<br />
				<div id="track-info">
				` + renderTrackInfo(trackInfo[0], trackInfo[1], trackInfo[2]) + `
				</div>
				<hr />

				<!-- BAND INFO -->
				<p class="font-weight-bold">BIO: </p>
				<p class="text-justify" id="artist-bio">` + bandInfo[0] + `</p>

				<p>FORM YEAR: <span id="artist-year">` + bandInfo[1] + ` ` + bandInfo[2] + `</span></p>

			</div>

			<script>
				document.getElementById('button').addEventListener('click', function(){
					refresh().then( (data) => { 
						console.log(data) 
						document.getElementById('title-artist-album').innerHTML = data[0] + ' ' + data[1] ;

						var imagesAlbumHTML = '';
						
						var wrapperAlbumImages = document.getElementById('wrapper-album-images');
						wrapperAlbumImages.innerHTML = data[2];

						var descriptionAlbum = document.getElementById('description-album');
						descriptionAlbum.innerHTML = data[3];

						var reviewAlbum = document.getElementById('review-album');
						reviewAlbum.innerHTML = data[4];

						var trackInfo = document.getElementById('track-info');
						trackInfo.innerHTML = data[5];

						var artistBio = document.getElementById('artist-bio');
						artistBio.innerHTML = data[6];

						var artistYear = document.getElementById('artist-year');
						artistYear.innerHTML = data[7];
					})
				})
			</script>
		</body>
		</html>
		`
	ui, _ := lorca.New("data:text/html,"+url.PathEscape(htmlBody), "", 900, 700)
	defer ui.Close()

	// Create a GoLang function callable from JS
	ui.Bind("hello", func() string { return "World!" })
	ui.Bind("refresh", func() []string {
		meta, err := spotify.GetMetadataSpotify()
		if err != nil {
			fmt.Println("Seems that you don't have the spotify app desktop installed  or is not open :(")
			log.Fatalf("failed getting metadata, err: %s", err.Error())
		}

		artistName := meta.ArtistName[0]
		albumInfo, err := client.GetAlbumInfo(artistName, meta.AlbumName)
		if err != nil {
			panic(err)
		}

		var albumImageHTML string = ""
		albumImageHTML += renderImage(albumInfo[0])
		albumImageHTML += renderImage(albumInfo[1])
		albumImageHTML += renderImage(albumInfo[2])

		albumReview := renderReview(albumInfo[4])

		trackInfo, err := client.GetTrackInfo(artistName, meta.TrackName)
		if err != nil {
			panic(err)
		}

		trackInfoHTML := renderTrackInfo(trackInfo[0], trackInfo[1], trackInfo[2])

		bandInfo, err := client.GetBandInfo(artistName)
		if err != nil {
			panic(err)
		}

		return []string{
			artistName,
			meta.AlbumName,
			albumImageHTML,
			albumInfo[3],
			albumReview,
			trackInfoHTML,
			bandInfo[0],
			bandInfo[1],
			bandInfo[2],
		}
	})
	// Call above `hello` function then log to the JS console
	ui.Eval("hello().then( (x) => { console.log(x) })")
	// Wait until UI window is closed
	<-ui.Done()
}
