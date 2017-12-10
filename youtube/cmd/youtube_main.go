package main

import (
	"fmt"

	"github.com/aracki/countgo/youtube"
	"github.com/aracki/countgo/youtube/service"
)

func main() {

	// init youtube service
	s, _ := youtube.InitYoutubeService()

	// getting IvannSerbia channel info
	info, _ := service.ChannelInfo(s, "IvannSerbia")
	fmt.Println(info)

	// getting all the lists
	lists, _ := service.AllPlaylists(s)
	// getting all the lists info concurrently
	service.PlaylistsInfo(s, lists)
	// getting all the videos of all playlists of mine
	service.AllVideos(s)
}