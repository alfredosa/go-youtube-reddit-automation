package tts

import (
	"fmt"
	"log"
	"sync"

	"github.com/alfredosa/go-youtube-reddit-automation/config"
	htgotts "github.com/hegedustibor/htgo-tts"
	rdt "github.com/vartanbeno/go-reddit/v2/reddit"
)

func CreateTTSFiles(posts []*rdt.Post, config config.Config) {
	var wg sync.WaitGroup

	log.Printf("Generating audio files with speech %s", config.TextToSpeechSetup.Voice_ID)
	speech := htgotts.Speech{Folder: "audio", Language: config.TextToSpeechSetup.Voice_ID}

	for _, post := range posts {
		wg.Add(1)

		go func(post *rdt.Post) {
			defer wg.Done()
			speech, err := speech.CreateSpeechFile(post.Title, fmt.Sprintf(post.Title))

			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Created audio file %s", speech)
		}(post)
	}

	wg.Wait()
}
