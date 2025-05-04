package api

import (
	"encoding/json"
	"fmt"
)

type (
	ChannelDetails struct {
		CustomURL string `json:"customUrl"`
		Title     string `json:"title"`
	}

	channelItem struct {
		ID      string         `json:"id"`
		Details ChannelDetails `json:"snippet"`
	}

	channelResponse struct {
		Channels []channelItem `json:"items"`
	}
)

func (api *youTubeAPI) Channel(channelID string) (*ChannelDetails, error) {
	body, err := api.get("channels", channelID)
	if err != nil {
		return nil, err
	}

	response := channelResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	if len(response.Channels) != 1 {
		return nil, fmt.Errorf("expected 1 channel, but got %d", len(response.Channels))
	}

	return &response.Channels[0].Details, nil
}
