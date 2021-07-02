package main

import (
	"html/template"
	"net/http"
	"strings"
)

func studentList(w http.ResponseWriter, r *http.Request) {
	//managing coockie
	cMap := cookieCheck(w, r)
	sessionUser := cMap["username"]

	//checking access type
	accType := getAccType(sessionUser)

	if sessionUser != "" && accType == "ADMIN" {
		//** process starts: preparing data for sending to frontend **//
		// using struct literal
		data := struct {
			Title      string
			IsLoggedIn bool
			Username   string
			Lists      []map[string]interface{}
		}{
			Title:    "Student List | MASTER-ACADEMY",
			Username: sessionUser,
			Lists:    getStudentList(),
		}
		//** process ends: preparing data for sending to frontend **//

		//** process starts: executing template **//
		tmpl, err := template.ParseFiles("template/index.gohtml")
		checkErr(err)
		tmpl, err = tmpl.ParseFiles("wpage/student_list.gohtml")
		checkErr(err)
		tmpl.Execute(w, data)
		//** process ends: executing template **//
	} else if sessionUser != "" && accType == "STUDENT" {
		http.Redirect(w, r, "/error", http.StatusSeeOther)
	} else { //if not logged in then redirecting to login page
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func studentSearch(w http.ResponseWriter, r *http.Request) {
	//managing coockie
	cMap := cookieCheck(w, r)
	sessionUser := cMap["username"]

	//checking access type
	accType := getAccType(sessionUser)

	if sessionUser != "" && accType == "ADMIN" {
		//** process starts: getting form data **//
		formData := make(map[string]string)
		formData["aid"] = strings.TrimSpace(r.FormValue("student-id"))
		formData["account_name"] = strings.ToLower(strings.TrimSpace(r.FormValue("name")))
		formData["email"] = strings.TrimSpace(r.FormValue("email"))
		formData["mobile"] = strings.TrimSpace(r.FormValue("phone"))
		//** process ends: getting form data **//

		//** process starts: preparing data for sending to frontend **//
		// using struct literal
		data := struct {
			Title      string
			IsLoggedIn bool
			Username   string
			Lists      []map[string]interface{}
		}{
			Title:    "Student List | MASTER-ACADEMY",
			Username: sessionUser,
			Lists:    doSearch(formData),
		}
		//** process ends: preparing data for sending to frontend **//

		//** process starts: executing template **//
		tmpl, err := template.ParseFiles("template/index.gohtml")
		checkErr(err)
		tmpl, err = tmpl.ParseFiles("wpage/student_list.gohtml")
		checkErr(err)
		tmpl.Execute(w, data)
		//** process ends: executing template **//
	} else {
		http.Redirect(w, r, "/error", http.StatusSeeOther)
	}
}
