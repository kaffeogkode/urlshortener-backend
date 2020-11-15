package main

import "net/http"

func wrap(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		// if r.Method == "OPTIONS" {
		// 	w.Header().Add("Access-Control-Allow-Headers", "Authorization, Content-Type") // You can add more headers here if needed
		// 	return
		// }

		h.ServeHTTP(w, r)
	})
}
