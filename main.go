package main

import (
	"os"
	"strconv"

	"example.com/feed_backend/src/db"
	"example.com/feed_backend/src/testing"
)

func main() {
	switch len(os.Args) {
	case 1:
		db.ConnectDB()
		db.InitializeAPI()
	case 2:
		args := os.Args[1]
		if args == "generate" {
			numsOfFeeds := testing.GenerateNumbersOfFeeds(10)
			testing.InsertFeedByFeed(numsOfFeeds)
		} else if args == "eraseAll" {
			testing.DeleteAll()
		} else if args == "deleteDuplicate" {
			testing.DeleteDuplicate()
		}
	case 3:
		args := os.Args[1]
		nums := os.Args[2]
		if args == "generate" {
			actualNums, err := strconv.Atoi(nums)
			if err != nil {
				panic(err)
			}
			numsOfFeeds := testing.GenerateNumbersOfFeeds(actualNums)
			testing.InsertFeedByFeed(numsOfFeeds)
		}
	}
}
