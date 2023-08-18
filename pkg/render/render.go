package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/simonetome/bookings/pkg/config"
	"github.com/simonetome/bookings/pkg/models"
)

var app *config.AppConfig

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultdata(td *models.TemplateData) *models.TemplateData {
	// I can add some data available to every page ind. from hand logic
	return td
}

func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {

	var tc map[string]*template.Template
	// create a template cache
	if app.UseCache {
		tc = app.TemplateCache
	} else {
		// rebuild the cache at the moment of the request
		// (eq to reading from disk)
		tc, _ = CreateTemplateCache()
	}

	// get tempate from cache
	t, _ := tc[tmpl]

	// more fine-grained error handling
	buf := new(bytes.Buffer)
	td = AddDefaultdata(td)
	err := t.Execute(buf, td)
	if err != nil {
		log.Println(err)
	}

	// render the template
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println(err)
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	// alternative to create a map
	myCache := map[string]*template.Template{}
	// I want to populate the entire cache with everything I have

	// get all page.tmpl from templates folder
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	// range through pages
	// render happens in this order page -> layout
	for _, page := range pages {
		// page is the full pathName
		name := filepath.Base(page)
		// create a new template with name "name" which is the filename
		ts, err := template.New(name).ParseFiles(page)

		if err != nil {
			return myCache, err
		}
		// Glob search for files that match a pattern in the folder
		// and return the full path
		layouts, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}

		// its not mandatory to have layouts for each page
		// parse glob parse everything that matches a pattern
		if len(layouts) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}
	return myCache, nil
}
