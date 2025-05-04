package service

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"rss-score/api"
	"rss-score/db"
	"rss-score/model"
)

type Service interface {
	Run(videoID string, score int) error
}

type service struct {
	api   api.API
	store db.ArticleRepository
}

func New(api api.API, store db.ArticleRepository) Service {
	return &service{api: api, store: store}
}

func (s *service) Run(videoID string, score int) error {
	// Fetch metadata
	video, err := s.api.Video(videoID)
	if err != nil {
		return err
	}

	pubDate, err := time.Parse(time.RFC3339Nano, video.PubDate)
	if err != nil {
		return err
	}

	// channel, err := api.Channel(video.ChannelID)
	// if err != nil {
	// 	return nil, err
	// }

	article := &model.Article{
		ArticleURL:   fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID),
		FeedURL:      fmt.Sprintf("https://www.youtube.com/feeds/videos.xml?channel_id=%s", video.ChannelID),
		ArticleTitle: video.VideoTitle,
		FeedTitle:    video.ChannelTitle,
		PubDate:      int(pubDate.Unix()),
		Score:        score,
	}

	// store
	// TODO: add transaction
	articleInDB, err := s.store.GetByArticleURL(article.ArticleURL)
	switch {
	case err == nil:
		return s.store.UpdateScore(articleInDB, article.Score)
	case errors.Is(err, sql.ErrNoRows):
		return s.store.Insert(article)
	default:
		return err
	}
}
