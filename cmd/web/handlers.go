package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

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
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	expires := r.PostForm.Get("expires")

	errors := make(map[string]string)
	if strings.TrimSpace(title) == "" {
		errors["title"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(title) > 100 {
		errors["title"] = "This field is too long (maximum is 100 characters)"
	}
	if strings.TrimSpace(content) == "" {
		errors["content"] = "This field cannot be blank"
	}
	if strings.TrimSpace(expires) == "" {
		errors["expires"] = "This field cannot be blank"
	} else if expires != "365" && expires != "7" && expires != "1" {
		errors["expires"] = "This field is invalid"
	}

	if len(errors) > 0 {
		app.render(w, r, "create.page.tmpl", &templateData{
			FormErrors: errors,
			FormData:   r.PostForm,
		})
		return
	}

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}

func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", nil)
}
