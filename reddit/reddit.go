package reddit

import (
	"context"

	"github.com/charmbracelet/log"

	"github.com/alfredosa/go-youtube-reddit-automation/config"
	dbmod "github.com/alfredosa/go-youtube-reddit-automation/db"
	"github.com/jmoiron/sqlx"
	"github.com/vartanbeno/go-reddit/v2/reddit"
)

var ctx = context.Background()

func PullLatestNews(config config.Config, db *sqlx.DB) ([]*reddit.Post, error) {
	posts, _, err := reddit.DefaultClient().Subreddit.TopPosts(ctx, "news", &reddit.ListPostOptions{
		ListOptions: reddit.ListOptions{
			Limit: 20,
		},
		Time: "day",
	})
	if err != nil {
		return nil, err
	}

	// TODO! FILTER OUT POSTS THAT ARE ALREADY POSTED. DB needs to be setup first
	log.Info("Found", "posts", len(posts))
	posts = dbmod.FilterPostedPosts(posts, db)
	log.Info("Found after filtering", "posts", len(posts))
	processedPosts := CreateTTSAndSSFiles(posts, config)

	return processedPosts, nil
}

// Questionable to save jsons? We have a db

// func save_posts_to_json(posts []*reddit.Post) {
// 	for _, post := range posts {
// 		filename := post.FullID + ".json"
// 		if !utils.CheckFileExists(filename, "data") {
// 			json_file, err := json.Marshal(posts)
// 			if err != nil {
// 				panic(err)
// 			}
// 			os.WriteFile("data/"+filename, json_file, 0644)
// 		}
// 	}
// }
