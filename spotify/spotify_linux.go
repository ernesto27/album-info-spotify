// +build linux

package spotify

import (
	"log"
	"reflect"

	"github.com/godbus/dbus"
)

const (
	sender          = "org.mpris.MediaPlayer2.spotify"
	path            = "/org/mpris/MediaPlayer2"
	member          = "org.mpris.MediaPlayer2.Player"
	metadataMessage = member + ".Metadata"
)

// Metadata contains Spotify player metadata
type SpotifyMetadata struct {
	ArtistName  []string `spotify:"xesam:artist"`
	TrackName   string   `spotify:"xesam:title"`
	AlbumName   string   `spotify:"xesam:album"`
	AlbumArtist []string `spotify:"xesam:albumArtist"`
	AutoRating  float64  `spotify:"xesam:autoRating"`
	DiskNumber  int32    `spotify:"xesam:discNumber"`
	TrackNumber int32    `spotify:"xesam:trackNumber"`
	URL         string   `spotify:"xesam:url"`
	TrackID     string   `spotify:"mpris:trackid"`
	Length      uint64   `spotify:"mpris:length"`
}

// parseMetadata returns a parsed Metadata struct
func parseMetadata(variant dbus.Variant) *SpotifyMetadata {
	metadataMap := variant.Value().(map[string]dbus.Variant)
	metadataStruct := new(SpotifyMetadata)

	valueOf := reflect.ValueOf(metadataStruct).Elem()
	typeOf := reflect.TypeOf(metadataStruct).Elem()

	for key, val := range metadataMap {
		for i := 0; i < typeOf.NumField(); i++ {
			field := typeOf.Field(i)
			if field.Tag.Get("spotify") == key {
				field := valueOf.Field(i)
				field.Set(reflect.ValueOf(val.Value()))
			}
		}
	}

	return metadataStruct
}

// GetMetadata returns the current metadata from the Spotify app
func GetMetadataSpotify() (*SpotifyMetadata, error) {
	conn := getConn()
	obj := conn.Object(sender, path)
	property, err := obj.GetProperty(metadataMessage)
	if err != nil {
		return nil, err
	}

	return parseMetadata(property), nil
}

func getConn() *dbus.Conn {
	conn, err := dbus.SessionBus()
	if err != nil {
		log.Fatal(err)
	}
	return conn
}
