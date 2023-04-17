package handlers

import (
	"html/template"
	"log"
	"net/http"
	"templator/services"
)

type TemplatesPageData struct {
	Templates []services.TemplateData
}

func TemplatesHandler(w http.ResponseWriter, r *http.Request) {
	// Rendering HTML page
	data := TemplatesPageData{ Templates: services.GetAvailableTemplates() }
	err := template.Must(template.ParseFiles("web/templates.html")).Execute(w, data)
	if err != nil {
		log.Fatal(err)
	}
}


