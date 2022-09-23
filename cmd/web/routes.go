package main

import "net/http"
func (app *application) routes() *http.ServeMux {
	// Use the http.NewServeMux() function to initialize a new servemux
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	// Create a file server which serves files out of the "./ui/static" directory.
	fileserver := http.FileServer(http.Dir("./ui/static/"))

	// For matching paths, we strip the "/static" prefix before the request reaches the file server.
	mux.Handle("/static/", http.StripPrefix("/static", fileserver))

	return mux
}