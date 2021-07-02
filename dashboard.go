package main

import (
	"html/template"
	"net/http"
	"os"
)

func dashboard(w http.ResponseWriter, r *http.Request) {
	//removing temp PDF file created at pdf view page
	os.RemoveAll("assets/temp/")
	os.MkdirAll("assets/temp/", 0755)

	//managing coockie
	cMap := cookieCheck(w, r)
	sessionUser := cMap["username"]

	//checking access type
	accType := getAccType(sessionUser)

	if sessionUser != "" && accType == "ADMIN" {
		//** process starts: preparing data for sending to frontend **//
		// using struct literal
		data := struct {
			Title               string
			Username            string
			TotalStudents       int
			CertRquestPending   int
			CertRquestDelivered int
		}{
			Title:         "Dashboard | MASTER-ACADEMY",
			Username:      sessionUser,
			TotalStudents: getTotalStudentsNumber(),
		}
		pending, deliverd := getTotalCertReqNumber()
		data.CertRquestPending = pending
		data.CertRquestDelivered = deliverd
		//** process ends: preparing data for sending to frontend **//

		//** process starts: executing template **//
		tmpl, err := template.ParseFiles("template/index.gohtml")
		checkErr(err)
		tmpl, err = tmpl.ParseFiles("wpage/dashboard.gohtml")
		checkErr(err)
		tmpl.Execute(w, data)
		//** process ends: executing template **//
	} else if sessionUser != "" && accType != "STUDENT" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else { //if not logged in then redirecting to login page
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
