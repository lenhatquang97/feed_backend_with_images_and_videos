package testing

import (
	"context"
	"fmt"
	"time"

	"example.com/feed_backend/src/db"
	"example.com/feed_backend/src/model"
	"go.mongodb.org/mongo-driver/bson"
)

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

func DeleteAll() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, err := db.FeedCollection.DeleteMany(ctx, bson.M{})
	if err == nil {
		fmt.Println(result)
	} else {
		fmt.Println(err)
	}
}
