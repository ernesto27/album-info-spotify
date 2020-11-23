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
	return `<img src="` + imageURL + `" width="230" height="200" class="mr-6">`
}

func renderReview(text string) string {
	if text == "" {
		return ""
	}

	return `<p class="text-3xl">Review album</p>
			<p class="mt-2 text-justify">` + text + `</p>`
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

		resp += `<h3>Track info</h3>`

		if thumb != "" {
			resp += `<img class=" mb-3" width=150 height=150 src="` + thumb + `" />`
		}

		resp += `<p class="text-justify">` + text + `</p>`

		if videoURL != "" {
			resp += `
			<div class="flex-col space-x-4 mt-10">
				<iframe width="560" height="315" src="https://www.youtube.com/embed/` + videoURL + `" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>	
			</div>
			`
		}
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

	resp += `<h1 class="mt-0 mb-2 text-white text-4xl">` + albumInfo.Album[0].Artist +
		` -  ` + albumInfo.Album[0].Name + ` -  ` + albumInfo.Album[0].ReleaseYear + `</h1>`

	if albumInfo.Album[0].Label != "" {
		resp += `<p class="text-gray-600 mb-2 text-sm">Label: ` + albumInfo.Album[0].Label + `</p>`
	}

	if albumInfo.Album[0].Style != "" {
		resp += `<p class="text-gray-600 text-sm">Genre: ` + albumInfo.Album[0].Style + `</p>`
	}

	if albumInfo.Album[0].Score != "" {
		resp += `<p class="text-gray-600 text-sm">Score: ` + albumInfo.Album[0].Score + `</p>`
	}

	return resp
}

func renderBandAlbumImages(items client.Items) string {
	var resp string = ""
	for _, value := range items.Images {
		imgURL := fmt.Sprint(value)
		imgURL = strings.ReplaceAll(imgURL, "{", "")
		imgURL = strings.ReplaceAll(imgURL, "}", "")
		resp += `<img class="mr-3 w-16 md:w-32 lg:w-48""
					style="max-height:200px"
					 src="` + imgURL + `"  />`
	}
	return resp
}

