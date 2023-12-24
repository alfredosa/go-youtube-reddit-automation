package main

// golang path github.com/alfredosa/go-youtube-reddit-automation

import (
	"log"

	"github.com/BurntSushi/toml"
	"github.com/alfredosa/go-youtube-reddit-automation/config"
	"github.com/alfredosa/go-youtube-reddit-automation/reddit"
)

func main() {
	var config config.Config

	_, err := toml.DecodeFile("config.toml", &config)
	if err != nil {
		log.Fatal(err)
	}

	val, err := reddit.PullLatestNews(config)
	if err != nil {
		log.Fatal(err)
	}
	println(val)

}
