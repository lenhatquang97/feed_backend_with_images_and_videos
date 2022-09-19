package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const OBJECT_URL = "https://customfeedbucket.s3.ap-southeast-1.amazonaws.com/"

func ConnectDB() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(EnvMongoURI()))
	if err != nil {
		log.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	//ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB")
	return client
}

var mongoDb *mongo.Client = ConnectDB()
var feedCollection *mongo.Collection = GetCollection(mongoDb, "feed")

func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("feed_database").Collection(collectionName)
	return collection
}

func initializeAPI() {
	r := gin.Default()
	r.GET("/feeds", GetAllFeeds)
	r.POST("/feeds/upload", UploadFeed)
	r.DELETE("/feeds/:id", DeleteFeed)

	r.Run(":8080")
}

func GetAllFeeds(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var feeds []Feed
	defer cancel()

	results, err := feedCollection.Find(ctx, bson.M{})

	if err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	}
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleFeed Feed
		if err := results.Decode(&singleFeed); err != nil {
			c.AbortWithStatus(404)
			fmt.Println(err)
		}
		feeds = append(feeds, singleFeed)
	}

	c.JSON(200, feeds)

}

func UploadFeed(c *gin.Context) {
	count := 0
	var feed Feed
	//Limit to 32 MB
	limit_err := c.Request.ParseMultipartForm(32 << 20)

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	//Corner case: Only allow to upload file <= 32mb
	if limit_err != nil {
		log.Fatal(limit_err)
	}
	formdata := c.Request.MultipartForm
	feedId := formdata.Value["feedId"]
	caption := formdata.Value["caption"]
	files := formdata.File["upload"]

	//Corner case: Feed id must only be once and unique
	if len(feedId) > 1 {
		fmt.Println("Oh no! FeedId is larger than 1")
		c.AbortWithStatus(404)
	}

	baseFolder := "files/" + feedId[0] + "/"

	//Create Folder based on id
	os.Mkdir(baseFolder, 0755)
	feed.Caption = caption[0]
	feed.FeedId = feedId[0]

	for _, file := range files {
		filename := filepath.Base(file.Filename)
		if err := c.SaveUploadedFile(file, baseFolder+filename); err != nil {
			c.String(http.StatusBadRequest, "upload file err: %s", err.Error())
			return
		}
		//upload to S3 Storage
		response, err := uploadFiles(c, baseFolder+filename)
		if err != nil {
			c.String(http.StatusBadRequest, "upload file err: %s", err.Error())
			return
		}
		count += 1
		fmt.Println(count)
		feed.ImageAndVideos = append(feed.ImageAndVideos, response.URL)
	}
	result, err := feedCollection.InsertOne(ctx, feed)
	if err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	}
	fmt.Println(result)

	c.JSON(200, feed)

}
func DeleteFeed(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()
	feedId := c.Param("id")
	result, err := feedCollection.DeleteOne(ctx, bson.M{"feedid": feedId})
	if err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	}
	fmt.Println(result)
	c.JSON(200, result)
}
