package youtube

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	rdt "github.com/vartanbeno/go-reddit/v2/reddit"
)

func GetVideoDescription(posts []*rdt.Post) string {
	var description string

	date := time.Now()
	description += fmt.Sprintf("Today's top breaking news (%d/%d/%d). Read more: \n\n #shorts #short \n\n", date.Year(), date.Month(), date.Day())

	for i, post := range posts {
		parts := strings.Split(post.URL, ".")
		description += strconv.Itoa(i) + ": " + post.Title + "by: " + parts[1] + " \n\n"
	}

	description += "this content is not my own, it belongs to the original authors. \n\n This is a compilation of the top breaking news of the day."
	return description
}

func GetKeywords() string {
	return "news,breaking,today,newsletter,stories"
}

func GetVideoTitle() string {
	return fmt.Sprintf("Top Breaking News of %d", time.Now().Day())
}
