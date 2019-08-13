package server

import (
	"net/http"
	"strings"

	"github.com/spf13/viper"
)

type neuteredFileSystem struct {
	fs http.FileSystem
}

// return 404 not found if index.html do not exist when passing url with path suffix /
func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if s != nil && s.IsDir() {
		index := strings.TrimSuffix(path, "/") + "/index.html"
		if _, err := nfs.fs.Open(index); err != nil {
			return nil, err
		}
	}

	return f, nil
}

func (app *application) routes() *http.ServeMux {

	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	fileServer := http.FileServer(neuteredFileSystem{http.Dir(viper.GetString("server.staticDir"))})
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}
