package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/pat"
	"gorm.io/gorm"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func createAPIHandler(db *gorm.DB) http.Handler {
	// use pat for dynamic routing
	router := pat.New()
	router.Get("/api/links", func(w http.ResponseWriter, r *http.Request) {
		var l []link
		err := db.Model(link{}).Find(&l).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound { // db ok, but none found = return 404
				http.Error(w, `¯\_(ツ)_/¯`, http.StatusNotFound)
				return
			}

			// unknown error happened = return 500
			http.Error(w, `(╯°□°）╯︵ ┻━┻`, http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(l)
		if err != nil {
			log.Println("error encoding JSON:", err)
			http.Error(w, `(╯°□°）╯︵ ┻━┻`, http.StatusInternalServerError)
		}
	})
	router.Get("/api/visits/{hash}", func(w http.ResponseWriter, r *http.Request) {
		hashStr := r.URL.Query().Get(":hash")
		var l link
		err := db.Where(link{Hash: hashStr}).First(&l).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound { // db ok, but hash not found = return 404
				http.Error(w, `¯\_(ツ)_/¯`, http.StatusNotFound)
				return
			}

			// unknown error happened = return 500
			http.Error(w, `(╯°□°）╯︵ ┻━┻`, http.StatusInternalServerError)
			return
		}

		var visits []linkVisit
		err = db.Where(linkVisit{LinkID: l.ID}).Find(&visits).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound { // db ok, but none found = return 404
				fmt.Fprint(w, "[]")
				return
			}

			// unknown error happened = return 500
			http.Error(w, `(╯°□°）╯︵ ┻━┻`, http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(visits)
		if err != nil {
			log.Println("error encoding JSON:", err)
			http.Error(w, `(╯°□°）╯︵ ┻━┻`, http.StatusInternalServerError)
		}
	})
	router.Post("/api/link", func(w http.ResponseWriter, r *http.Request) {
		// this route creates a new link and inserts it into the db
		var posted linkPost
		err := json.NewDecoder(r.Body).Decode(&posted) // decode JSON from frontend
		if err != nil {
			log.Println("error decoding JSON:", err)
			http.Error(w, `(╯°□°）╯︵ ┻━┻`, http.StatusBadRequest)
			return
		}

		_, err = url.Parse(posted.URL) // check that the URL field is an URL
		if err != nil {
			log.Println("error parsing URL:", err)
			http.Error(w, `(╯°□°）╯︵ ┻━┻`, http.StatusBadRequest)
			return
		}

		// create some random bytes for the shortlink (unique (hopefully))
		randomData := make([]byte, 5)
		c, err := rand.Read(randomData)
		if err != nil || c != 5 {
			log.Println("error getting random data:", err)
			http.Error(w, "╰(*°▽°*)╯", http.StatusInternalServerError)
			return
		}

		// encode the bytes to make 'em URL friendly
		hash := base64.URLEncoding.EncodeToString(randomData)

		newLink := link{
			Hash: hash,
			URL:  posted.URL,
		}

		// insert in the db
		err = db.Create(&newLink).Error // error here probz be unique violation
		if err != nil {
			log.Println("error creating link in db:", err)
			http.Error(w, "╰(*°▽°*)╯", http.StatusInternalServerError)
			return
		}

		// send it back as JSON
		err = json.NewEncoder(w).Encode(newLink)
		if err != nil {
			log.Println("error encoding JSON:", err)
			http.Error(w, "╰(*°▽°*)╯", http.StatusInternalServerError)
			return
		}
	})

	// static handler for the frontend - must be the last line in the handler setup!
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./frontend/")))

	return router
}
