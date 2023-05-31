package main

import (
	"net/http"
	"strings"

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

	filesDir := http.Dir("./ui/static")
	fileServer(r, "/static", filesDir)

	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standardMiddleware.Then(r)
}

func fileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
