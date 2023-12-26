package utils

import (
	"log"
	"os"
	"strings"
)

func CheckFileExists(substr string, dir string) bool {
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if strings.Contains(f.Name(), substr) {
			return true
		}
	}
	return false
}

func GetAudios() []string {
	var audios []string
	files, err := os.ReadDir("audio")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if strings.Contains(f.Name(), "mp3") {
			audios = append(audios, "audio/"+f.Name())
		}
	}
	return audios
}
