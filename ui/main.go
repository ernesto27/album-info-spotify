package main

import (
	"album-info-spotify/client"
	"album-info-spotify/spotify"
	"fmt"
	"log"
	"net/url"
	"strings"
	"sync"

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

func renderTrackInfo(trackInfo []string) string {
	if len(trackInfo) == 0 {
		return ""
	}

	text := trackInfo[0]
	videoURL := trackInfo[1]
	thumb := trackInfo[2]

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

func renderTitle(albumInfo *client.ResponseAlbum) string {
	var resp string = ""
	resp += albumInfo.Album[0].Artist + ` - ` + albumInfo.Album[0].Name
	return resp
}

func renderAlbumMetadata(albumInfo *client.ResponseAlbum) string {
	var resp string = ""

	resp += `<p class="font-weight-bold" id="title-artist-album">` + albumInfo.Album[0].Artist +
		` -  ` + albumInfo.Album[0].Name + ` -  ` + albumInfo.Album[0].ReleaseYear + `</p>`

	if albumInfo.Album[0].Label != "" {
		resp += `<p class="font-weight-bold">Label: ` + albumInfo.Album[0].Label + `</p>`
	}
	if albumInfo.Album[0].Style != "" {
		resp += `<p class="font-weight-bold">Genre: ` + albumInfo.Album[0].Style + `</p>`
	}

	if albumInfo.Album[0].Score != "" {
		resp += `<p class="font-weight-bold">Score: ` + albumInfo.Album[0].Score + `</p>`
	}

	return resp
}

func renderBandAlbumImages(items client.Items) string {
	var resp string = ""
	for _, value := range items.Images {
		imgURL := fmt.Sprint(value)
		imgURL = strings.ReplaceAll(imgURL, "{", "")
		imgURL = strings.ReplaceAll(imgURL, "}", "")
		resp += `<img 
			src="` + imgURL + `" 
			class="img-fluid mr-2 mb-2" 
			alt="Responsive image"
			style="max-width: 250px; max-height: 250px"
		>`
	}
	return resp
}

var wg sync.WaitGroup

func main() {
	meta, err := spotify.GetMetadataSpotify()
	if err != nil {
		fmt.Println("Seems that you don't have the spotify app desktop installed  or is not open :(")
		log.Fatalf("failed getting metadata, err: %s", err.Error())
	}

	artistName := meta.ArtistName[0]

	if artistName == "" {
		// Not type music album
		// TODO SHOW OTHER VIEW
		panic("No album type")
	}

	wg.Add(3)
	albumChannel := make(chan *client.ResponseAlbum)
	go client.GetAlbumInfo(artistName, meta.AlbumName, &wg, albumChannel)
	albumInfo := <-albumChannel

	trackChannel := make(chan []string)
	go client.GetTrackInfo(artistName, meta.TrackName, &wg, trackChannel)
	trackInfo := <-trackChannel

	bandChannel := make(chan []string)
	go client.GetBandInfo(artistName, &wg, bandChannel)
	bandInfo := <-bandChannel
	wg.Wait()

	// Get images from album, band
	items := client.GetImagesBand(artistName, albumInfo.Album[0].ReleaseYear)

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

			<title>` + renderTitle(albumInfo) + `</title>
		</head>
		<body>
			<div class="container mt-5" >

				
				<h2>SPOTIFY - ALBUM BAND INFO</h2>

				<img id="loading" src="https://i.pinimg.com/originals/1c/13/f3/1c13f3fe7a6bba370007aea254e195e3.gif" width="50" height="50" style="display:none"/>
				<button type="button" class="btn btn-primary" id="button">REFRESH DATA</button>
				<br /> <br />
				
				<div id="wrapper-metadata-album">
				` + renderAlbumMetadata(albumInfo) + `
				</div>

				<div class="row" id="wrapper-album-images">
					` + renderImage(albumInfo.Album[0].ThumbFront) + renderImage(albumInfo.Album[0].ThumbBack) + renderImage(albumInfo.Album[0].ThumbCD) + `
				</div>

				<br>
				<p class="font-weight-bold">Description album: </p>
				<p class="text-justify" id="description-album">` + albumInfo.Album[0].Description + `</p>

				<br>
				<div id="review-album">
				` + renderReview(albumInfo.Album[0].Review) + `
				</div>
				<hr />
				<!-- TRACK INFO -->
				<br />
				<div id="track-info">
				` + renderTrackInfo(trackInfo) + `
				</div>
				<hr />

				<p class="font-weight-bold">BAND/ARTIST IMAGES: </p>

				<div id="wrapper-images-band" class="row">
				` + renderBandAlbumImages(*items) + `
				</div>
				<hr />

				<!-- BAND INFO -->
				<p class="font-weight-bold">BIO: </p>
				<p class="text-justify" id="artist-bio">` + bandInfo[0] + `</p>

				<p>FORM YEAR: <span id="artist-year">` + bandInfo[1] + ` ` + bandInfo[2] + `</span></p>

			</div>

			<script>
				var loading = document.getElementById('loading');

				document.getElementById('button').addEventListener('click', function(){
					loading.style.display = 'block';

					refresh().then( (data) => { 
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

						var metadataAlbum = document.getElementById('wrapper-metadata-album');
						metadataAlbum.innerHTML = data[9];

						document.title = data[10]

						var artistImages = document.getElementById('wrapper-images-band');
						artistImages.innerHTML = data[11];

						loading.style.display = 'none';
					})
				})
			</script>
		</body>
		</html>
		`
	ui, _ := lorca.New("data:text/html,"+url.PathEscape(htmlBody), "", 1100, 800)
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

		wg.Add(3)

		albumChannel := make(chan *client.ResponseAlbum)
		go client.GetAlbumInfo(artistName, meta.AlbumName, &wg, albumChannel)
		albumInfo := <-albumChannel

		trackChannel := make(chan []string)
		go client.GetTrackInfo(artistName, meta.TrackName, &wg, trackChannel)
		trackInfo := <-trackChannel

		bandChannel := make(chan []string)
		go client.GetBandInfo(artistName, &wg, bandChannel)
		bandInfo := <-bandChannel

		wg.Wait()

		albumMetadata := renderAlbumMetadata(albumInfo)

		var albumImageHTML string = ""
		albumImageHTML += renderImage(albumInfo.Album[0].ThumbFront)
		albumImageHTML += renderImage(albumInfo.Album[0].ThumbBack)
		albumImageHTML += renderImage(albumInfo.Album[0].ThumbCD)

		albumReview := renderReview(albumInfo.Album[0].Review)

		trackInfoHTML := renderTrackInfo(trackInfo)

		title := renderTitle(albumInfo)

		// Get images from album, band
		items := client.GetImagesBand(artistName, albumInfo.Album[0].ReleaseYear)
		bandImages := renderBandAlbumImages(*items)

		return []string{
			artistName,
			meta.AlbumName,
			albumImageHTML,
			albumInfo.Album[0].Description,
			albumReview,
			trackInfoHTML,
			bandInfo[0],
			bandInfo[1],
			bandInfo[2],
			albumMetadata,
			title,
			bandImages,
		}
	})
	// Call above `hello` function then log to the JS console
	ui.Eval("hello().then( (x) => { console.log(x) })")
	// Wait until UI window is closed
	<-ui.Done()
}
