package main

type Feed struct {
	FeedId         string   `json:"feedId"`
	ImageAndVideos []string `json:"imagesAndVideos"`
}

type Result struct {
	Status string `json:"status"`
}