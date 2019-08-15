package server

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/mgeri/snippetbox/store"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// check for home
	if r.URL.Path != "/" {
		app.notFound(w) // Use the notFound() helper
		return
	}

	// s, err := app.snippetStore.Latest()
	// if err != nil {
	// 	app.serverError(w, err)
	// 	return
	// }

	// for _, snippet := range s {
	//     fmt.Fprintf(w, "%v\n", snippet)
	// }

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

	// Use the SnippetModel object's Get method to retrieve the data for a
	// specific record based on its ID. If no matching record is found,
	// return a 404 Not Found response.
	s, err := app.snippetStore.Get(id)
	if err == store.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	// Write the snippet data as a plain-text HTTP response body.
	fmt.Fprintf(w, "%v", s)
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	// check for POST
	if r.Method != "POST" {
		// header must be set before any WriteHeader/Write
		w.Header().Set("Allow", "POST")

		app.clientError(w, http.StatusMethodNotAllowed) // Use the clientError() helper.
		return
	}

	// Create some variables holding dummy data. We'll remove these later on
	// during the build.
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := "7"

	// Pass the data to the SnippetModel.Insert() method, receiving the
	// ID of the new record back.
	id, err := app.snippetStore.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Redirect the user to the relevant page for the snippet.
	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)

}
