package src

import (
	"fmt"
	"github.com/go-chi/chi"
	"net/http"
)

// FourZeroFour 404
func FourZeroFour(w http.ResponseWriter, r *http.Request) {
	tpl, err := RenderFile("404")
	if err != nil {
		fmt.Println(fmt.Sprintf("parse 404 err: %v", err))
	}

	err = tpl.ExecuteTemplate(w, "layout", nil)
	if err != nil {
		fmt.Println(fmt.Sprintf("template 404 err: %v", err))
	}
}

// StaticFiles gets the static files
func StaticFiles(r chi.Router, path string, root http.FileSystem) {
	fs := http.StripPrefix(path, http.FileServer(root))
	if path != "/" && path[len(path) - 1] != '/' {
		r.Get(path, http.RedirectHandler(path + "/", 301).ServeHTTP)
	}
	path += "*"
	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}
