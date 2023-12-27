package main

// golang path github.com/alfredosa/go-youtube-reddit-automation

import (
	"log"

	"github.com/BurntSushi/toml"
	"github.com/alfredosa/go-youtube-reddit-automation/config"
	"github.com/alfredosa/go-youtube-reddit-automation/reddit"
	"github.com/alfredosa/go-youtube-reddit-automation/utils"
	"github.com/alfredosa/go-youtube-reddit-automation/video"
)

func main() {
	var config config.Config

	_, err := toml.DecodeFile("config.toml", &config)
	if err != nil {
		log.Fatal(err)
	}

	posts, err := reddit.PullLatestNews(config)
	if err != nil {
		log.Fatal(err)
	}

	if utils.CheckFileExists("final_cut", "audio/result") && !utils.CheckFileExists("resultwsound", "studio/staging") {
		video.CreateVideo(posts, config)
	} else {
		println("Final Video already exists, skipping video creation")
	}

}
