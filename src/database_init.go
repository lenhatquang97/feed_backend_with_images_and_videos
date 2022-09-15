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
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Feed struct {
	FeedId         string   `json:"feedId"`
	ImageAndVideos []string `json:"imagesAndVideos"`
}

type Result struct {
	Status string `json:"status"`
}

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
	r.POST("/feeds", PostFeed)
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
func PostFeed(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var feed Feed
	if err := c.BindJSON(&feed); err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
		return
	}

	feed.FeedId = uuid.New().String()

	fmt.Println(feed.ImageAndVideos[0])

	result, err := feedCollection.InsertOne(ctx, feed)
	if err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	}
	println(result)
	c.JSON(200, feed)
}
func DeleteFeed(c *gin.Context) {

}

func UploadFeed(c *gin.Context) {
	limit_err := c.Request.ParseMultipartForm(32 << 20)

	//Corner case: Only allow to upload file <= 32mb
	if limit_err != nil {
		log.Fatal(limit_err)
	}
	formdata := c.Request.MultipartForm
	feedId := formdata.Value["feedId"]
	files := formdata.File["upload"]

	//Corner case: Feed id must only be once and unique
	if len(feedId) > 1 {
		c.AbortWithStatus(404)
	}

	baseFolder := "files/" + feedId[0] + "/"

	//Create Folder based on id
	os.Mkdir(baseFolder, 0755)

	for _, file := range files {
		filename := filepath.Base(file.Filename)
		if err := c.SaveUploadedFile(file, baseFolder+filename); err != nil {
			c.String(http.StatusBadRequest, "upload file err: %s", err.Error())
			return
		}
	}
	c.JSON(200, Result{Status: "Successful!"})

}
