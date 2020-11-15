package main

import (
	"log"
	"net/http"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	log.Println("Connecting to database 🚀")
	db, err := gorm.Open(mysql.Open("root@tcp(127.0.0.1:3306)/kaffeogkode"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	log.Println("Performing DB migration 🤯")
	err = db.AutoMigrate(
		&link{},
		&linkVisit{},
	)
	if err != nil {
		panic(err)
	}

	log.Println("Starting server 🐱‍🏍👏")
	http.Handle("/v/", createLinksHandler(db))
	http.Handle("/api/", wrap(createAPIHandler(db)))
	err = http.ListenAndServe("127.0.0.1:9000", nil)
	if err != nil {
		panic(err)
	}
}
