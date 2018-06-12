package main

import (
	"os"
	"net/http"
	"time"
	"fmt"
	"log"
	"encoding/json"
	"net/url"
	"io/ioutil"
	"errors"
	"sync"
)

type timeRangeDescriptor string

var (
	baseUrl    = "https://api.spotify.com/v1/me/top"
	artistsURL = baseUrl + "/artists"
	tracksURL  = baseUrl + "/tracks"

	longTerm  timeRangeDescriptor = "long_term"
	shortTerm timeRangeDescriptor = "short_term"

	limit      = "5"
	spotifyKey = os.Getenv("SPOTIFY_KEY")
)

func main() {

	client := http.Client{Timeout: 2 * time.Second}

	handleArtists(client)
	fmt.Printf("\n")
	handleTracks(client)
}

func handleArtists(client http.Client) {
	wg := sync.WaitGroup{}

	wg.Add(2)

	var (
		artistsAllTime   *spotifyArtistResponse
		artistsShortTime *spotifyArtistResponse
		requestErr       error
	)

	go func() {
		body, err := getArtistsAllTime(client)
		if err != nil {
			requestErr = err
			wg.Done()
			return
		}
		artistsAllTime = body
		wg.Done()
	}()

	go func() {
		body, err := getArtistsShortTime(client)
		if err != nil {
			requestErr = err
			wg.Done()
			return
		}
		artistsShortTime = body
		wg.Done()
	}()

	wg.Wait()

	if requestErr != nil {
		log.Fatal(requestErr)
	}

	fmt.Println("Top five artists all time:")
	for i, artist := range artistsAllTime.Items {
		fmt.Printf("%d. %s\n", i+1, artist.Name)
	}

	fmt.Println("\nTop five artists last four weeks:")
	for i, artist := range artistsShortTime.Items {
		fmt.Printf("%d. %s\n", i+1, artist.Name)
	}
}

func handleTracks(client http.Client) {
	wg := sync.WaitGroup{}

	wg.Add(2)

	var (
		tracksAllTime   *spotifyTrackResponse
		tracksShortTime *spotifyTrackResponse
		requestErr      error
	)

	go func() {
		body, err := getTracksAllTime(client)
		if err != nil {
			requestErr = err
			wg.Done()
			return
		}
		tracksAllTime = body
		wg.Done()
	}()

	go func() {
		body, err := getTracksShortTime(client)
		if err != nil {
			requestErr = err
			wg.Done()
			return
		}
		tracksShortTime = body
		wg.Done()
	}()

	wg.Wait()

	if requestErr != nil {
		log.Fatal(requestErr)
	}

	fmt.Println("Top five trackks all time:")
	for i, track := range tracksAllTime.Items {
		fmt.Printf("%d. %s - %s\n", i+1, track.Name, track.Artists[0].Name)
	}

	fmt.Println("\nTop five tracks last four weeks:")
	for i, track := range tracksShortTime.Items {
		fmt.Printf("%d. %s - %s\n", i+1, track.Name, track.Artists[0].Name)
	}
}

