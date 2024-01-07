package news

import (
	"context"
	"net/http"

	"github.com/charmbracelet/log"

	"github.com/alfredosa/go-youtube-reddit-automation/config"
	dbmod "github.com/alfredosa/go-youtube-reddit-automation/db"
	"github.com/barthr/newsapi"
	"github.com/jmoiron/sqlx"
)

func PullLatestNews(config config.Config, db *sqlx.DB) ([]newsapi.Article, error) {
	var ctx = context.Background()
	newsApiClient := newsapi.NewClient(config.NewsAPI.API_Key, newsapi.WithHTTPClient(http.DefaultClient), newsapi.WithUserAgent("go-youtube-reddit-automation"))

	articles, err := newsApiClient.GetEverything(ctx, &newsapi.EverythingParameters{
		Language: "en",
		Keywords: "tech",
		SortBy:   "popularity",
	})

	if err != nil {
		return nil, err
	}

	log.Info("Found", "posts", len(articles.Articles))

	posts := dbmod.FilterPostedPosts(articles.Articles, db)
	log.Info("Found after filtering", "posts", len(posts))
	processedPosts := CreateTTSAndSSFiles(posts, config)
	return processedPosts, nil
}
