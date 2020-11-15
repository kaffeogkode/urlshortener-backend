package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/pat"
	"gorm.io/gorm"
)

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

	return router
}
