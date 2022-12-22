package main

import "example.com/feed_backend/src/db"

func main() {
	db.ConnectDB()
	db.InitializeAPI()
}

// func main() {
// 	result := testing.GetAllFeedsWithTesting()
// 	numsOfFeeds := testing.GenerateNumbersOfFeeds(result, 30)
// 	testing.InsertFeedByFeed(numsOfFeeds)
// }
