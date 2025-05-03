package api

import (
	"encoding/json"
	"fmt"
)

type (
	videoDetails struct {
		VideoTitle   string `json:"title"`
		ChannelTitle string `json:"channelTitle"`
		ChannelID    string `json:"channelId"`
		PubDate      string `json:"publishedAt"`
	}

	videoItem struct {
		ID      string       `json:"id"`
		Details videoDetails `json:"snippet"`
	}

	videoResponse struct {
		Videos []videoItem `json:"items"`
	}
)

func (api *googleAPI) video(videoID string) (*videoDetails, error) {
	body, err := api.get("videos", videoID)
	if err != nil {
		return nil, err
	}

	response := videoResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	if len(response.Videos) != 1 {
		return nil, fmt.Errorf("expected 1 video, but got %d", len(response.Videos))
	}

	return &response.Videos[0].Details, nil
}
