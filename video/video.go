package video

import (
	"log"
	"os/exec"
	"strconv"
	"sync"

	"github.com/alfredosa/go-youtube-reddit-automation/config"
	"github.com/alfredosa/go-youtube-reddit-automation/reddit"
	"github.com/alfredosa/go-youtube-reddit-automation/utils"
	rdt "github.com/vartanbeno/go-reddit/v2/reddit"
)

func CreateVideo(posts []*rdt.Post, config config.Config) {
	final_cut_duration, err := reddit.GetMP3Length("audio/final_cut.mp3")

	if err != nil {
		log.Fatal(err)
	}

	audios := utils.GetAudios()
	var wg sync.WaitGroup

	for _, audio := range audios {
		wg.Add(1)

		go func(audio string) {
			defer wg.Done()

			// Get the duration of the audio file
			duration, err := reddit.GetMP3Length(audio)
			if err != nil {
				log.Fatal(err)
			}

			// Create a video with the same duration as the audio
			video := CreateVideoWithLength(duration, config)

			log.Printf("Created video %s", video)
		}(audio)
	}

	wg.Wait()

	log.Printf("Finished creating videos, now concatenating them and adding audio")
	CreateFinalCutVideo()
}

func CreateVideoWithLength(duration int, config config.Config) string {
	// ffmpeg -ss 160 -i video/gta.mp4 -t 126 -c copy sample/output.mp4
	// execute command
	cmd := exec.Command("ffmpeg", "-ss", "160", "-i", "studio/gta.mp4", "-t", strconv.Itoa(duration), "-c", "copy", "studio/background.mp4")

}
