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

func renderTitle(albumInfo *client.ResponseAlbum) string {
	var resp string = ""
	resp += albumInfo.Album[0].Artist + ` - ` + albumInfo.Album[0].Name
	return resp
}

func renderHeader(albumInfo *client.ResponseAlbum) string {
	resp := renderImage(albumInfo.Album[0].ThumbFront) + renderImage(albumInfo.Album[0].ThumbBack) + renderImage(albumInfo.Album[0].ThumbCD) + `
	<div class="flex flex-col justify-center">
		<!-- content -->
		<h4 class="mt-0 mb-2 uppercase text-gray-500 tracking-widest text-xs">Album metadata spotify</h4>
		` + renderAlbumMetadata(albumInfo) + `
	</div>`
	return resp
}

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

func renderTrackInfo(trackInfo *client.ResponseTrack) string {
	if trackInfo.Track[0].Description == "" {
		return ""
	}
	var resp string = ""

	if trackInfo.Track[0].Description != "" {

		resp += `<p class="text-3xl mb-3">Track info</p>`

		if trackInfo.Track[0].Thumb != "" {
			resp += `<img class=" mb-3" width=150 height=150 src="` + trackInfo.Track[0].Thumb + `" />`
		}

		resp += `<p class="text-justify mt-3">` + trackInfo.Track[0].Description + `</p>`

		if trackInfo.Track[0].YoutubeURL != "" {
			idVideo := strings.Split(trackInfo.Track[0].YoutubeURL, "v=")[1]

			resp += `
			<div class="flex-col space-x-4 mt-10">
				<iframe width="560" height="315" src="https://www.youtube.com/embed/` + idVideo + `" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"></iframe>	
			</div>

			<p class="mt-2"><a class="text-blue-400" target="blank" href="` + trackInfo.Track[0].YoutubeURL + `">` + trackInfo.Track[0].YoutubeURL + `</a></p>
			`
		}
	}

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

	trackChannel := make(chan *client.ResponseTrack)
	go client.GetTrackInfo(artistName, meta.TrackName, &wg, trackChannel)
	trackInfo := <-trackChannel

	bandChannel := make(chan *client.ResponseBand)
	go client.GetBandInfo(artistName, &wg, bandChannel)
	bandInfo := <-bandChannel

	wg.Wait()

	// Get images from album, band
	items := client.GetImagesBand(artistName, albumInfo.Album[0].ReleaseYear)

	htmlBody := `
	<html>
	<title>` + renderTitle(albumInfo) + `</title>
	<head>
		<link href="https://unpkg.com/tailwindcss@^2/dist/tailwind.min.css" rel="stylesheet">

	</head>
	<body>
		<div class="bg-black text-gray-300 min-h-screen p-10">
	
			<!-- header -->
			<div id="wrapper-header" class="flex">
				` + renderHeader(albumInfo) + `
				
			</div>

			<!-- action buttons -->
			<div class="mt-6 flex justify-between">
				<div class="flex">
					<button id="refresh-data" class="mr-2 bg-green-500 text-green-100 block py-2 px-8 rounded-full">Refresh data</button>

				</div>
			</div>

			<div class="container mt-10 ">
				<p class="text-3xl">Album description</p>
				<p id="wrapper-album-description" class="mt-2 text-justify">
					` + albumInfo.Album[0].Description + `
				</p>
			</div>

			<div id="wrapper-album-review" class="container mt-10 ">
			` + renderReview(albumInfo.Album[0].Review) + `
			</div>

			<div id="wrapper-album-track" class="container mt-10 ">
			` + renderTrackInfo(trackInfo) + `
			</div>

			<div class="container mt-10 ">
				<p class="text-3xl">Band/Artist images</p>

				<div id="wrapper-artist-images" class="flex overflow-x-scroll mt-4">
					` + renderBandAlbumImages(*items) + `
		
				</div>
			</div>

			<div id="wrapper-artist-bio"  class="container mt-10 ">
				` + renderBio(bandInfo) + `
			</div>
		</div>

		<script>
		document.getElementById('refresh-data').addEventListener('click', function(){

			var wrapperHeader = document.getElementById('wrapper-header');
			var wrapperAlbumDescription = document.getElementById('wrapper-album-description');
			var wrapperAlbumReview = document.getElementById('wrapper-album-review');
			var wrapperAlbumTrack = document.getElementById('wrapper-album-track');
			var wrapperArtistBio = document.getElementById('wrapper-artist-bio');
			var wrapperArtistImages = document.getElementById('wrapper-artist-images');

			refresh().then( (data) => { 
				console.log(data)

				document.title = data.title;

				wrapperHeader.innerHTML = data.header;
				wrapperAlbumDescription.innerHTML = data.albumDescription;
				wrapperAlbumReview.innerHTML = data.review;
				wrapperAlbumTrack.innerHTML = data.track;
				wrapperArtistBio.innerHTML = data.bio;
				wrapperArtistImages.innerHTML = data.artistImages;

			})
		})
		</script>
	</body>
	</html>
	`

	ui, _ := lorca.New("data:text/html,"+url.PathEscape(htmlBody), "", 1100, 800)
	defer ui.Close()

	ui.Bind("refresh", func() map[string]string {
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

		trackChannel := make(chan *client.ResponseTrack)
		go client.GetTrackInfo(artistName, meta.TrackName, &wg, trackChannel)
		trackInfo := <-trackChannel

		// bandChannel := make(chan []string)
		go client.GetBandInfo(artistName, &wg, bandChannel)
		bandInfo := <-bandChannel

		wg.Wait()

		title := renderTitle(albumInfo)
		header := renderHeader(albumInfo)
		review := renderReview(albumInfo.Album[0].Review)
		track := renderTrackInfo(trackInfo)
		bio := renderBio(bandInfo)

		items := client.GetImagesBand(artistName, albumInfo.Album[0].ReleaseYear)
		artistImages := renderBandAlbumImages(*items)

		n := map[string]string{
			"title":            title,
			"header":           header,
			"albumDescription": albumInfo.Album[0].Description,
			"review":           review,
			"track":            track,
			"bio":              bio,
			"artistImages":     artistImages,
		}
		return n
	})

	// Wait until UI window is closed
	<-ui.Done()
}
