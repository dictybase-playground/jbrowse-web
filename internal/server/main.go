package server

import (
	"context"
	"log"
	"net/http"
	"time"
)

func logRequest(h http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			h.ServeHTTP(w, r)
			log.Printf(
				"%s %s %v",
				r.Method,
				r.URL.Path,
				time.Since(start),
			)
		},
	)
}

func StartServer(ctx context.Context) error {
	appFiles := http.FileServer(http.Dir("./jbrowse2"))
	dataFiles := http.FileServer(http.Dir("./test_data"))

	mux := http.NewServeMux()
	mux.Handle("/", logRequest(appFiles))
	mux.Handle(
		"/test_data/",
		http.StripPrefix("/test_data/", logRequest(dataFiles)),
	)

	srv := &http.Server{
		Addr:              ":3000",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Print("Listening on :3000...")
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down server...")

	return srv.Shutdown(ctx)
}
