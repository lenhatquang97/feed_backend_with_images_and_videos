package main

import (
	"os"

	"example.com/feed_backend/src/db"
	"example.com/feed_backend/src/testing"
)

func main() {
	args := os.Args[1]
	if args == "test" {
		res := testing.GetAllFeedsWithTesting()
		numsOfFeeds := testing.GenerateNumbersOfFeeds(res, 2)
		testing.InsertFeedByFeed(numsOfFeeds)
	} else {
		db.ConnectDB()
		db.InitializeAPI()
	}
}
