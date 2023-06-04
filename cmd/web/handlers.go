package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"agmmtoo.me/snippetbox/pkg/forms"
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

	form := forms.New(r.PostForm)
	form.Required("title", "content", "expires")
	form.MaxLength("title", 100)
	form.PermittedValues("expires", "30", "7", "1")

	if !form.Valid() {
		app.render(w, r, "create.page.tmpl", &templateData{Form: form})
		return
	}

	id, err := app.snippets.Insert(form.Get("title"), form.Get("content"), form.Get("expires"))
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.session.Put(r, "flash", "Snippet successfully created!")
	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}

func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Display the user signup form")
}

func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Create a new user")
}

func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Display the user login form")
}

func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Authenticate and login the user")
}

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Logout the user")
}
