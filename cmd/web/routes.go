package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()

	r.Get("/", app.home)

	r.Route("/snippet", func(r chi.Router) {
		r.Get("/create", app.createSnippetForm)
		r.Post("/create", app.createSnippet)

		r.Route("/{id}", func(r chi.Router) {
			r.Use(app.snippetCtx)
			r.Get("/", app.showSnippet)
		})
	})

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	r.Handle("/static/", http.StripPrefix("/static", fileServer))

	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standardMiddleware.Then(r)
}
