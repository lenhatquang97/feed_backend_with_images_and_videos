package main

import "example.com/feed_backend/src/db"

func main() {
	db.ConnectDB()
	db.InitializeAPI()

}
