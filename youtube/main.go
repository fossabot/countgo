package main

import (
	"context"
	"fmt"
	"sync"

	"google.golang.org/api/youtube/v3"
)

type Videos []Video

type Video struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	PublishedAt string `json:"publishedAt"`
	ResourceId  string `json:"title"`
	Thumbnail   string `json:"title"`
}

type Playlist struct {
	Title       string `json:"title"`
	VideosCount string `json:"videos_count"`
}

const (
	snippetContentDetailsStatistics = "snippet,contentDetails,statistics"
	snippetContentDetails           = "snippet,contentDetails"
)

// The getChannelInfo uses forUsername
// to get info (id, tittle, totalViews and description)
func getChannelInfo(service *youtube.Service, part string, forUsername string) {
	call := service.Channels.List(part)
	call = call.ForUsername(forUsername)
	response, err := call.Do()
	handleError(err, "")
	fmt.Println(fmt.Sprintf("This channel's ID is %s. Its title is '%s', "+
		"and it has %d views. \n",
		response.Items[0].Id,
		response.Items[0].Snippet.Title,
		response.Items[0].Statistics.ViewCount))
	fmt.Println(response.Items[0].Snippet.Description, "\n")
}

// The getAllPlaylists uses current user
// maxResult is set to 50 (default is 5)
// returns all playlists
func getAllPlaylists(service *youtube.Service, part string) (playlists []*youtube.Playlist) {

	call := service.Playlists.List(part)
	call = call.MaxResults(50).Mine(true)
	response, err := call.Do()
	handleError(err, "")

	var lists []*youtube.Playlist
	for _, item := range response.Items {
		lists = append(lists, item)
	}
	return lists
}

// The getPlaylistsInfo runs go routines for each playlist
// and call appendPlaylistInfo which populates plInfo array.
// Different goroutines are appending the same slice,
// WaitGroup waits for all goroutines to finish
func getPlaylistsInfo(service *youtube.Service, part string, playlists []*youtube.Playlist) {

	var wg sync.WaitGroup
	wg.Add(len(playlists))

	var pls []Playlist
	for _, playlist := range playlists {
		go func(pl *youtube.Playlist) {
			appendPlaylistInfo(service, part, pl, &pls)
			wg.Done()
		}(playlist)
	}
	wg.Wait()

	fmt.Println(pls)
}

// Gets all the videos of specific youtube.Playlist
func getAllVideos(service *youtube.Service, part string, pl *youtube.Playlist) (videos Videos) {

	var vds Videos
	pageToken := ""

	for {
		call := service.PlaylistItems.List(part)
		call = call.PlaylistId(pl.Id).MaxResults(50)
		response, err := call.PageToken(pageToken).Do()
		handleError(err, "")

		// move pageToken to another page
		pageToken = response.NextPageToken

		for _, item := range response.Items {
			t := ""
			if item.Snippet.Thumbnails != nil && item.Snippet.Thumbnails.Medium != nil {
				t = item.Snippet.Thumbnails.Medium.Url
			}
			vds = append(vds, Video{
				item.Snippet.Title,
				item.Snippet.Description,
				item.Snippet.PublishedAt,
				item.Snippet.ResourceId.VideoId,
				t,
			})
		}
		// if there are no pages
		if pageToken == "" {
			break
		}
	}
	return vds
}

func main() {
	ctx := context.Background()

	config := readConfigFile()

	client := getClient(ctx, config)
	service, err := youtube.New(client)

	handleError(err, "Error creating YouTube client")

	// getting IvannSerbia channel info
	getChannelInfo(service, snippetContentDetailsStatistics, "IvannSerbia")

	// getting all the lists
	lists := getAllPlaylists(service, snippetContentDetails)
	// getting all the lists info concurrently
	getPlaylistsInfo(service, snippetContentDetails, lists)
}
