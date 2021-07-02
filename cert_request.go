package main

import (
	"html/template"
	"net/http"
)

func certRequestList(w http.ResponseWriter, r *http.Request) {
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
			Lists:    getCertRequestList(),
		}
		//** process ends: preparing data for sending to frontend **//
		//fmt.Println(getCertRequestList("company::1"))

		//** process starts: executing template **//
		tmpl, err := template.ParseFiles("template/index.gohtml")
		checkErr(err)
		tmpl, err = tmpl.ParseFiles("wpage/cert_req_list.gohtml")
		checkErr(err)
		tmpl.Execute(w, data)
		//** process ends: executing template **//
	} else if sessionUser != "" && accType != "STUDENT" {
		http.Redirect(w, r, "/error", http.StatusSeeOther)
	} else { //if not logged in then redirecting to login page
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
