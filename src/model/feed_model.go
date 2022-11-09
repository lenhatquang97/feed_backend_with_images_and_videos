package model

type Feed struct {
	FeedId         string   `json:"feedId"`
	Name           string   `json:"name"`
	Avatar         string   `json:"avatar"`
	CreatedTime    string   `json:"createdTime"`
	Caption        string   `json:"caption"`
	ImageAndVideos []string `json:"imagesAndVideos"`
	FirstWidth     int      `json:"firstWidth"`
	FirstHeight    int      `json:"firstHeight"`
}

type Result struct {
	Status string `json:"status"`
}