func getArtistsAllTime(client http.Client) (*spotifyArtistResponse, error) {
	resp, err := artistRequest(longTerm, client)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func getArtistsShortTime(client http.Client) (*spotifyArtistResponse, error) {
	resp, err := artistRequest(shortTerm, client)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func getTracksAllTime(client http.Client) (*spotifyTrackResponse, error) {
	resp, err := trackRequest(longTerm, client)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func getTracksShortTime(client http.Client) (*spotifyTrackResponse, error) {
	resp, err := trackRequest(shortTerm, client)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func artistRequest(timeRange timeRangeDescriptor, client http.Client) (*spotifyArtistResponse, error) {
	form := url.Values{}
	form.Add("time_range", string(timeRange))
	form.Add("limit", string(limit))

	fullUrl := fmt.Sprintf("%s?%s", artistsURL, form.Encode())

	resp, err := spotifyRequest(fullUrl, client)

	if err != nil {
		return nil, err
	}

	body := &spotifyArtistResponse{}

	bytes, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(bytes, body); err != nil {
		return nil, err
	}
	return body, nil
}

func trackRequest(timeRange timeRangeDescriptor, client http.Client) (*spotifyTrackResponse, error) {
	form := url.Values{}
	form.Add("time_range", string(timeRange))
	form.Add("limit", string(limit))

	fullUrl := fmt.Sprintf("%s?%s", tracksURL, form.Encode())

	resp, err := spotifyRequest(fullUrl, client)

	if err != nil {
		return nil, err
	}

	body := &spotifyTrackResponse{}

	bytes, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(bytes, body); err != nil {
		return nil, err
	}
	return body, nil
}

func spotifyRequest(url string, client http.Client) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", spotifyKey))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if 200 != resp.StatusCode {
		return nil, errors.New(fmt.Sprintf("resp had non 200 status code: %s", resp.Status))
	}

	return resp, nil
}

type spotifyArtistResponse struct {
	Items []struct {
		ExternalUrls struct {
			Spotify string `json:"spotify"`
		} `json:"external_urls"`
		Followers struct {
			Href  interface{} `json:"href"`
			Total int         `json:"total"`
		} `json:"followers"`
		Genres []string `json:"genres"`
		Href   string   `json:"href"`
		ID     string   `json:"id"`
		Images []struct {
			Height int    `json:"height"`
			URL    string `json:"url"`
			Width  int    `json:"width"`
		} `json:"images"`
		Name       string `json:"name"`
		Popularity int    `json:"popularity"`
		Type       string `json:"type"`
		URI        string `json:"uri"`
	} `json:"items"`
	Total    int         `json:"total"`
	Limit    int         `json:"limit"`
	Offset   int         `json:"offset"`
	Href     string      `json:"href"`
	Previous interface{} `json:"previous"`
	Next     string      `json:"next"`
}

type spotifyTrackResponse struct {
	Items []struct {
		Album struct {
			AlbumType string `json:"album_type"`
			Artists []struct {
				ExternalUrls struct {
					Spotify string `json:"spotify"`
				} `json:"external_urls"`
				Href string `json:"href"`
				ID   string `json:"id"`
				Name string `json:"name"`
				Type string `json:"type"`
				URI  string `json:"uri"`
			} `json:"artists"`
			ExternalUrls struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			Href string `json:"href"`
			ID   string `json:"id"`
			Images []struct {
				Height int    `json:"height"`
				URL    string `json:"url"`
				Width  int    `json:"width"`
			} `json:"images"`
			Name                 string `json:"name"`
			ReleaseDate          string `json:"release_date"`
			ReleaseDatePrecision string `json:"release_date_precision"`
			Type                 string `json:"type"`
			URI                  string `json:"uri"`
		} `json:"album"`
		Artists []struct {
			ExternalUrls struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			Href string `json:"href"`
			ID   string `json:"id"`
			Name string `json:"name"`
			Type string `json:"type"`
			URI  string `json:"uri"`
		} `json:"artists"`
		DiscNumber int  `json:"disc_number"`
		DurationMs int  `json:"duration_ms"`
		Explicit   bool `json:"explicit"`
		ExternalIds struct {
			Isrc string `json:"isrc"`
		} `json:"external_ids"`
		ExternalUrls struct {
			Spotify string `json:"spotify"`
		} `json:"external_urls"`
		Href        string `json:"href"`
		ID          string `json:"id"`
		IsLocal     bool   `json:"is_local"`
		IsPlayable  bool   `json:"is_playable"`
		Name        string `json:"name"`
		Popularity  int    `json:"popularity"`
		PreviewURL  string `json:"preview_url"`
		TrackNumber int    `json:"track_number"`
		Type        string `json:"type"`
		URI         string `json:"uri"`
		LinkedFrom struct {
			ExternalUrls struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			Href string `json:"href"`
			ID   string `json:"id"`
			Type string `json:"type"`
			URI  string `json:"uri"`
		} `json:"linked_from,omitempty"`
	} `json:"items"`
	Total    int         `json:"total"`
	Limit    int         `json:"limit"`
	Offset   int         `json:"offset"`
	Href     string      `json:"href"`
	Previous interface{} `json:"previous"`
	Next     string      `json:"next"`
}
