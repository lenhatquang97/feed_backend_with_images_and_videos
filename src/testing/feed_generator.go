package testing

import (
	"context"
	"fmt"
	"time"

	"example.com/feed_backend/src/db"
	"example.com/feed_backend/src/model"
)

func GenerateNumbersOfFeeds(nums int) []model.Feed {
	var outputFeeds []model.Feed
	for i := 0; i < nums; i++ {
		inputFeed := model.GenerateARandomFeed()
		outputFeeds = append(outputFeeds, inputFeed)
		fmt.Println(inputFeed.ImageAndVideos)
	}
	return outputFeeds
}

func InsertFeedByFeed(feeds []model.Feed) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	faultCount := 0
	for _, feed := range feeds {
		result, err := db.FeedCollection.InsertOne(ctx, feed)
		if err == nil {
			fmt.Println(result)
		} else {
			faultCount += 1
			fmt.Printf("Fault count: %d\n", faultCount)
			fmt.Println(err)
		}
	}
}
