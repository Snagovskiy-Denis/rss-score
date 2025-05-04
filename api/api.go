package api

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

type GoogleAPI interface {
	Video(ID string) (*VideoDetails, error)
	Channel(ID string) (*ChannelDetails, error)
}

type googleAPI struct {
	apiKey  string
	baseURL string
}

func New(apiKey []byte) GoogleAPI {
	return &googleAPI{
		baseURL: "https://www.googleapis.com/youtube/v3/%s?id=%s&key=%s&part=snippet",
		apiKey:  strings.TrimSpace(string(apiKey)),
	}
}

func (api *googleAPI) get(method string, id string) ([]byte, error) {
	URL := fmt.Sprintf(api.baseURL, method, id, api.apiKey)

	response, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("%s from %s", response.Status, URL)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
