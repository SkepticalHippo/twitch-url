package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"

	"github.com/grafov/m3u8"
)

type TokenResponse struct {
	Signature string `json:"sig"`
	Token     string `json:"token"`
}

const TOKEN_API string = "http://api.twitch.tv/api/channels/%s/access_token?client_id=%s"
const PLAYLIST_API string = "http://usher.twitch.tv/api/channel/hls/%s.m3u8?player=twitchweb&&token=%s&sig=%s&allow_audio_only=true&allow_source=true&type=any&p=%s"

func main() {
	channel := os.Args[1]
	client_id := os.Args[2]

	res, err := http.Get(fmt.Sprintf(TOKEN_API, channel, client_id))

	if err != nil {
		fmt.Println(err)
	}

	defer res.Body.Close()

	var tokenResponse TokenResponse

	json.NewDecoder(res.Body).Decode(&tokenResponse)

	resp, err := http.Get(fmt.Sprintf(
		PLAYLIST_API,
		channel,
		tokenResponse.Token,
		tokenResponse.Signature,
		strconv.Itoa(rand.Intn(999999)), //random int up to 6 digits
	))

	if err != nil {
		fmt.Println(err)
	}

	if err != nil {
		fmt.Println(err)
	}

	p, listType, err := m3u8.DecodeFrom(resp.Body, false)

	if listType != m3u8.MASTER {
		fmt.Println("Stream offline or does not exist.")
	}

	masterpl := p.(*m3u8.MasterPlaylist)
	for _, variant := range masterpl.Variants {
		fmt.Printf("[%s] %s\n\n", variant.VariantParams.Video, variant.URI)
	}
}
