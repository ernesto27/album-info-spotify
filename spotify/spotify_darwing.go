// +build darwin

package spotify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
)

// Metadata contains Spotify player metadata
type SpotifyMetadata struct {
	ArtistName []string `json:"artist"`
	AlbumName  string   `json:"album"`
	TrackName  string   `json:"name"`
}

// GetMetadata returns the current metadata from the Spotify app
func GetMetadataSpotify() (SpotifyMetadata, error) {
	command := `
		on escape_quotes(string_to_escape)
			set AppleScript's text item delimiters to the "\""
			set the item_list to every text item of string_to_escape
			set AppleScript's text item delimiters to the "\\\""
			set string_to_escape to the item_list as string
			set AppleScript's text item delimiters to ""
			return string_to_escape
		end escape_quotes

		tell application "Spotify"
			set ctrack to "{"
			set ctrack to ctrack & "\"artist\": \"" & my escape_quotes(current track's artist) & "\""
			set ctrack to ctrack & ",\"album\": \"" & my escape_quotes(current track's album) & "\""
			set ctrack to ctrack & ",\"disc_number\": " & current track's disc number
			set ctrack to ctrack & ",\"duration\": " & current track's duration
			set ctrack to ctrack & ",\"played_count\": " & current track's played count
			set ctrack to ctrack & ",\"track_number\": " & current track's track number
			set ctrack to ctrack & ",\"popularity\": " & current track's popularity
			set ctrack to ctrack & ",\"id\": \"" & current track's id & "\""
			set ctrack to ctrack & ",\"name\": \"" & my escape_quotes(current track's name) & "\""
			set ctrack to ctrack & ",\"album_artist\": \"" & my escape_quotes(current track's album artist) & "\""
			set ctrack to ctrack & ",\"artwork_url\": \"" & current track's artwork url & "\""
			set ctrack to ctrack & ",\"spotify_url\": \"" & current track's spotify url & "\""
			set ctrack to ctrack & "}"
		end tell
		`

	cmd := exec.Command("osascript", "-e", command)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Run()
	info := out.String()

	r := regexp.MustCompile("\"artist\": \"([a-zA-Z 0-9 '/]+)\"")
	res := r.ReplaceAllString(info, "\"artist\": [\"$1\"]")

	data := SpotifyMetadata{}
	json.Unmarshal([]byte(res), &data)
	return data, nil

}

func Darwin() {
	fmt.Println("Im from mac")
}
