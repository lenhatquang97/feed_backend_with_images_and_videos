package db

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"

	"example.com/feed_backend/src/cdn"
	"example.com/feed_backend/src/model"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func InitializeAPI() {
	r := gin.Default()
	r.GET("/feeds", GetAllFeeds)
	r.POST("/feeds/upload", UploadFeed)
	r.DELETE("/feeds/:id", DeleteFeed)

	r.Run(":8080")
}

func GetAllFeeds(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var feeds []model.Feed
	defer cancel()

	results, err := feedCollection.Find(ctx, bson.M{})

	if err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	}
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleFeed model.Feed
		if err := results.Decode(&singleFeed); err != nil {
			c.AbortWithStatus(404)
			fmt.Println(err)
		}
		feeds = append(feeds, singleFeed)
	}
	sort.Slice(feeds, func(i, j int) bool {
		it1, _ := strconv.ParseInt(feeds[i].CreatedTime, 10, 64)
		it2, _ := strconv.ParseInt(feeds[j].CreatedTime, 10, 64)
		return it1 > it2
	})

	c.JSON(200, feeds)

}

func UploadFeed(c *gin.Context) {
	count := 0
	var feed model.Feed
	//Limit to 120 MB
	limit_err := c.Request.ParseMultipartForm(120 << 20)

	//Timeout is 120 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	//Corner case: Only allow to upload file <= 32mb
	if limit_err != nil {
		log.Fatal(limit_err)
	}
	formdata := c.Request.MultipartForm

	feedId := formdata.Value["feedId"]
	name := formdata.Value["name"]
	avatar := formdata.Value["avatar"]
	createdTime := formdata.Value["createdTime"]
	caption := formdata.Value["caption"]
	files := formdata.File["upload"]
	//Corner case: Feed id must only be once and unique
	if len(feedId) > 1 {
		fmt.Println("Nooooo")
		fmt.Println("Oh no! FeedId is larger than 1")
		c.AbortWithStatus(404)
	}
	baseFolder := "files/" + feedId[0] + "/"

	//Create Folder based on id
	os.Mkdir(baseFolder, 0755)

	feed.FeedId = feedId[0]
	feed.Name = name[0]
	feed.Avatar = avatar[0]
	feed.CreatedTime = createdTime[0]
	feed.Caption = caption[0]
	feed.FirstWidth = 0
	feed.FirstHeight = 0

	fmt.Println("Checkpoint 1")

	for _, file := range files {
		filename := filepath.Base(file.Filename)
		if err := c.SaveUploadedFile(file, baseFolder+filename); err != nil {
			fmt.Printf("Error 2 %s\n", err.Error())
			c.String(http.StatusBadRequest, "upload file err: %s", err.Error())
			return
		}
		response, err := cdn.UploadFiles(c, baseFolder+filename)
		if count == 0 {
			feed.FirstWidth = response.Width
			feed.FirstHeight = response.Height
		}
		if err != nil {
			fmt.Printf("Error 3 %s\n", err.Error())
			c.String(http.StatusBadRequest, "upload file err: %s", err.Error())
			return
		}
		count += 1
		fmt.Printf("You have uploaded %d files\n", count)
		feed.ImageAndVideos = append(feed.ImageAndVideos, response.URL)
	}

	result, err := feedCollection.InsertOne(ctx, feed)
	if err != nil {
		fmt.Println("Checkpoint 3")
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
