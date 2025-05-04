package service

import (
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
	api   api.GoogleAPI
	store db.ArticleRepository
}

func New(api api.GoogleAPI, store db.ArticleRepository) Service {
	return &service{api: api, store: store}
}

func (s *service) Run(videoID string, score int) error {
	// Update score if article is presented in the database
	articleURL := fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID)
	articleInDB, err := s.store.GetByArticleURL(articleURL)
	// TODO add transaction and with update lock
	if err == nil {
		return s.store.UpdateScore(articleInDB, score)
	}

	// Fetch metadata
	video, err := s.api.Video(videoID)
	if err != nil {
		return err
	}

	pubDate, err := time.Parse(time.RFC3339Nano, video.PubDate)
	if err != nil {
		return err
	}

	article := &model.Article{
		ArticleURL:   articleURL,
		FeedURL:      fmt.Sprintf("https://www.youtube.com/feeds/videos.xml?channel_id=%s", video.ChannelID),
		ArticleTitle: video.VideoTitle,
		FeedTitle:    video.ChannelTitle,
		PubDate:      int(pubDate.Unix()),
		Score:        score,
	}

	// store
	return s.store.Upsert(article)
}
