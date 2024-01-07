package instagram

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/barthr/newsapi"
	"github.com/charmbracelet/log"

	"github.com/alfredosa/go-youtube-reddit-automation/config"
	"github.com/alfredosa/go-youtube-reddit-automation/youtube"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
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

func UploadFileToS3(sess *session.Session, bucketName, filePath, key string) error {
	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(sess)

	// Open the file for use
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Upload the file to S3.
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
		Body:   file,
	})
	time.Sleep(15 * time.Second)
	log.Info("Waiting for file to be uploaded to S3")
	return err
}

func UploadInstagramVideo(config config.Config, posts []newsapi.Article) {

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("eu-west-2"), // replace with your preferred region
	}))

	err := UploadFileToS3(sess, "one-tech-stack-prod", "studio/staging/resultwsound.mp4", "resultwsound.mp4")
	if err != nil {
		log.Fatal(err)
	}

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
	offset := 10                      // your offset as an integer
	offsetStr := strconv.Itoa(offset) // convert to string
	data.Set("thumb_offset", offsetStr)

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

	data = url.Values{}
	apiURL = "https://graph.facebook.com/v18.0/" + igUserID + "/media_publish"
	data.Set("creation_id", container.ID)
	data.Set("access_token", accessToken)

	for {
		resp, err := http.Post(apiURL, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)

		var mediaPublish400 MediaPublish400
		err = json.Unmarshal(bodyBytes, &mediaPublish400)
		if err != nil {
			log.Fatal(err)
		}

		if strings.Contains(mediaPublish400.Error.Message, "Media ID is not available") {
			log.Info("Media ID is not available yet, sleeping for 10 seconds")
			time.Sleep(10 * time.Second)
			continue
		}

		log.Info("Instagram response: ", resp.Status, bodyString)
		break
	}
	os.Remove("studio/staging/resultwsound.mp4")
	log.Info("Cleaning, after success from instagram.")
}
