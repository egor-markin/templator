package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"templator/services"
)

func GeneratorHandler(w http.ResponseWriter, r *http.Request) {
	// Parsing provided params
	err := r.ParseForm()
	if err != nil {
		returnError(w, fmt.Sprintf("Unable to parse provided form data: %s", err))
		return
	}

	// Figuring out 'template' param value
	templateName := r.Form.Get("template")
	if len(strings.TrimSpace(templateName)) == 0 {
		returnError(w, fmt.Sprintf("Required parameter '%s' was not provided", "template"))
		return
	}

	// Figuring out a list of required parameters for the provided template
	params, err := services.CollectTemplateParams(templateName)
	if err != nil {
		returnError(w, fmt.Sprintf("Unable to figure out a list of required parameters for the provided template ('%s'): %s", templateName, err))
		return
	}

	// Saving provided params into a map
	paramsMap := make(map[string]string)
	for _, paramName := range params {
		paramValue := r.Form.Get(paramName)
		if len(strings.TrimSpace(paramValue)) == 0 {
			paramsMap[paramName] = ""
			//_, err := w.Write([]byte(fmt.Sprintf("required parameter '%s' was not provided", paramName)))
			//if err != nil {
			//	log.Fatal(err)
			//}
			//return
		} else {
			paramsMap[paramName] = strings.TrimSpace(paramValue)
		}
	}

	// Generating a new document from the template
	newHtmlDocumentPath := services.GenerateDocument(templateName, paramsMap)
	defer os.Remove(newHtmlDocumentPath)

	// Converting generated HTML file into PDF
	newPdfDocumentPath := services.ConvertHtmlToPdf(newHtmlDocumentPath)
	defer os.Remove(newPdfDocumentPath)

	// Collecting template's meta data
	templateData, err := services.GetTemplateData(templateName)
	if err != nil {
		returnError(w, fmt.Sprintf("Unable to collect provided template's data ('%s'): %s", templateName, err))
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.pdf\"", templateData.Title))
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeFile(w, r, newPdfDocumentPath)
}

func returnError(w http.ResponseWriter, errorMessage string) {
	_, err := w.Write([]byte(errorMessage))
	if err != nil {
		log.Fatal(err)
	}
}