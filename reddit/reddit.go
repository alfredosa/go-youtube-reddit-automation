package reddit

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/alfredosa/go-youtube-reddit-automation/config"
	"github.com/alfredosa/go-youtube-reddit-automation/utils"
	"github.com/vartanbeno/go-reddit/v2/reddit"
)

var ctx = context.Background()

func PullLatestNews(config config.Config) ([]*reddit.Post, error) {
	posts, resp, err := reddit.DefaultClient().Subreddit.TopPosts(ctx, "news", &reddit.ListPostOptions{
		ListOptions: reddit.ListOptions{
			Limit: 20,
		},
		Time: "day",
	})
	if err != nil {
		return nil, err
	}

	// TODO! FILTER OUT POSTS THAT ARE ALREADY POSTED. DB needs to be setup first

	processedPosts := CreateTTSAndSSFiles(posts, config)
	fmt.Printf("resp: %s", resp.After)
	save_posts_to_json(processedPosts)

	return processedPosts, nil
}

func save_posts_to_json(posts []*reddit.Post) {

	for _, post := range posts {
		filename := post.ID + ".json"
		if !utils.CheckFileExists(filename, "data") {
			json_file, err := json.Marshal(posts)
			if err != nil {
				panic(err)
			}
			os.WriteFile("data/"+filename, json_file, 0644)
		}
	}
}