func renderBio(bandInfo *client.ResponseBand) string {
	var resp string = ""
	resp += `<p class="text-3xl">BIO</p>
				<p class="mt-2 text-justify">
					` + bandInfo.Band[0].Biograhpy + `
				</p>

				<p class="mt-3">Form year: ` + bandInfo.Band[0].FormedYear + `</p>`
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

	bandChannel := make(chan *client.ResponseBand)
	go client.GetBandInfo(artistName, &wg, bandChannel)
	bandInfo := <-bandChannel

	wg.Wait()

	// Get images from album, band
	items := client.GetImagesBand(artistName, albumInfo.Album[0].ReleaseYear)

	fmt.Print(trackInfo)

	// Create UI with data URI
	/*
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
	*/
	htmlBody := `
	<html>
	<title>` + renderTitle(albumInfo) + `</title>
	<head>
		<link href="https://unpkg.com/tailwindcss@^2/dist/tailwind.min.css" rel="stylesheet">

	</head>
	<body>
		<div class="bg-black text-gray-300 min-h-screen p-10">
	
			<!-- header -->
			<div class="flex">
				` + renderImage(albumInfo.Album[0].ThumbFront) + renderImage(albumInfo.Album[0].ThumbBack) + renderImage(albumInfo.Album[0].ThumbCD) + `
			

				<div class="flex flex-col justify-center">
					<!-- content -->
					<h4 class="mt-0 mb-2 uppercase text-gray-500 tracking-widest text-xs">Album metadata spotify</h4>
					` + renderAlbumMetadata(albumInfo) + `
				</div>

				<!-- <div class="flex flex-row">
					
				</div> -->
			</div>
			
			<!-- action buttons -->
			<div class="mt-6 flex justify-between">
			<div class="flex">
				<button id="refresh-data" class="mr-2 bg-green-500 text-green-100 block py-2 px-8 rounded-full">Refresh data</button>

			</div>
			<!-- <div class="text-gray-600 text-sm tracking-widest text-right">
				<h5 class="mb-1">Followers</h5>
				<p>5,055</p>
			</div> -->
			</div>

			<div class="container mt-10 ">
				<p class="text-3xl">Album description</p>

				<p class="mt-2 text-justify">
					` + albumInfo.Album[0].Description + `
				</p>
			</div>

			<div class="container mt-10 ">
			` + renderReview(albumInfo.Album[0].Review) + `
			</div>

			<div class="container mt-10 ">
			` + renderTrackInfo(trackInfo) + `
			</div>

			<div class="container mt-10 ">
				<p class="text-3xl">Band/Artist images</p>

				<div class="flex overflow-x-scroll mt-4">
					` + renderBandAlbumImages(*items) + `
		
				</div>
			</div>

			<div class="container mt-10 ">
				` + renderBio(bandInfo) + `
			</div>
		</div>

		<script>
		document.getElementById('refresh-data').addEventListener('click', function(){


			refresh().then( (data) => { 
				console.log(data)

				document.title = data.title;

			})
		})
		</script>
	</body>
	</html>
	`

	ui, _ := lorca.New("data:text/html,"+url.PathEscape(htmlBody), "", 1100, 800)
	defer ui.Close()

	ui.Bind("refresh", func() map[string]string {
		// meta, err := spotify.GetMetadataSpotify()
		// if err != nil {
		// 	fmt.Println("Seems that you don't have the spotify app desktop installed  or is not open :(")
		// 	log.Fatalf("failed getting metadata, err: %s", err.Error())
		// }

		// artistName := meta.ArtistName[0]

		// wg.Add(3)

		// albumChannel := make(chan *client.ResponseAlbum)
		// go client.GetAlbumInfo(artistName, meta.AlbumName, &wg, albumChannel)
		// albumInfo := <-albumChannel

		// trackChannel := make(chan []string)
		// go client.GetTrackInfo(artistName, meta.TrackName, &wg, trackChannel)
		// trackInfo := <-trackChannel

		// bandChannel := make(chan []string)
		// go client.GetBandInfo(artistName, &wg, bandChannel)
		// bandInfo := <-bandChannel

		// wg.Wait()

		// albumMetadata := renderAlbumMetadata(albumInfo)

		// var albumImageHTML string = ""
		// albumImageHTML += renderImage(albumInfo.Album[0].ThumbFront)
		// albumImageHTML += renderImage(albumInfo.Album[0].ThumbBack)
		// albumImageHTML += renderImage(albumInfo.Album[0].ThumbCD)

		// albumReview := renderReview(albumInfo.Album[0].Review)

		// trackInfoHTML := renderTrackInfo(trackInfo)

		// title := renderTitle(albumInfo)

		// // Get images from album, band
		// items := client.GetImagesBand(artistName, albumInfo.Album[0].ReleaseYear)
		// bandImages := renderBandAlbumImages(*items)

		// return []string{
		// 	artistName,
		// 	meta.AlbumName,
		// 	albumImageHTML,
		// 	albumInfo.Album[0].Description,
		// 	albumReview,
		// 	trackInfoHTML,
		// 	bandInfo[0],
		// 	bandInfo[1],
		// 	bandInfo[2],
		// 	albumMetadata,
		// 	title,
		// 	bandImages,
		// }

		meta, err := spotify.GetMetadataSpotify()
		if err != nil {
			fmt.Println("Seems that you don't have the spotify app desktop installed  or is not open :(")
			log.Fatalf("failed getting metadata, err: %s", err.Error())
		}

		artistName := meta.ArtistName[0]

		wg.Add(1)

		albumChannel := make(chan *client.ResponseAlbum)
		go client.GetAlbumInfo(artistName, meta.AlbumName, &wg, albumChannel)
		albumInfo := <-albumChannel

		// trackChannel := make(chan []string)
		// go client.GetTrackInfo(artistName, meta.TrackName, &wg, trackChannel)
		// trackInfo := <-trackChannel

		// bandChannel := make(chan []string)
		// go client.GetBandInfo(artistName, &wg, bandChannel)
		// bandInfo := <-bandChannel

		wg.Wait()

		title := renderTitle(albumInfo)
		n := map[string]string{
			"title": title,
			"bar":   "some bar",
		}
		return n
	})

	// Wait until UI window is closed
	<-ui.Done()
}
