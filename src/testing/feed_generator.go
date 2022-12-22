package testing

import (
	"context"
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"time"

	"example.com/feed_backend/src/db"
	"example.com/feed_backend/src/model"
	"github.com/google/uuid"

	"go.mongodb.org/mongo-driver/bson"
)

//Function return int

func GetAllFeedsWithTesting() []model.Feed {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var feeds []model.Feed
	defer cancel()

	results, err := db.FeedCollection.Find(ctx, bson.M{})

	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleFeed model.Feed
		if err := results.Decode(&singleFeed); err != nil {
			fmt.Println(err)
			return nil
		}
		feeds = append(feeds, singleFeed)
	}
	sort.Slice(feeds, func(i, j int) bool {
		it1, _ := strconv.ParseInt(feeds[i].CreatedTime, 10, 64)
		it2, _ := strconv.ParseInt(feeds[j].CreatedTime, 10, 64)
		return it1 > it2
	})
	return feeds
}

func GenerateNumbersOfFeeds(inputFeeds []model.Feed, nums int) []model.Feed {
	var outputFeeds []model.Feed
	for i := 0; i < nums; i++ {
		randomId := uuid.New().String()
		randomIndex := rand.Intn(len(inputFeeds))
		inputFeed := inputFeeds[randomIndex]
		inputFeed.FeedId = randomId
		inputFeed.CreatedTime = strconv.FormatInt(time.Now().UnixMilli(), 10)
		outputFeeds = append(outputFeeds, inputFeed)
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

func DeleteWithId(feedId string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, err := db.FeedCollection.DeleteOne(ctx, bson.M{"feedid": feedId})
	if err == nil {
		fmt.Println(result)
	} else {
		fmt.Println(err)
	}
}

// Delete duplicate with same feedId in MongoDB
func DeleteDuplicate() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	idList, _ := db.FeedCollection.Distinct(ctx, "feedid", bson.M{})
	var addedItems []model.Feed
	var result model.Feed
	for _, feedId := range idList {
		err := db.FeedCollection.FindOne(ctx, bson.M{"feedid": feedId}).Decode(&result)
		if err == nil {
			addedItems = append(addedItems, result)
		}
	}
	fmt.Println(addedItems)

	for _, value := range idList {
		result, err := db.FeedCollection.DeleteMany(ctx, bson.M{"feedid": value})
		if err == nil {
			fmt.Println(result)
		} else {
			fmt.Println(err)
		}
	}

	InsertFeedByFeed(addedItems)
}
