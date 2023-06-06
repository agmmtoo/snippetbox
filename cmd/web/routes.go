package main

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()

	dynamicMiddleware := alice.New(app.session.Enable, noSurf, app.authenticate)

	r.Get("/", dynamicMiddleware.ThenFunc(app.home).(http.HandlerFunc))

	r.Route("/snippet", func(r chi.Router) {
		r.Get("/create", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.createSnippetForm).(http.HandlerFunc))
		r.Post("/create", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.createSnippet).(http.HandlerFunc))

		r.Route("/{id}", func(r chi.Router) {
			r.Use(app.snippetCtx)
			r.Get("/", dynamicMiddleware.ThenFunc(app.showSnippet).(http.HandlerFunc))
		})
	})

	r.Route("/user", func(r chi.Router) {
		r.Get("/signup", dynamicMiddleware.ThenFunc(app.signupUserForm).(http.HandlerFunc))
		r.Post("/signup", dynamicMiddleware.ThenFunc(app.signupUser).(http.HandlerFunc))
		r.Get("/login", dynamicMiddleware.ThenFunc(app.loginUserForm).(http.HandlerFunc))
		r.Post("/login", dynamicMiddleware.ThenFunc(app.loginUser).(http.HandlerFunc))
		r.Post("/logout", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.logoutUser).(http.HandlerFunc))
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
