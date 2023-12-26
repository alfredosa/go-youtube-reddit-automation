package video

import (
	"log"
	"math/rand"
	"os/exec"
	"strconv"
	"sync"

	"github.com/alfredosa/go-youtube-reddit-automation/config"
	"github.com/alfredosa/go-youtube-reddit-automation/reddit"
	rdt "github.com/vartanbeno/go-reddit/v2/reddit"
)

func CreateVideo(posts []*rdt.Post, config config.Config) {
	final_cut_duration, err := reddit.GetMP3Length("audio/result/final_cut.mp3")

	log.Printf("Final cut duration: %d", final_cut_duration)
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup

	for _, post := range posts {
		wg.Add(1)

		go func(post *rdt.Post) {
			defer wg.Done()
			audio := "audio/" + post.ID + ".mp3"
			// Get the duration of the audio file
			duration, err := reddit.GetMP3Length(audio)
			if err != nil {
				log.Fatal(err)
			}

			// Create a video with the same duration as the audio
			CreateVideoWithLength(duration, post.ID)
		}(post)
	}

	wg.Wait()

	log.Printf("Finished creating videos, now concatenating them and adding audio")
}

func CreateVideoWithLength(duration int, id string) {
	// ffmpeg -ss 160 -i studio/gta4_hd.mp4 -t 126 -vf "scale=-1:1920,crop=1080:1920:(iw-1080)/2:0" studio/staging/output.mp4
	// execute command
	actualDruration := duration + 1 // space between videos

	// get random start from 0 to 1 ho
	randomNumber := rand.Intn(3600)

	stuidoPath := "studio/staging/" + id + ".mp4"
	log.Printf("Creating video %s", stuidoPath)
	cmd := exec.Command("ffmpeg", "-ss", strconv.Itoa(randomNumber), "-i", "studio/gta4_hd.mp4", "-t", strconv.Itoa(actualDruration), "-vf", "scale=-1:1920,crop=1080:1920:(iw-1080)/2:0", stuidoPath)
	err := cmd.Run()

	if err != nil {
		log.Fatal(err)
	}

	// ffmpeg -i studio/staging/18qimw8.mp4 -i screenshots/18qimw8_0.jpg -filter_complex "[1:v]scale=640:-1,format=rgba,colorchannelmixer=aa=0.9[img1];[0:v][img1]overlay=(main_w-overlay_w)/2:(main_h-overlay_h)/4" output1.mp4
	// ffmpeg -i output1.mp4 -i screenshots/18qimw8_1.jpg -filter_complex "[1:v]scale=640:-1,format=rgba,colorchannelmixer=aa=0.9[img2];[0:v][img2]overlay=(main_w-overlay_w)/2:(main_h-overlay_h)/2" output.mp4

}
