package main

import (
	"album-info-spotify/client"
	"album-info-spotify/spotify"
	"fmt"
	"log"
	"os"

	"github.com/godbus/dbus"
)

var nameBand string
var albumBand string
var trackName string

func main() {
	if len(os.Args) < 3 {
		conn := getConn()
		var meta *spotify.SpotifyMetadata
		meta, err := spotify.GetMetadataSpotify(conn)
		if err != nil {
			fmt.Println("Seems that you don't have the spotify app desktop installed  or is not open :(")
			log.Fatalf("failed getting metadata, err: %s", err.Error())
		}

		trackName = meta.Title
		nameBand = meta.Artist[0]
		albumBand = meta.Album
	} else {
		// Search by band , album enter by user
		nameBand = os.Args[1]
		albumBand = os.Args[2]
		trackName = ""
	}

	albumInfo, err := client.GetAlbumInfo(nameBand, albumBand)
	if err != nil {
		panic(err)
	}

	fmt.Println()
	fmt.Println(albumInfo[5] + " - " + albumInfo[6])
	fmt.Println(albumInfo[7] + " - " + albumInfo[8])
	if albumInfo[9] != "" {
		fmt.Println("Label: " + albumInfo[9])
	}
	fmt.Println()
	fmt.Println("----------------------------------------------------------------------------------------------------------------------------------------")
	fmt.Println()

	fmt.Println("ALBUM DESCRIPTION:")
	fmt.Println(string(albumInfo[3]))

	if albumInfo[4] != "" {
		fmt.Println()
		fmt.Println("----------------------------------------------------------------------------------------------------------------------------------------")
		fmt.Println()
		fmt.Println("REVIEW:")
		fmt.Println(string(albumInfo[4]))
	}

	fmt.Println()
	fmt.Println("----------------------------------------------------------------------------------------------------------------------------------------")
	fmt.Println()

	// BAND INFO
	bandInfo, err := client.GetBandInfo(nameBand)
	if err != nil {
		panic(err)
	}
	fmt.Println()
	fmt.Println("BAND INFO:")
	fmt.Println(bandInfo[0])
	fmt.Println()
	fmt.Println("YEAR: " + bandInfo[1])
	fmt.Println()
	fmt.Println("COUNTRY: " + bandInfo[2])
	fmt.Println()

	// TRACK DESCRIPTION
	trackInfo, err := client.GetTrackInfo(nameBand, trackName)
	if err != nil {
		panic(err)
	}

	if trackInfo[0] != "" {
		fmt.Println("----------------------------------------------------------------------------------------------------------------------------------------")
		fmt.Println()
		fmt.Println("TRACK DESCRIPTION:")
		fmt.Println(trackInfo[0])

		fmt.Println()

		fmt.Println("TRACK YOUTUBE URL:")
		fmt.Println(trackInfo[3])
	}
}

func getConn() *dbus.Conn {
	conn, err := dbus.SessionBus()
	if err != nil {
		log.Fatal(err)
	}
	return conn
}
