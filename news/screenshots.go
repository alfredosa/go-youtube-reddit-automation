package news

import (
	"fmt"
	_ "image/png"
	"io"
	"os"

	"github.com/charmbracelet/log"

	"encoding/json"
	"net/http"
	"net/url"

	"github.com/alfredosa/go-youtube-reddit-automation/config"
	"github.com/alfredosa/go-youtube-reddit-automation/utils"
)

// TakeScreenShot takes a screenshot of the first 2 images from Google search results based on the title and saves them as id_(number).png.
func TakeScreenShot(title string, id string, config config.Config) {
	query := url.QueryEscape(title)
	// search for files, if at least one id exists, skip
	if utils.CheckFileExists(id, "screenshots") {
		return
	}
	url := "https://www.googleapis.com/customsearch/v1?key=" + config.Goggle.API_Key + "&cx=" + config.Goggle.CX + "&q=" + query + "&searchType=image&num=2"

	resp, err := http.Get(url)

	if err != nil {
		log.Fatalf("could not make API request: %v", err)
	}
	defer resp.Body.Close()

	var data struct {
		Items []struct {
			Link string `json:"link"`
		} `json:"items"`
	}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		log.Fatalf("could not decode API response: %v", err)
	}

	for i, item := range data.Items {
		imgPath := "screenshots/" + id + "_" + fmt.Sprint(i)

		_, err := DownloadFile(imgPath, item.Link)
		if err != nil {
			log.Info("could not download image: %v", err)
		}
	}
}

func DownloadFile(filepath string, url string) (string, error) {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Get the Content-Type header of the HTTP response
	contentType := resp.Header.Get("Content-Type")

	// Determine the file extension based on the Content-Type header
	var extension string
	switch contentType {
	case "image/jpeg":
		extension = ".jpg"
	case "image/png":
		extension = ".png"
	case "image/gif":
		extension = ".gif"
	default:
		return "", fmt.Errorf("unsupported content type: %s", contentType)
	}

	filepath += extension

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return extension, err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return extension, err
}
