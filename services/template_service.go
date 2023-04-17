package services

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const templatesFolder = "templates"

type TemplateData struct {
	Title string
	TemplateName string
	Type string
}

func getTemplatePath(templateName string) string {
	return fmt.Sprintf("%s/%s", templatesFolder, templateName)
}

func GetTemplateData(templateName string) (TemplateData, error) {
	templatePath := getTemplatePath(templateName)

	// Reading HTML file
	buf, err := ioutil.ReadFile(templatePath)
	if err != nil {
		return TemplateData{}, err
	}

	// Parsing the HTML document
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(buf)))
	if err != nil {
		return TemplateData{}, err
	}

	// Looking for document's "Title" tag
	var title = ""
	selector := doc.Find("html > head > title")
	if selector.Length() > 0 {
		title = selector.Text()
	}

	return TemplateData { Title: title, TemplateName: templateName, Type: filepath.Ext(templatePath) }, nil
}

func GetAvailableTemplates() []TemplateData {
	var templates []TemplateData

	err := filepath.Walk(templatesFolder, func(path string, info os.FileInfo, err error) error {
		if info != nil && !info.IsDir() && filepath.Ext(info.Name()) == ".html" {
			// Collecting template's data
			templateData, err := GetTemplateData(info.Name())
			if err == nil {
				// Adding a new item into the list of templates
				templates = append(templates, templateData)
			} else {
				log.Printf("Unable to collect template's data ('%s'): %s", info.Name(), err)
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	return templates
}

func CollectTemplateParams(templateName string) ([]string, error) {
	// Reading input file
	buf, err := ioutil.ReadFile(getTemplatePath(templateName))
	if err != nil {
		return make([]string, 0), err
	}
	str := string(buf)

	// Looking for all ${...} occurrences in the source string
	list := regexp.MustCompile(`\${([0-9A-ZА-Яa-zа-я.,/_ ]+)}`).FindAllStringSubmatch(str, -1)
	// Getting only second elements of the resulting list
	var list2 []string
	for _, x := range list {
		if len(x) > 1 {
			list2 = append(list2, x[1])
		}
	}

	// Removing duplicates from the list
	list2 = removeDuplicateStr(list2, true)

	return list2, nil
}

func removeDuplicateStr(strSlice []string, caseSensitive bool) []string {
	allKeys := make(map[string]bool)
	var list []string
	for _, item := range strSlice {
		var key string
		if !caseSensitive {
			key = strings.ToLower(item)
		} else {
			key = item
		}
		if _, value := allKeys[key]; !value {
			allKeys[key] = true
			list = append(list, item)
		}
	}
	return list
}

func GenerateDocument(templateName string, params map[string]string) string {
	tmpFile, err := ioutil.TempFile("", "templator-html-document-*.html")
	if err != nil {
		log.Fatal(err)
	}
	fileReplaceAllStrings(getTemplatePath(templateName), tmpFile.Name(), params)

	return tmpFile.Name()
}

func fileReplaceAllStrings(inputFile string, outputFile string, params map[string]string) {
	// Reading input file
	buf, err := ioutil.ReadFile(inputFile)
	if err != nil {
		log.Fatal(err)
	}
	documentBody := string(buf)

	// Replacing all the placeholders with their values
	for k, v := range params {
		documentBody = strings.ReplaceAll(documentBody, fmt.Sprintf("${%s}", k), v)
	}

	// Writing output file
	err = ioutil.WriteFile(outputFile, []byte(documentBody), 0644)
	if err != nil {
		log.Fatal(err)
	}
}
