package services

import (
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"io/ioutil"
	"log"
)

func ConvertHtmlToPdf(htmlFilePath string) string {
	// PDF generator initialization
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		log.Fatal(err)
	}

	// Set global options
	pdfg.Dpi.Set(300)
	pdfg.Orientation.Set(wkhtmltopdf.OrientationPortrait)
	pdfg.Grayscale.Set(true)

	// Create a new input page from a HTML file
	page := wkhtmltopdf.NewPage(htmlFilePath)

	// Set options for this page
	//page.FooterRight.Set("[page]")
	//page.FooterFontSize.Set(10)
	//page.Zoom.Set(0.95)

	// Add to document
	pdfg.AddPage(page)

	// Create PDF document in internal buffer
	err = pdfg.Create()
	if err != nil {
		log.Fatal(err)
	}

	// Generating a temp file for storing resulting PDF
	tmpPdfFile, err := ioutil.TempFile("", "templator-pdf-document-*.pdf")
	if err != nil {
		log.Fatal(err)
	}

	// Write buffer contents to file on disk
	err = pdfg.WriteFile(tmpPdfFile.Name())
	if err != nil {
		log.Fatal(err)
	}

	return tmpPdfFile.Name()
}


