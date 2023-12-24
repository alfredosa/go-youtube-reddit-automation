package main

// golang path github.com/alfredosa/go-youtube-reddit-automation

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/alfredosa/go-youtube-reddit-automation/config"
	"github.com/alfredosa/go-youtube-reddit-automation/reddit"
	htgotts "github.com/hegedustibor/htgo-tts"
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
	log.Printf("Generating audio files with speech %s", config.TextToSpeechSetup.Voice_ID)
	speech := htgotts.Speech{Folder: "audio", Language: config.TextToSpeechSetup.Voice_ID}

	for _, post := range posts {
		speech, err := speech.CreateSpeechFile(post.Title, fmt.Sprintf("%s.mp3", post.Title))
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Created audio file %s", speech)
	}
	log.Printf("Finished generating audio files")
	log.Printf("Cleaning up audio files")
	cleanUp()

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
