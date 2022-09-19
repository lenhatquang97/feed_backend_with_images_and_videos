package main

type Feed struct {
	FeedId         string   `json:"feedId"`
	Caption        string   `json:"caption"`
	ImageAndVideos []string `json:"imagesAndVideos"`
}

type Result struct {
	Status string `json:"status"`
}
