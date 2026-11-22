package main

import (
	"ascii-art-web/internal/handlers"
	"log"
	"net/http"
)

func main() {
	mux := newRouter()

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func newRouter() *http.ServeMux {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	mux.HandleFunc("/", handlers.Home)
	mux.HandleFunc("/ascii-art", handlers.Generate)

	return mux
}
