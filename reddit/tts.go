package reddit

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"sync"

	"github.com/Vernacular-ai/godub"
	"github.com/Vernacular-ai/godub/converter"
	"github.com/alfredosa/go-youtube-reddit-automation/config"
	"github.com/alfredosa/go-youtube-reddit-automation/utils"
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

			if utils.CheckFileExists(post.ID, "audio") {
				return
			}

			speech, err := speech.CreateSpeechFile(post.Title, post.ID)

			if err != nil {
				log.Fatal(err)
			}

			TakeScreenShot(post.Title, post.ID, config)
			log.Printf("Created audio file %s", speech)
		}(post)

	}

	wg.Wait()

	if utils.CheckFileExists("final_cut", "audio/result") {
		log.Println("Final cut already exists, skipping")
	} else {
		log.Printf("Finished generating audio files, now concatenating them")
		ConcatAllAudiosWithPause()
	}
}

func GetMP3Length(filename string) (int, error) {
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

		segmentSilence, err := godub.NewLoader().Load("studio/1sec_silence.mp3")

		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Concatenating %s", audio)
		segment, err = segment.Append(segment2, segmentSilence)

		if err != nil {
			log.Fatal(err)
		}
	}

	log.Printf("Concatenating %s", filePath)
	newPth := path.Join("audio", "result", "final_cut.mp3")
	err := godub.NewExporter(newPth).WithDstFormat("mp3").WithBitRate(converter.MP3BitRatePerfect).Export(segment)

	if err != nil {
		log.Fatal(err)
	}
}
