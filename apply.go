package main

import (
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

//preparing data
type pdfData struct {
	StudentsName string
	CourseName   string
	Duration     string
	QRImagePath  template.URL
	QRText       string
	LedgerCode   string
}

func apply(w http.ResponseWriter, r *http.Request) {
	//managing coockie
	cMap := cookieCheck(w, r)
	sessionUser := cMap["username"]

	if r.Method != "POST" {
		//** process starts: preparing data for sending to frontend **//
		// using struct literal
		data := struct {
			Title    string
			Username string
		}{
			Title:    "Application | MASTER-ACADEMY",
			Username: sessionUser,
		}
		//** process ends: preparing data for sending to frontend **//
		//** process starts: executing template **//
		tmpl, err := template.ParseFiles("template/index.gohtml")
		checkErr(err)
		tmpl, err = tmpl.ParseFiles("wpage/apply.gohtml")
		checkErr(err)
		tmpl.Execute(w, data)
		//** process ends: executing template **//
	} else {
		studentsName := strings.TrimSpace(r.FormValue("studentsName"))
		email := strings.TrimSpace(r.FormValue("email"))
		courseName := strings.TrimSpace(r.FormValue("courseName"))
		duration := strings.TrimSpace(r.FormValue("duration"))
		codewarsUsername := strings.TrimSpace(r.FormValue("codewarsUsername"))

		// studentsName := "Nahid Hasan"
		// email := "mnh.nahid35@gmail.com"
		// courseName := "Golang Coding Bootcamp"
		// duration := "January - March, 2021"
		// codewarsUsername := "nahidhasan98"

		//checking if already requested or not
		exist := isAlreadyRequested(email, courseName, duration)
		//fmt.Println(exist)
		var flag int = 0

		if !exist {
			ledgerCode := getLedgerCode(email)

			if codewars(codewarsUsername) > 0 && attendence(email) > 0 && ledgerCode != "Not Found" {
				randomLedger := makeRandomString(16)
				domainName := "http://localhost:9001"
				qrLink := fmt.Sprintf("%s/cert/%s", domainName, randomLedger)
				//qrLink = domainName + "/cert/" + randomLedger

				qrRes := GenerateQrCode(qrLink, ledgerCode)
				var accID string

				if qrRes {
					accID = updateWithQR(email, randomLedger)
				}

				if accID != "failed" {
					res := insertCert(accID, ledgerCode, randomLedger, email, studentsName, courseName, duration)

					if res != "failed" {
						flag = 1
					}
				}
			}
		}

		//** process starts: preparing data for sending to frontend **//
		// using struct literal
		data := struct {
			Title      string
			IsLoggedIn bool
			Username   string
			Message    string
		}{
			Title:    "Certificate Generator | MASTER-ACADEMY",
			Username: sessionUser,
			Message:  "",
		}
		//** process ends: preparing data for sending to frontend **//

		if exist { //already have a request
			data.Message = "You have already a pending request. Please wait for the mail."
		} else if flag == 1 { //if success
			data.Message = "Thank you for requesting Cetificate. Your application status is pending now. We will sent a mail attached with your certificate when its ready."
		} else {
			data.Message = "Sorry, your request cannot be completed right now. You have either low attendance or low codewars score."
		}
		fmt.Fprintln(w, data.Message)

		// //** process starts: executing template **//
		// tmpl, err := template.ParseFiles("template/index.gohtml")
		// checkErr(err)
		// tmpl.Execute(w, data)
		// //** process ends: executing template **//
	}
}

func makeRandomString(length int) string {
	randomLedger := ""

	//determining a seed value
	rand.Seed(time.Now().UnixNano())

	//setting up a character set
	charSet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	//making random string of given length
	for i := 0; i < length; i++ {
		random := rand.Intn(len(charSet))
		randomChar := charSet[random]
		randomLedger += string(randomChar)
	}
	//fmt.Println(randomLedger)

	return randomLedger
}
