package reddit

import (
	"context"

	"github.com/alfredosa/go-youtube-reddit-automation/config"
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

	// for _, post := range posts {
	// 	fmt.Printf("post Title: %s\n", post.Title)
	// 	fmt.Printf("post URL: %s\n", post.URL)
	// 	fmt.Printf("post BODY: %s\n", post.Body)
	// 	fmt.Printf("\n")
	// }

	println("After from response: %s", resp.After)
	return posts, nil
}
