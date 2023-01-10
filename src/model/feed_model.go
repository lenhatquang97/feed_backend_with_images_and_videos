package model

import (
	"math/rand"
	"strconv"
	"time"

	"example.com/feed_backend/src/randomized"
	"example.com/feed_backend/src/utility"
	"github.com/google/uuid"
)

type Feed struct {
	FeedId         string   `json:"feedId"`
	Name           string   `json:"name"`
	Avatar         string   `json:"avatar"`
	CreatedTime    string   `json:"createdTime"`
	Caption        string   `json:"caption"`
	ImageAndVideos []string `json:"imagesAndVideos"`
	FirstWidth     int      `json:"firstWidth"`
	FirstHeight    int      `json:"firstHeight"`
}

type Result struct {
	Status string `json:"status"`
}

func GenerateARandomFeed() Feed {
	var feed Feed
	feed.FeedId = uuid.New().String()
	feed.Name = "Le Nhat Quang"
	feed.Avatar = "https://picsum.photos/200"
	feed.CreatedTime = strconv.FormatInt(time.Now().UnixMilli(), 10)

	randomStringLength := rand.Intn(100-30+1) + 30
	feed.Caption = randomized.RandStringRunes(randomStringLength)

	randomImagesNumber := rand.Intn(15-3+1) + 3
	feed.ImageAndVideos = utility.GenerateBatchImages(randomImagesNumber, feed.FeedId)
	feed.FirstWidth = 0
	feed.FirstHeight = 0
	return feed
}
