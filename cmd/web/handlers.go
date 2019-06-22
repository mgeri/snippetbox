package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// check for home
	if r.URL.Path != "/" {
		app.notFound(w) // Use the notFound() helper
		return
	}

	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}
	// parse template
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err) // Use the serverError() helper.
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		app.serverError(w, err) // Use the serverError() helper.
		return
	}
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.URL.Query().Get("id"))

	if err != nil || id < 1 {
		app.notFound(w) // Use the notFound() helper.
		return
	}

	fmt.Fprintf(w, "Display a snippet with ID %d...", id)
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	// check for POST
	if r.Method != "POST" {
		// header must be set before any WriteHeader/Write
		w.Header().Set("Allow", "POST")

		app.clientError(w, http.StatusMethodNotAllowed) // Use the clientError() helper.
		return
	}
	w.Write([]byte("Create a new snippet..."))
}
