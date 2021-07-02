package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
)

func cert(w http.ResponseWriter, r *http.Request) {
	vMap := mux.Vars(r)
	certCode := vMap["certCode"]

	certInfo := getCertInfo(certCode)

	if len(certInfo) == 0 {
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}

	//managing coockie
	cMap := cookieCheck(w, r)
	sessionUser := cMap["username"]

	//** process starts: preparing data for sending to frontend **//
	// using struct literal
	data := struct {
		Title        string
		Username     string
		Message      string
		StudentsName string
		CourseName   string
		Duration     string
	}{
		Title:        "Certificate Generator | MASTER-ACADEMY",
		Username:     sessionUser,
		StudentsName: fmt.Sprintf("%s", certInfo[0]["students_name"]),
		CourseName:   fmt.Sprintf("%s", certInfo[0]["course_name"]),
		Duration:     fmt.Sprintf("%s", certInfo[0]["duration"]),
	}

	if len(certInfo) == 0 {
		data.Message = "Couldn't find any info about this certificate."
	} else {
		data.Message = "Found in our database."
		fmt.Fprintln(w, data.StudentsName, data.CourseName, data.Duration)
	}
	//** process ends: preparing data for sending to frontend **//

	fmt.Fprintln(w, data.Message)

	// //** process starts: executing template **//
	// tmpl, err := template.ParseFiles("template/index.gohtml")
	// checkErr(err)
	// tmpl.Execute(w, data)
	// //** process ends: executing template **//
}

func view(w http.ResponseWriter, r *http.Request) {
	vMap := mux.Vars(r)
	certID := vMap["certID"]

	//managing coockie
	cMap := cookieCheck(w, r)
	sessionUser := cMap["username"]

	//checking access type
	accType := getAccType(sessionUser)

	if sessionUser != "" && accType == "ADMIN" {
		pdfFileName := fmt.Sprintf("%s.pdf", certID)

		wDir, _ := os.Getwd()
		pdfFilePath := filepath.Join(wDir, "print", "certificate", pdfFileName)
		_, err := os.Stat(pdfFilePath)

		if err == nil { // if file exist
			src := fmt.Sprintf("print/certificate/%s", pdfFileName)
			dst := fmt.Sprintf("assets/temp/%s", pdfFileName)
			copyFile(src, dst)

			tempPath := fmt.Sprintf("/resources/temp/%s", pdfFileName)

			// using struct literal
			data := struct {
				Title       string
				PDFFilePath template.URL
			}{
				Title:       "Certificate Dispaly | MASTER-ACADEMY",
				PDFFilePath: template.URL(tempPath),
			}

			//** process starts: executing template **//
			tmpl, err := template.ParseFiles("wpage/view_cert.gohtml")
			checkErr(err)
			tmpl.Execute(w, data)
			//** process ends: executing template **//
		} else {
			http.Redirect(w, r, "/error", http.StatusSeeOther)
		}
	} else {
		http.Redirect(w, r, "/error", http.StatusSeeOther)
	}
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	checkErr(err)
	defer in.Close()

	out, err := os.Create(dst)
	checkErr(err)
	defer out.Close()

	_, err = io.Copy(out, in)
	checkErr(err)
	return out.Close()
}

// defer func() {
// 	if req := recover(); req != nil {
// 		fmt.Println("cert panicking with value >>", req)
// 	}
// }()
