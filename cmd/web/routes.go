package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/justinas/alice"
)

// func (app *application) routes() http.Handler {
// 	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
// 	mux := http.NewServeMux()

// 	mux.HandleFunc("/", app.home)
// 	mux.HandleFunc("/snippet", app.showSnippet)
// 	mux.HandleFunc("/snippet/create", app.createSnippet)

// 	fileServer := http.FileServer(http.Dir("./ui/static/"))
// 	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

// 	// return app.recoverPanic(app.logRequest(secureHeaders(mux)))
// 	return standardMiddleware.Then(mux)
// }

func (app *application) routes() http.Handler {
	r := chi.NewRouter()

	r.Get("/", app.home)

	r.Route("/snippet", func(r chi.Router) {
		r.Get("/{id}", app.showSnippet)
		r.Get("/create", app.createSnippetForm)
		r.Post("/create", app.createSnippet)
	})

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	r.Handle("/static/", http.StripPrefix("/static", fileServer))

	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standardMiddleware.Then(r)
}
