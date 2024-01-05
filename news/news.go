package news

import (
	"context"
	"net/http"
	"os"

	"github.com/charmbracelet/log"

	"github.com/alfredosa/go-youtube-reddit-automation/config"
	dbmod "github.com/alfredosa/go-youtube-reddit-automation/db"
	"github.com/barthr/newsapi"
	"github.com/jmoiron/sqlx"
)

var ctx = context.Background()

func PullLatestNews(config config.Config, db *sqlx.DB) ([]newsapi.Article, error) {
	newsApiClient := newsapi.NewClient(config.NewsAPI.API_Key, newsapi.WithHTTPClient(http.DefaultClient))

	articles, err := newsApiClient.GetEverything(context.Background(), &newsapi.EverythingParameters{
		Language: "en",
		Keywords: "tech",
		SortBy:   "popularity",
	})

	if err != nil {
		return nil, err
	}

	// TODO! FILTER OUT POSTS THAT ARE ALREADY POSTED. DB needs to be setup first
	log.Info("Found", "posts", len(articles.Articles))
	os.Exit(0)
	posts = dbmod.FilterPostedPosts(posts, db)
	log.Info("Found after filtering", "posts", len(posts))
	processedPosts := CreateTTSAndSSFiles(posts, config)
	return processedPosts, nil
}
