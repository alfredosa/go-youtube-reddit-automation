package news

import (
	"os"
	"os/exec"
	"path"
	"strconv"
	"sync"

	"github.com/charmbracelet/log"

	"github.com/Vernacular-ai/godub"
	"github.com/Vernacular-ai/godub/converter"
	"github.com/alfredosa/go-youtube-reddit-automation/config"
	"github.com/alfredosa/go-youtube-reddit-automation/utils"
	"github.com/hajimehoshi/go-mp3"
	htgotts "github.com/hegedustibor/htgo-tts"
	rdt "github.com/vartanbeno/go-reddit/v2/reddit"
)

func CreateTTSAndSSFiles(posts []*rdt.Post, config config.Config) []*rdt.Post {
	var processedPosts []*rdt.Post
	var wg sync.WaitGroup

	speech := htgotts.Speech{Folder: "audio", Language: config.TextToSpeechSetup.Voice_ID}

	length := 6
	maxLength := 55

	for _, post := range posts {
		wg.Add(1)

		go func(post *rdt.Post) {
			defer wg.Done()
			TakeScreenShot(post.Title, post.FullID, config)
		}(post)

		audioLength := CreateAudioFile(post, config, speech)
		if length+audioLength > maxLength {
			os.Remove("audio/" + post.FullID + ".mp3")
			log.Warn("Audio file %s is too long, skipping and all subsequent posts", post.FullID)
			break
		} else {
			length += audioLength
			processedPosts = append(processedPosts, post)
		}
	}

	wg.Wait()

	if utils.CheckFileExists("final_cut", "audio/result") {
		log.Warn("Final cut already exists, skipping")
	} else {
		log.Info("Finished generating audio files, now concatenating them")
		ConcatAllAudiosWithPause()
	}

	return processedPosts
}

func CreateAudioFile(post *rdt.Post, config config.Config, speech htgotts.Speech) int {
	if utils.CheckFileExists(post.FullID, "audio") {
		length, err := GetMP3Length("audio/" + post.FullID + ".mp3")
		if err != nil {
			log.Fatal(err)
		}
		// add 1 second for the pause
		return length + 1
	}

	audio, err := speech.CreateSpeechFile(post.Title, post.FullID)

	if err != nil {
		log.Fatal(err)
	}

	length, err := GetMP3Length(audio)
	if err != nil {
		log.Fatal(err)
	}
	// add 1 second for the pause
	return length + 1
}

func GetMP3Length(filename string) (int, error) {
	file, err := os.Open(filename)
	if err != nil {
		log.Error("Error opening file: ", err)
		return 0, err
	}
	defer file.Close()

	mp3Decoder, err := mp3.NewDecoder(file)
	if err != nil {
		log.Error("Error creating decoder: ", err)
		return 0, err
	}
	// 4 bytes per sample
	samples := mp3Decoder.Length() / 4

	// Samples divided by sample rate gives length in seconds
	audioLength := samples / int64(mp3Decoder.SampleRate())

	return int(audioLength), nil
}

func ConvertSampleRate(audioPath string, sampleRate int) error {
	outputPath := audioPath + ".converted.mp3"
	cmd := exec.Command("ffmpeg", "-i", audioPath, "-ar", strconv.Itoa(sampleRate), outputPath)
	err := cmd.Run()
	if err != nil {
		return err
	}

	// Replace the original file with the converted one
	err = os.Rename(outputPath, audioPath)
	if err != nil {
		return err
	}

	return nil
}

func ConcatAllAudiosWithPause() {
	filePath := path.Join("studio", "news_intro.mp3")
	segment, _ := godub.NewLoader().Load(filePath)

	// append news transition to news intro
	segment2, err := godub.NewLoader().Load("studio/news_transition.mp3")
	if err != nil {
		log.Fatal(err)
	}

	segment, err = segment.Append(segment2)
	if err != nil {
		log.Fatal(err)
	}

	for _, audio := range utils.GetAudios() {
		// Convert the sample rate of the audio file to 24000 Hz
		err := ConvertSampleRate(audio, 44100)
		if err != nil {
			log.Fatal(err)
		}

		segment2, err := godub.NewLoader().Load(audio)

		if err != nil {
			log.Fatal(err)
		}

		segmentSilence, err := godub.NewLoader().Load("studio/news_transition.mp3")

		if err != nil {
			log.Fatal(err)
		}

		log.Info("Concatenating %s", audio)
		segment, err = segment.Append(segment2, segmentSilence)

		if err != nil {
			log.Fatal(err)
		}
	}

	log.Info("Concatenating %s", filePath)
	newPth := path.Join("audio", "result", "final_cut.mp3")
	err = godub.NewExporter(newPth).WithDstFormat("mp3").WithBitRate(converter.MP3BitRatePerfect).Export(segment)

	if err != nil {
		log.Fatal(err)
	}
}
