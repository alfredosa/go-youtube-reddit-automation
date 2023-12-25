package main

// golang path github.com/alfredosa/go-youtube-reddit-automation

import (
	"log"
	"os"
	"path/filepath"

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

	posts, err := reddit.PullLatestNews(config)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Found %d posts", len(posts))

	log.Printf("Finished generating audio files")
	// log.Printf("Cleaning up audio files")
	// cleanUp()

}

func cleanUp() {
	dirname := "audio/"

	d, err := os.Open(dirname)
	if err != nil {
		log.Fatal(err)
	}
	defer d.Close()

	files, err := d.Readdir(-1)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if err := os.RemoveAll(filepath.Join(dirname, file.Name())); err != nil {
			log.Fatal(err)
		}
	}
}
