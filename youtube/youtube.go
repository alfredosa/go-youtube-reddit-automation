package youtube

import (
	"context"
	"encoding/json"
	"os"
	"strings"

	"github.com/alfredosa/go-youtube-reddit-automation/config"
	"github.com/charmbracelet/log"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type WebConfig struct {
	Web struct {
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		RedirectURL  string `json:"redirect_url"`
		AuthURI      string `json:"auth_uri"`
		TokenURI     string `json:"token_uri"`
	} `json:"web"`
}

func YoutubeUpload(config config.Config, title, description, category, privacy string, keywords string, filename string) {
	ctx := context.Background()
	b, err := os.ReadFile("google-youtube.json")
	if err != nil {
		log.Fatalf("Unable to read service account key file: %v", err)
	}

	// Create an oauth2 config from the client secret file
	// Load the redirect URL from the JSON file
	var conf WebConfig
	if err := json.Unmarshal(b, &conf); err != nil {
		log.Fatalf("Unable to parse JSON file: %v", err)
	}

	// Create an oauth2 config from the client secret file
	webconf := &oauth2.Config{
		ClientID:     conf.Web.ClientID,
		ClientSecret: conf.Web.ClientSecret,
		RedirectURL:  conf.Web.RedirectURL,
		Scopes:       []string{youtube.YoutubeUploadScope},
		Endpoint: oauth2.Endpoint{
			AuthURL:  conf.Web.AuthURI,
			TokenURL: conf.Web.TokenURI,
		},
	}

	// Replace "YOUR_REFRESH_TOKEN" with your actual refresh token
	token := &oauth2.Token{RefreshToken: config.Goggle.Refresh_Token}

	// Create a YouTube service with the token source
	service, err := youtube.NewService(ctx, option.WithTokenSource(webconf.TokenSource(ctx, token)))
	if err != nil {
		log.Fatalf("Unable to create YouTube service: %v", err)
	}
	upload := &youtube.Video{
		Snippet: &youtube.VideoSnippet{
			Title:       title,
			Description: description,
			CategoryId:  category,
		},
		Status: &youtube.VideoStatus{
			PrivacyStatus: privacy,
			Embeddable:    true,
			MadeForKids:   true},
	}

	// The API returns a 400 Bad Request response if tags is an empty string.
	if strings.Trim(keywords, "") != "" {
		upload.Snippet.Tags = strings.Split(keywords, ",")
	}

	part := []string{"snippet", "status"}
	call := service.Videos.Insert(part, upload)

	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Error opening %v: %v", filename, err)
	}
	defer file.Close()
	if err != nil {
		log.Fatalf("Error opening %v: %v", filename, err)
	}

	response, err := call.Media(file).Do()
	if err != nil {
		log.Fatalf("Error making YouTube API call: %v", err)
	}
	log.Info("Upload successful! Video ID: ", "video", response.Id)
	os.Remove(filename)

}
