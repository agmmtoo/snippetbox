package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"agmmtoo.me/snippetbox/pkg/models"
	"github.com/go-chi/chi/v5"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "home.page.tmpl", &templateData{
		Snippets: s,
	})
}

// load snippet based on id parameter
func (app *application) snippetCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		id, err := strconv.Atoi(chi.URLParam(r, "id"))

		if err != nil {
			app.errorLog.Println(err)
			app.notFound(w)
			return
		}

		s, err := app.snippets.Get(id)

		if err == models.ErrNoRecord {
			app.notFound(w)
			return
		} else if err != nil {
			app.serverError(w, err)
			return
		}

		ctx := context.WithValue(r.Context(), "snippet", s)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	// snippet is loaded by the middleware
	snippet := r.Context().Value("snippet").(*models.Snippet)

	app.render(w, r, "show.page.tmpl", &templateData{
		Snippet: snippet,
	})
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {

	title := "0 snail"
	content := "0 snail\nClimb Mount Fjui,\nBut slowly, slowly!\n\n- Kobayashi"
	expires := "7"

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
}

func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Create snippet form"))
}
