package db

import (
	"io"
	"os"
	"time"

	"github.com/alfredosa/go-youtube-reddit-automation/utils"
	"github.com/barthr/newsapi"
	"github.com/charmbracelet/log"
	"github.com/jmoiron/sqlx"
)

type DBPost struct {
	PostId      string    `json:"post_id" db:"post_id"`
	SourceId    string    `json:"source_id" db:"source_id"`
	SourceName  string    `json:"source_name" db:"source_name"`
	Author      string    `json:"author" db:"author"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	Url         string    `json:"url" db:"url"`
	UrlToImage  string    `json:"url_to_image" db:"url_to_image"`
	PublishedAt time.Time `json:"published_at" db:"published_at"`
	Content     string    `json:"content" db:"content"`
}

func GetPostByID(id string, db *sqlx.DB) (DBPost, error) {
	var post DBPost
	err := db.Get(&post, "SELECT post_id FROM posts WHERE post_id=$1", id)
	if err != nil {
		return DBPost{}, err
	}
	return post, nil
}

func FilterPostedPosts(posts []newsapi.Article, db *sqlx.DB) []newsapi.Article {
	var filteredPosts []newsapi.Article
	for _, post := range posts {
		postId := utils.StringToHex(post.Title)
		_, err := GetPostByID(postId, db)
		if err != nil {
			log.Info("Post %s not found in DB, adding to filteredPosts", postId)
			filteredPosts = append(filteredPosts, post)
		} else {
			log.Info("Post %s found in DB, skipping", postId)
		}
	}

	return filteredPosts
}

func InsertPostsFromReddit(posts []newsapi.Article, db *sqlx.DB) error {

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
		postId := utils.StringToHex(post.Title)
		_, err := db.Exec(sql, postId, post.Source.ID, post.Source.Name, post.Author, post.Title, post.Description, post.URL, post.URLToImage, post.PublishedAt, post.Content)
		if err != nil {
			log.Error("Could not insert post", "post", postId)
			return err
		}
	}

	log.Info("Inserted %d posts into DB", len(posts))

	return nil
}
