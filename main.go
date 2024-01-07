package main

// golang path github.com/alfredosa/go-youtube-reddit-automation

import (
	"os"

	"github.com/charmbracelet/log"

	"github.com/BurntSushi/toml"
	"github.com/alfredosa/go-youtube-reddit-automation/config"
	"github.com/alfredosa/go-youtube-reddit-automation/db"
	dbmod "github.com/alfredosa/go-youtube-reddit-automation/db"
	"github.com/alfredosa/go-youtube-reddit-automation/instagram"
	"github.com/alfredosa/go-youtube-reddit-automation/news"
	"github.com/alfredosa/go-youtube-reddit-automation/utils"
	"github.com/alfredosa/go-youtube-reddit-automation/video"
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

func main() {
	var config config.Config

	_, err := toml.DecodeFile("config.toml", &config)
	if err != nil {
		log.Fatal(err)
	}
	db := db.Connect(config)
	log.Info("Connected to DB", "driver", db.DriverName())
	CreateVideo(config, db)
}

func CreateVideo(config config.Config, db *sqlx.DB) {
	for {
		log.Info("Starting new iteration")
		posts, err := news.PullLatestNews(config, db)
		if len(posts) == 0 {
			log.Info("No new posts found, exiting")
			os.Exit(0)
		}
		if err != nil {
			log.Fatal(err)
		}

		if !utils.CheckFileExists("resultwsound", "studio/staging") {
			err := video.CreateVideo(posts, config)
			if err != nil {
				log.Fatal(err, "file: ", err.Error())
			}
			log.Info("Finished creating video, now adding posts to db")
			err = dbmod.InsertPostsFromReddit(posts, db)
			if err != nil {
				log.Fatal(err)
			}

		} else {
			log.Info("Final Video already exists, skipping video creation")
		}

		// keywords := youtube.GetKeywords()
		// title := youtube.GetVideoTitle()
		// description := youtube.GetVideoDescription(posts)
		// filename := "studio/staging/resultwsound.mp4"
		// youtube.YoutubeUpload(config, title, description, "25", "public", keywords, filename)
		instagram.UploadInstagramVideo(config, posts)
	}
}
