package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"templator/services"
)

type paramsPageData struct {
	TemplateTitle 	string
	TemplateName 	string
	Params     		[]string
}

func ParamsHandler(w http.ResponseWriter, r *http.Request) {
	// Figuring out 'template' param value
	templateNameValues, present := r.URL.Query()["template"]
	if !present || len(templateNameValues) == 0 {
		_, err := w.Write([]byte(fmt.Sprintf("required parameter '%s' was not provided", "template")))
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	templateName := templateNameValues[0]

	// Collecting requested template's data
	templateData, err := services.GetTemplateData(templateName)
	if err != nil {
		returnError(w, fmt.Sprintf("Unable to collect provided template's data ('%s'): %s", templateName, err))
		return
	}

	// Gathering a list of required parameters for the selected template
	params, err := services.CollectTemplateParams(templateName)
	if err != nil {
		returnError(w, fmt.Sprintf("Unable to figure out a list of required parameters for the provided template ('%s'): %s", templateName, err))
		return
	}

	// Rendering HTML page
	data := paramsPageData { TemplateTitle: templateData.Title, TemplateName: templateName, Params: params }
	err = template.Must(template.ParseFiles("web/params.html")).Execute(w, data)
	if err != nil {
		log.Fatal(err)
	}
}


