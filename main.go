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

	linksHandler := createLinksHandler(db)
	apiHandler := wrap(createAPIHandler(db))

	// this anonymous function splits the work between the handlers based on
	// hostname in the request
	err = http.ListenAndServe(":80", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Host {
		case "admin.l.kaffeogkode.dk":
			apiHandler.ServeHTTP(w, r)
			break
		case "l.kaffeogkode.dk":
			linksHandler.ServeHTTP(w, r)
			break
		default:
			http.Error(w, `¯\_(ツ)_/¯`, http.StatusNotFound)
		}
	}))
	if err != nil {
		panic(err)
	}
}
