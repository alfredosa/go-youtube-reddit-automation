package reddit

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/alfredosa/go-youtube-reddit-automation/config"
	"github.com/hajimehoshi/go-mp3"
	htgotts "github.com/hegedustibor/htgo-tts"
	rdt "github.com/vartanbeno/go-reddit/v2/reddit"
)

func CreateTTSAndSSFiles(posts []*rdt.Post, config config.Config) {
	var wg sync.WaitGroup

	log.Printf("Generating audio files with speech %s", config.TextToSpeechSetup.Voice_ID)
	speech := htgotts.Speech{Folder: "audio", Language: config.TextToSpeechSetup.Voice_ID}

	for _, post := range posts {
		wg.Add(1)

		go func(post *rdt.Post) {
			defer wg.Done()
			speech, err := speech.CreateSpeechFile(post.Title, post.ID)

			if err != nil {
				log.Fatal(err)
			}

			TakeScreenShot(post.Title, post.ID, config)
			log.Printf("Created audio file %s", speech)
		}(post)

	}

	wg.Wait()
}

func getMP3Length(filename string) (int, error) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file: ", err)
		return 0, err
	}
	defer file.Close()

	mp3Decoder, err := mp3.NewDecoder(file)
	if err != nil {
		fmt.Println("Error creating decoder: ", err)
		return 0, err
	}
	// 4 bytes per sample
	samples := mp3Decoder.Length() / 4

	// Samples divided by sample rate gives length in seconds
	audioLength := samples / int64(mp3Decoder.SampleRate())
	log.Println("Length in seconds: ", audioLength)

	return int(audioLength), nil
}
