package utils

import (
	"encoding/hex"
	"log"
	"os"
	"strings"
)

// check if file exists given a substring
// example:
// CheckFileExists("test", "audio")
// will return true if there is a file in the audio folder
// that contains the substring "test"
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

// GetAudios returns a slice of strings with the names of all the mp3 files
// in the "audio" folder
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

func RemoveFilesWithSubstr(substr string, dir string) {
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if strings.Contains(f.Name(), substr) {
			os.Remove(dir + f.Name())
		}
	}
}

func StringToHex(s string) string {
	return hex.EncodeToString([]byte(s))
}
