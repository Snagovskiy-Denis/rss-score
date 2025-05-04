package db

import "rss-score/model"

type ArticleRepository interface {
	GetByArticleURL(articleURL string) (*model.Article, error)
	UpdateScore(articleInDB *model.Article, score int) error
	Insert(article *model.Article) error
}

func (store Store) GetByArticleURL(articleURL string) (*model.Article, error) {
	row := store.db.QueryRow(
		`SELECT
            feed_url,
            feed_title,
            article_url,
            article_title,
            pub_date,
			score
        FROM rss_scores
        WHERE article_url = ?`,
		articleURL,
	)

	var as model.Article
	err := row.Scan(
		&as.FeedURL,
		&as.FeedTitle,
		&as.ArticleURL,
		&as.ArticleTitle,
		&as.PubDate,
		&as.Score,
	)

	return &as, err
}

func (store Store) UpdateScore(articleInDB *model.Article, score int) error {
	_, err := store.db.Exec(
		`UPDATE rss_scores
		SET score = ?
		WHERE feed_url = ? AND article_url = ?`,
		score, articleInDB.FeedURL, articleInDB.ArticleURL,
	)

	return err
}

func (store Store) Insert(article *model.Article) error {
	_, err := store.db.Exec(
		`INSERT INTO rss_scores (
	        feed_url,
	        feed_title,
	        article_url,
	        article_title,
	        pub_date,
	        score
	    )
	    VALUES (?, ?, ?, ?, ?, ?)`,
		article.FeedURL,
		article.FeedTitle,
		article.ArticleURL,
		article.ArticleTitle,
		article.PubDate,
		article.Score,
	)

	return err
}
