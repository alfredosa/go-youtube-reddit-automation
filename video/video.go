package video

import (
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"github.com/alfredosa/go-youtube-reddit-automation/config"
	"github.com/alfredosa/go-youtube-reddit-automation/reddit"
	"github.com/alfredosa/go-youtube-reddit-automation/utils"
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
	preSound := ConcatAllVideos()
	audio := "audio/result/final_cut.mp3"
	AddAudioToVideo(preSound, audio, "studio/staging/resultwsound.mp4")
}

func CreateVideoWithLength(duration int, id string) {
	// ffmpeg -ss 160 -i studio/gta4_hd.mp4 -t 126 -vf "scale=-1:1920,crop=1080:1920:(iw-1080)/2:0" studio/staging/output.mp4
	// execute command
	actualDruration := duration + 1 // space between videos

	// get random start from 100 to 3600
	randomNumber := rand.Intn(3600-100) + 100

	studioStaging := "studio/staging/"

	if utils.CheckFileExists(id, studioStaging) {
		log.Printf("Enhanced Video %s already exists, skipping", id)
		return
	}

	log.Printf("Creating video %s", id)

	stuidoPath := studioStaging + id + ".mp4"
	cmd := exec.Command("ffmpeg", "-ss", strconv.Itoa(randomNumber), "-i", "studio/gta4_hd.mp4", "-t", strconv.Itoa(actualDruration), "-vf", "scale=-1:1920,crop=1080:1920:(iw-1080)/2:0", stuidoPath)
	err := cmd.Run()

	if err != nil {
		log.Fatal(err)
	}
	videoPath := studioStaging + id + ".mp4"
	OutputPath := studioStaging + id + "_enhanced.mp4"

	images := GetAllImagesByID(id)

	if len(images) == 1 {
		log.Printf("Adding one image to video %s", OutputPath)
		AddOneImageToVideo(videoPath, images[0], id)
	}
	if len(images) == 2 {
		log.Printf("Adding two images to video %s", OutputPath)
		AddTwoImagesToVideo(videoPath, images[0], images[1], id)
	}
	log.Printf("Finished creating video %s", OutputPath)
}

func AddOneImageToVideo(videoPath string, imagePath string, id string) {
	videoPathEnhanced := "studio/staging/" + id + "_enhanced.mp4"
	cmd := exec.Command("ffmpeg", "-i", videoPath, "-i", imagePath, "-filter_complex", "[1:v]scale=640:-1,format=rgba,colorchannelmixer=aa=0.9[img1];[0:v][img1]overlay=(main_w-overlay_w)/2:(main_h-overlay_h)/4", videoPathEnhanced)
	err := cmd.Run()

	if err != nil {
		log.Fatal(err)
	}
	os.Remove(videoPath)
}

func AddTwoImagesToVideo(videoPath string, imagePath1 string, imagePath2 string, id string) {

	videoPathEnhanced := "studio/staging/" + id + "pre_enhanced.mp4"
	cmd := exec.Command("ffmpeg", "-i", videoPath, "-i", imagePath1, "-filter_complex", "[1:v]scale=640:-1,format=rgba,colorchannelmixer=aa=0.9[img1];[0:v][img1]overlay=(main_w-overlay_w)/2:(main_h-overlay_h)/4", videoPathEnhanced)
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	os.Remove(videoPath)
	// secomd cmmand: ffmpeg -i output1.mp4 -i screenshots/18qimw8_1.jpg -filter_complex "[1:v]scale=640:-1,format=rgba,colorchannelmixer=aa=0.9[img2];[0:v][img2]overlay=(main_w-overlay_w)/2:(main_h-overlay_h)/2" output.mp4
	finalOutput := "studio/staging/" + id + "_enhanced.mp4"
	cmd = exec.Command("ffmpeg", "-i", videoPathEnhanced, "-i", imagePath2, "-filter_complex", "[1:v]scale=640:-1,format=rgba,colorchannelmixer=aa=0.9[img2];[0:v][img2]overlay=(main_w-overlay_w)/2:(main_h-overlay_h)/2", finalOutput)
	err = cmd.Run()

	if err != nil {
		log.Fatal(err)
	}
	os.Remove(videoPathEnhanced)
}

func GetAllImagesByID(id string) []string {
	var images []string
	files, err := os.ReadDir("screenshots")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if strings.Contains(f.Name(), id) {
			images = append(images, "screenshots/"+f.Name())
		}
	}
	return images
}

func ConcatAllVideos() string {
	files, err := os.ReadDir("studio/staging")
	if err != nil {
		log.Fatal(err)
	}

	// Create a list of files for FFmpeg
	var fileList strings.Builder
	// add intro video:
	fileList.WriteString("file 'studio/news_intro.mp4'\n")

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".mp4") {
			fileList.WriteString("file '" + "studio/staging/" + file.Name() + "'\n")
		}
	}

	// Write the list to a file
	err = os.WriteFile("filelist.txt", []byte(fileList.String()), 0644)
	if err != nil {
		log.Fatal(err)
	}

	videoPreSound := "studio/staging/presoundresult.mp4"
	// Concatenate the videos
	cmd := exec.Command("ffmpeg", "-f", "concat", "-safe", "0", "-i", "filelist.txt", "-c", "copy", videoPreSound)
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	return videoPreSound
}

func AddAudioToVideo(videoFile string, audioFile string, outputFile string) {
	// Create the FFmpeg command
	cmd := exec.Command("ffmpeg", "-i", videoFile, "-i", audioFile, "-c:v", "copy", "-c:a", "aac", "-map", "0:v:0", "-map", "1:a:0", outputFile)

	// Run the command
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
