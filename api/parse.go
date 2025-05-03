package api

import (
	"fmt"
	"time"

	"rss-score/model"
)

func (api *googleAPI) FetchMetadata(videoID string) (*model.Article, error) {
	video, err := api.video(videoID)
	if err != nil {
		return nil, err
	}

	pubDate, err := time.Parse(time.RFC3339Nano, video.PubDate)
	if err != nil {
		return nil, err
	}

	// channel, err := api.channel(video.ChannelID)
	// if err != nil {
	// 	return nil, err
	// }

	return &model.Article{
		ArticleURL:   fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID),
		FeedURL:      fmt.Sprintf("https://www.youtube.com/feeds/videos.xml?channel_id=%s", video.ChannelID),
		ArticleTitle: video.VideoTitle,
		FeedTitle:    video.ChannelTitle,
		PubDate:      int(pubDate.Unix()),
	}, nil
}
