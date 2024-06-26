package video

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"github.com/barthr/newsapi"
	"github.com/charmbracelet/log"

	"github.com/alfredosa/go-youtube-reddit-automation/config"
	"github.com/alfredosa/go-youtube-reddit-automation/news"
	"github.com/alfredosa/go-youtube-reddit-automation/utils"
)

func CreateVideo(posts []newsapi.Article, config config.Config) error {
	final_cut_duration, err := news.GetMP3Length("audio/result/final_cut.mp3")

	log.Info("Final cut duration: %d", final_cut_duration)
	if err != nil {
		log.Error("Could not get final cut duration")
		return err
	}

	var wg sync.WaitGroup

	for _, post := range posts {
		postId := utils.StringToHex(post.Title)
		wg.Add(1)

		go func(post newsapi.Article) {
			defer wg.Done()
			audio := "audio/" + postId + ".mp3"
			// Get the duration of the audio file
			duration, err := news.GetMP3Length(audio)
			if err != nil {
				log.Fatal(err, "file: ", err.Error())
			}

			// Create a video with the same duration as the audio
			CreateVideoWithLength(duration, postId, post.Title)
		}(post)
	}

	log.Info("Waiting for all videos to be created")
	wg.Wait()

	preSound := ConcatAllVideos()
	audio := "audio/result/final_cut.mp3"
	AddAudioToVideo(preSound, audio, "studio/staging/resultwsound.mp4")

	os.Remove(preSound)
	CleanUp()
	return nil
}

func CleanUp() {
	os.Remove("filelist.txt")
	dir := "studio/staging/"
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err, "file: ", err.Error())
	}

	for _, file := range files {
		if strings.Contains(file.Name(), "_enhanced.mp4") {
			os.Remove(dir + file.Name())
		}
	}

	dirname := "audio/"
	utils.RemoveFilesWithSubstr(".mp3", dirname)

	dirname = "audio/result/"
	utils.RemoveFilesWithSubstr(".mp3", dirname)
}

func CreateVideoWithLength(duration int, id string, title string) {
	actualDruration := duration + 2 // space between videos

	// get random start from 100 to 3600
	randomNumber := rand.Intn(3600-100) + 100

	studioStaging := "studio/staging/"

	if utils.CheckFileExists(id, studioStaging) {
		log.Info("Enhanced Video %s already exists, skipping", id)
		return
	}

	log.Info("Creating", "video", id)

	stuidoPath := studioStaging + id + ".mp4"
	cmd := exec.Command("ffmpeg", "-ss", strconv.Itoa(randomNumber), "-i", "studio/gta4_hd.mp4", "-t", strconv.Itoa(actualDruration), "-vf", "scale=-1:1920,crop=1080:1920:(iw-1080)/2:0", stuidoPath)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		log.Fatal(fmt.Sprint(err) + ": " + stderr.String())
	}

	videoPath := studioStaging + id + ".mp4"
	OutputPath := studioStaging + id + "_enhanced.mp4"
	log.Info("Getting Images for", "video", id)
	images := GetAllImagesByID(id)

	if len(images) == 1 {
		AddOneImageToVideo(videoPath, images[0], id, title)
	}
	if len(images) == 2 {
		AddTwoImagesToVideo(videoPath, images[0], images[1], id, title)
	}
	log.Info("Finished creating video %s", OutputPath)
}

func cleanString(s string) string {
	escapedText := strings.ReplaceAll(s, "The Morning After: ", "")
	escapedText = strings.ReplaceAll(escapedText, ":", "")
	escapedText = strings.ReplaceAll(escapedText, "'", "")
	escapedText = strings.ReplaceAll(escapedText, "=", "")
	escapedText = strings.ReplaceAll(escapedText, "\\", "")
	escapedText = strings.ReplaceAll(escapedText, "/", "")
	return escapedText
}

