package handlers

import (
	"fmt"
	"net/http"

	"github.com/simonetome/bookings/pkg/config"
	"github.com/simonetome/bookings/pkg/models"
	"github.com/simonetome/bookings/pkg/render"
)

// Using repository pattern

// Repository
var Repo *Repository

type Repository struct {
	App *config.AppConfig
}

// Creates a repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// Sets repository for the handlers
func NewHandler(r *Repository) {
	Repo = r
}

// Link all the handlers to the repository with a receiver
// Home handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {

	remoteIp := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIp)

	render.RenderTemplate(w, "home.page.tmpl", &models.TemplateData{})
}

// About page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	// logic
	stringMap := make(map[string]string)
	stringMap["test"] = "Test data"

	remoteIp := m.App.Session.Get(r.Context(), "remote_ip")
	fmt.Println(remoteIp)
	// send logic to template
	render.RenderTemplate(w, "about.page.tmpl", &models.TemplateData{StringMap: stringMap})
}
