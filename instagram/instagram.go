package instagram

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	rdt "github.com/vartanbeno/go-reddit/v2/reddit"

	"github.com/alfredosa/go-youtube-reddit-automation/config"
	"github.com/alfredosa/go-youtube-reddit-automation/youtube"
)

type Container struct {
	ID string `json:"id"`
}

type MediaPublish400 struct {
	Error struct {
		Message   string `json:"message"`
		Type      string `json:"type"`
		ErrorUser string `json:"error_user_title"`
	} `json:"error"`
}

func UploadInstagramVideo(config config.Config, posts []*rdt.Post) {

	igUserID := "17841463853910397"               // Your Instagram user ID
	accessToken := config.Instagram.Access_Token  // Your access token
	videoURL := config.Instagram.Post_URL         // URL of your video
	caption := youtube.GetVideoDescription(posts) // Caption for the reel

	log.Info("Uploading to Instagram", "videoURL", videoURL, "caption", caption, "igUserID", igUserID, "accessToken", accessToken)
	apiURL := "https://graph.facebook.com/v18.0/" + igUserID + "/media"

	data := url.Values{}
	data.Set("media_type", "REELS")
	data.Set("video_url", videoURL)
	data.Set("caption", caption)
	data.Set("access_token", accessToken)

	resp, err := http.Post(apiURL, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		log.Fatal(err)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)
	log.Info("Instagram response: ", resp.Status, bodyString)
	resp.Body.Close()
	// id element of the response is the carousel ID
	var container Container
	err = json.Unmarshal(bodyBytes, &container)
	if err != nil {
		log.Fatal(err)
	}
	/// THEN WE NEED TO SEND A REQUEST USING media_publish
	// sleep for 10 seconds
	time.Sleep(30 * time.Second)

	data = url.Values{}
	apiURL = "https://graph.facebook.com/v18.0/" + igUserID + "/media_publish"
	data.Set("creation_id", container.ID)
	data.Set("access_token", accessToken)

	resp, err = http.Post(apiURL, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	bodyBytes, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString = string(bodyBytes)
	log.Info("Instagram response: ", resp.Status, bodyString)
	resp.Body.Close()
}
