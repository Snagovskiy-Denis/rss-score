package api

import (
	"encoding/json"
	"fmt"
)

type (
	VideoDetails struct {
		VideoTitle   string `json:"title"`
		ChannelTitle string `json:"channelTitle"`
		ChannelID    string `json:"channelId"`
		PubDate      string `json:"publishedAt"`
	}

	videoItem struct {
		ID      string       `json:"id"`
		Details VideoDetails `json:"snippet"`
	}

	videoResponse struct {
		Videos []videoItem `json:"items"`
	}
)

func (api *googleAPI) Video(ID string) (*VideoDetails, error) {
	body, err := api.get("videos", ID)
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
