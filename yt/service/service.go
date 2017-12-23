package service

import (
	"fmt"
	"sync"

	"github.com/aracki/countgo/models"
	"google.golang.org/api/youtube/v3"
)

const (
	snippetContentDetailsStatistics = "snippet,contentDetails,statistics"
	snippetContentDetails           = "snippet,contentDetails"
)

// The getChannelInfo uses forUsername
// to get info (id, tittle, totalViews and description)
func ChannelInfo(service *youtube.Service, forUsername string) (string, error) {

	call := service.Channels.List(snippetContentDetailsStatistics)
	call = call.ForUsername(forUsername)
	response, err := call.Do()
	if err != nil {
		return "", err
	}

	var info string

	info = fmt.Sprintf("This channel's ID is %s. Its title is '%s', "+
		"and it has %d views. \n",
		response.Items[0].Id,
		response.Items[0].Snippet.Title,
		response.Items[0].Statistics.ViewCount)
	info += fmt.Sprintf(response.Items[0].Snippet.Description)

	return info, nil
}

// Gets all playlists of current user - maxResult is set to 50 (default is 5)
// returns array of all playlists (id, name, count)
func Playlists(service *youtube.Service) ([]models.Playlist, error) {

	// get all playlists
	call := service.Playlists.List(snippetContentDetails)
	call = call.MaxResults(50).Mine(true)
	response, err := call.Do()
	if err != nil {
		return nil, err
	}

	var pls []models.Playlist
	for _, pl := range response.Items {
		pls = append(pls, models.Playlist{
			Id:          pl.Id,
			Title:       pl.Snippet.Title,
			VideosCount: int(pl.ContentDetails.ItemCount),
		})
	}

	return pls, nil
}

// Gets all the videos of all playlists of mine
// Different goroutines are appending the same vds slice;
// WaitGroup waits for all goroutines to finish
func Videos(service *youtube.Service) (vds []models.Video, err error) {

	// get all playlists of mine
	call := service.Playlists.List(snippetContentDetails)
	call = call.MaxResults(50).Mine(true)
	response, err := call.Do()
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	wg.Add(len(response.Items))

	for _, pl := range response.Items {
		go func(p *youtube.Playlist) {
			v, _ := AllVideosByPlaylist(service, p)
			vds = append(vds, v...)
			wg.Done()
		}(pl)
	}
	wg.Wait()

	return vds, nil
}

// Gets all the videos of specific youtube.Playlist
func AllVideosByPlaylist(service *youtube.Service, pl *youtube.Playlist) ([]models.Video, error) {

	var vds []models.Video
	pageToken := ""

	for {
		call := service.PlaylistItems.List(snippetContentDetails)
		call = call.PlaylistId(pl.Id).MaxResults(50)
		response, err := call.PageToken(pageToken).Do()
		if err != nil {
			return nil, err
		}

		// move pageToken to another page
		pageToken = response.NextPageToken

		for _, item := range response.Items {
			t := ""
			if item.Snippet.Thumbnails != nil && item.Snippet.Thumbnails.Medium != nil {
				t = item.Snippet.Thumbnails.Medium.Url
			}
			vds = append(vds, models.Video{
				Title:       item.Snippet.Title,
				Description: item.Snippet.Description,
				PublishedAt: item.Snippet.PublishedAt,
				ResourceId:  item.Snippet.ResourceId.VideoId,
				Thumbnail:   t,
			})
		}
		// if there are no pages
		if pageToken == "" {
			break
		}
	}
	return vds, nil
}