func CreateNewsBannerAndAdd(title string, videoPath string, id string) {
	lines := splitIntoLines(title, 64)

	// only add 1 lines max
	if len(lines) >= 2 {
		lines = lines[:1]
	}
	videoBannerPath := "studio/staging/" + id + "_banner.png"

	const baseBanner = "studio/banners/onestackbanner.png"
	// create banner with text
	bannercmd := exec.Command("ffmpeg", "-i", baseBanner, "-vf", "drawtext=fontfile=studio/font/timesnewroman.ttf:text='"+cleanString(lines[0])+"':x=(w-text_w)*4/10:y=(h-text_h)/2:fontsize=24:fontcolor=black", videoBannerPath)

	var out bytes.Buffer
	var stderr bytes.Buffer
	bannercmd.Stdout = &out
	bannercmd.Stderr = &stderr
	err := bannercmd.Run()

	if err != nil {
		log.Fatal(fmt.Sprint(err) + ": " + stderr.String())
	}

	finalVideoOutputPath := "studio/staging/" + id + "_enhanced.mp4"
	// add banner to video with ffmpeg
	cmd := exec.Command("ffmpeg", "-i", videoPath, "-i", videoBannerPath, "-filter_complex", "[1:v]scale=1060:-1,format=rgba,colorchannelmixer=aa=0.9[img2];[0:v][img2]overlay=(main_w-overlay_w)/2:(main_h-overlay_h)*6/7", finalVideoOutputPath)
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()

	if err != nil {
		log.Fatal(fmt.Sprint(err) + ": " + stderr.String())
	}

	if err != nil {
		log.Fatal(err, "file: ", err.Error())
	}
	os.Remove(videoPath)
	os.Remove(videoBannerPath)
}

func AddOneImageToVideo(videoPath string, imagePath string, id string, title string) {
	videoPathEnhanced := "studio/staging/" + id + "pre_banner_enhanced.mp4"
	cmd := exec.Command("ffmpeg", "-i", videoPath, "-i", imagePath, "-filter_complex", "[1:v]scale=640:-1,format=rgba,colorchannelmixer=aa=0.9[img1];[0:v][img1]overlay=(main_w-overlay_w)/2:(main_h-overlay_h)/4", videoPathEnhanced)
	err := cmd.Run()

	if err != nil {
		log.Fatal(err, "file: ", err.Error())
	}
	os.Remove(videoPath)
	CreateNewsBannerAndAdd(title, videoPathEnhanced, id)
}

func AddTwoImagesToVideo(videoPath string, imagePath1 string, imagePath2 string, id string, title string) {

	videoPathEnhanced := "studio/staging/" + id + "pre_enhanced.mp4"
	cmd := exec.Command("ffmpeg", "-i", videoPath, "-i", imagePath1, "-filter_complex", "[1:v]scale=640:-1,format=rgba,colorchannelmixer=aa=0.9[img1];[0:v][img1]overlay=(main_w-overlay_w)/2:(main_h-overlay_h)/4", videoPathEnhanced)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		log.Fatal(fmt.Sprint(err) + ": " + stderr.String())
	}

	os.Remove(videoPath)
	finalOutput := "studio/staging/" + id + "pre_banner_enhanced.mp4"
	cmd = exec.Command("ffmpeg", "-i", videoPathEnhanced, "-i", imagePath2, "-filter_complex", "[1:v]scale=640:-1,format=rgba,colorchannelmixer=aa=0.9[img2];[0:v][img2]overlay=(main_w-overlay_w)/2:(main_h-overlay_h)*2/3", finalOutput)

	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()

	if err != nil {
		log.Fatal(fmt.Sprint(err) + ": " + stderr.String())
	}
	os.Remove(videoPathEnhanced)
	CreateNewsBannerAndAdd(title, finalOutput, id)
}

func GetAllImagesByID(id string) []string {
	var images []string
	files, err := os.ReadDir("screenshots")
	if err != nil {
		log.Fatal(err, "file: ", err.Error())
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
		log.Fatal(err, "file: ", err.Error())
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
		log.Fatal(err, "file: ", err.Error())
	}

	videoPreSound := "studio/staging/presoundresult.mp4"
	// Concatenate the videos
	cmd := exec.Command("ffmpeg", "-f", "concat", "-safe", "0", "-i", "filelist.txt", "-c", "copy", videoPreSound)
	err = cmd.Run()
	if err != nil {
		log.Fatal(err, "file: ", err.Error())
	}
	return videoPreSound
}

func AddAudioToVideo(videoFile string, audioFile string, outputFile string) {
	// Create the FFmpeg command
	cmd := exec.Command("ffmpeg", "-i", videoFile, "-i", audioFile, "-c:v", "copy", "-c:a", "aac", "-map", "0:v:0", "-map", "1:a:0", outputFile)

	// Run the command
	err := cmd.Run()
	if err != nil {
		log.Fatal(err, "file: ", err.Error())
	}
}

func splitIntoLines(s string, maxLen int) []string {
	words := strings.Fields(s)
	var lines []string
	var line string

	for _, word := range words {
		if len(line+" "+word) <= maxLen {
			line += " " + word
		} else {
			lines = append(lines, strings.TrimSpace(line))
			line = word
		}
	}

	lines = append(lines, strings.TrimSpace(line))

	return lines
}
