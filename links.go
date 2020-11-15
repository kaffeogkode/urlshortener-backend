package main

import (
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/pat"
	"gorm.io/gorm"
)

// createLinksHandler returns a handler that redirects to the URL's stored in
// the database. If they are not found a 404 is returned. Any other error
// returns a 500.
func createLinksHandler(db *gorm.DB) http.Handler {
	// use pat for dynamic routing
	router := pat.New()
	// the part after the slash is dynamic
	router.Get("/v/{hash}", func(w http.ResponseWriter, r *http.Request) {
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

		// spawn in new goroutine to save visit async
		go saveVisit(db, r.RemoteAddr, l.ID)

		// redirect the request to the shortened URL from the DB
		http.Redirect(w, r, l.URL, http.StatusSeeOther)
	})
	return router
}

// saveVisits writes the visitor info (IP + timestamp) to the database
func saveVisit(db *gorm.DB, remoteAddr string, linkID int64) {
	// split IP address from port number
	host, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		log.Println("error getting host from visitor:", err)
		host = "unknown"
	}

	// put data into struct
	visit := linkVisit{
		LinkID:    linkID,
		VisitTime: time.Now().Unix(),
		VisitorIP: host,
	}

	err = db.Create(&visit).Error // save in DB
	if err != nil {
		log.Println("error inserting visit into db:", err)
	}
}
