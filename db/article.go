package db

import (
	"database/sql"
	"rss-score/model"
)

func InsertOrUpdate(db *sql.DB, article *model.Article, score int) error {
	articleInDB, err := GetArticleScore(db, article.ArticleURL)
	if err == nil {
		_, err := db.Exec(
			`UPDATE rss_scores SET score = ? WHERE feed_url = ? AND article_url = ?`,
			score, articleInDB.FeedURL, articleInDB.ArticleURL,
		)
		return err
	}

	_, err = db.Exec(
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
		score,
	)

	return err
}

func GetArticleScore(db *sql.DB, articleURL string) (model.Article, error) {
	row := db.QueryRow(
		`SELECT
            feed_url,
            feed_title,
            article_url,
            article_title,
            pub_date
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
	)
	return as, err
}
