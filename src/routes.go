package src

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/keloran/go-probe"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// Routes self explanatory
func Routes() chi.Router {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))

	// Probe
	router.Get("/probe", probe.HTTP)
	router.Get(fmt.Sprintf("%s/probe", os.Getenv("SITE_PREFIX")), probe.HTTP)

	// Routes
	router.Route(fmt.Sprintf("%s/", os.Getenv("SITE_PREFIX")), func(r chi.Router) {
		r.HandleFunc("/", renderPage)
	})

	// Static
	staticDir := "static"
	StaticFiles(router, fmt.Sprintf("/%s", staticDir), http.Dir(staticDir))

	// Invalid paths
	router.NotFound(FourZeroFour)
	router.MethodNotAllowed(FourZeroFour)

	return router
}

// RenderFile template parse
func RenderFile(fileName string) (*template.Template, error) {
	tplDir := "./static/tpl"
	pagesDir := "./static/pages"

	// pages itself
	pages := filepath.Join(pagesDir, fmt.Sprintf("%s.html", fileName))

	// layout stuff
	layout := filepath.Join(tplDir, "layout.html")
	header := filepath.Join(tplDir, "header.html")
	footer := filepath.Join(tplDir, "footer.html")

	return template.ParseFiles(layout, header, footer, pages)
}

func renderPage(w http.ResponseWriter, r *http.Request) {
	// pages itself
	path := filepath.Clean(r.URL.Path)
	if path == "/" {
		path = "index"
	}

	tpl, err := RenderFile(path)
	if err != nil {
		fmt.Println(fmt.Sprintf("parse err: %v", err))
	}

	err = tpl.ExecuteTemplate(w, "layout", nil)
	if err != nil {
		fmt.Println(fmt.Sprintf("template err: %v", err))
	}
}
