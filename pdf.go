package main

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

func generatePDF(data pdfData) bool {
	//wkhtmltopdf.SetPath(`C:\Program Files\wkhtmltopdf\bin\wkhtmltopdf.exe`)

	//preparing path for store pdf
	wDir, _ := os.Getwd()
	certificatePath := filepath.Join(wDir, "print", "certificate", fmt.Sprintf("%s.pdf", data.LedgerCode))

	//html template path of certificate
	templateFile := "template/certificate/index.html"

	htmlBody := parseTemplate(templateFile, data)
	res := convertToPDF(htmlBody, certificatePath)

	return res
}

func parseTemplate(templateFile string, data pdfData) string {
	tmpl, err := template.ParseFiles(templateFile)
	checkErr(err)

	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, data)
	checkErr(err)

	body := buf.String()
	return body
}

func convertToPDF(htmlBody, certificatePath string) bool {
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	checkErr(err)

	//setting up pdf page
	pdfg.AddPage(wkhtmltopdf.NewPageReader(strings.NewReader(htmlBody)))
	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)
	pdfg.Dpi.Set(300)
	pdfg.Orientation.Set(wkhtmltopdf.OrientationLandscape)
	pdfg.MarginTop.Set(0)
	pdfg.MarginBottom.Set(0)
	pdfg.MarginLeft.Set(0)
	pdfg.MarginRight.Set(0)

	// Create PDF document in internal buffer
	err = pdfg.Create()
	checkErr(err)

	// writing pdf document to a file
	err = pdfg.WriteFile(certificatePath)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}
