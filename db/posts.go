package db

import (
	"io"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/jmoiron/sqlx"
	rdt "github.com/vartanbeno/go-reddit/v2/reddit"
)

type DBPost struct {
	ID                    string     `db:"id"`
	FullID                string     `db:"full_id"`
	Created               *time.Time `db:"created"`
	Edited                *time.Time `db:"edited"`
	Permalink             string     `db:"permalink"`
	URL                   string     `db:"url"`
	Title                 string     `db:"title"`
	Body                  string     `db:"body"`
	Likes                 *bool      `db:"likes"`
	Score                 int        `db:"score"`
	UpvoteRatio           float32    `db:"upvote_ratio"`
	NumberOfComments      int        `db:"number_of_comments"`
	SubredditName         string     `db:"subreddit_name"`
	SubredditNamePrefixed string     `db:"subreddit_name_prefixed"`
	SubredditID           string     `db:"subreddit_id"`
	SubredditSubscribers  int        `db:"subreddit_subscribers"`
	Author                string     `db:"author"`
	AuthorID              string     `db:"author_id"`
	Spoiler               bool       `db:"spoiler"`
	Locked                bool       `db:"locked"`
	NSFW                  bool       `db:"nsfw"`
	IsSelfPost            bool       `db:"is_self_post"`
	Saved                 bool       `db:"saved"`
	Stickied              bool       `db:"stickied"`
	Posted                bool       `db:"posted"`
}

func GetPostByID(id string, db *sqlx.DB) (DBPost, error) {
	var post DBPost
	err := db.Get(&post, "SELECT full_id FROM posts WHERE full_id=$1", id)
	if err != nil {
		return DBPost{}, err
	}
	return post, nil
}

func FilterPostedPosts(posts []*rdt.Post, db *sqlx.DB) []*rdt.Post {
	var filteredPosts []*rdt.Post
	for _, post := range posts {
		_, err := GetPostByID(post.FullID, db)
		if err != nil {
			log.Info("Post %s not found in DB, adding to filteredPosts", post.FullID)
			filteredPosts = append(filteredPosts, post)
		} else {
			log.Info("Post %s found in DB, skipping", post.FullID)
		}
	}

	return filteredPosts
}

func InsertPostsFromReddit(posts []*rdt.Post, db *sqlx.DB) error {

	file := "sql/transactions/insert_post.sql"

	openedFile, err := os.Open(file)
	if err != nil {
		log.Error("Could not open file", "file", file)
		return err
	}
	defer openedFile.Close()

	content, err := io.ReadAll(openedFile)
	if err != nil {
		log.Error("Could not read file", "file", file)
		return err
	}

	sql := string(content)

	for _, post := range posts {
		_, err := db.Exec(sql, post.ID, post.FullID, post.Created.Time, post.Edited.Time, post.Permalink, post.URL, post.Title, post.Body, post.Likes, post.Score, post.UpvoteRatio, post.NumberOfComments, post.SubredditName, post.SubredditNamePrefixed, post.SubredditID, post.SubredditSubscribers, post.Author, post.AuthorID, post.Spoiler, post.Locked, post.NSFW, post.IsSelfPost, post.Saved, post.Stickied)
		if err != nil {
			log.Error("Could not insert post", "post", post.FullID)
			return err
		}
	}

	log.Info("Inserted %d posts into DB", len(posts))

	return nil
}
