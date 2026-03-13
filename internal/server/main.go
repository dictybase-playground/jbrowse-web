package main

import (
	"log"
	"net/http"
	"time"
)

func logRequest(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		h.ServeHTTP(w, r)
		log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
	})
}
func main() {
	appFiles := http.FileServer(http.Dir("./jbrowse2"))
	dataFiles := http.FileServer(http.Dir("./test_data"))

	http.Handle("/", logRequest(appFiles))
	http.Handle("/test_data/", http.StripPrefix("/test_data/", logRequest(dataFiles)))

	log.Print("Listening on :3000...")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
