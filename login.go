package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/mateors/mtool"
)

func cookieCheck(w http.ResponseWriter, r *http.Request) (cMap map[string]string) {
	log_session, err := r.Cookie("log_session")

	if err != nil { //got error/not exist -> create one coockie
		userData := make(map[string]string)
		userData["username"] = ""

		sessionStr := mtool.MapToString(userData)
		sessionValue := mtool.EncodeStr(sessionStr, PassWordEncryptionDecryption)

		c := &http.Cookie{
			Name:   "log_session",
			Value:  sessionValue,
			MaxAge: 86400,
		}
		http.SetCookie(w, c)
	} else { //already exist one log_session coockie
		cText := mtool.DecodeStr(log_session.Value, PassWordEncryptionDecryption)
		cMap = mtool.StringToMap(cText)
	}
	return cMap
}
func login(w http.ResponseWriter, r *http.Request) {
	//managing coockie
	cMap := cookieCheck(w, r)
	sessionUser := cMap["username"]

	if r.Method != "POST" {
		if sessionUser != "" { // if already logged in, then redirecting to homepage
			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else {
			//** process starts: preparing data for sending to frontend **//
			// using struct literal
			data := struct {
				Title    string
				Username string
			}{
				Title:    "Login | MASTER-ACADEMY",
				Username: sessionUser,
			}
			//** process ends: preparing data for sending to frontend **//

			//** process starts: executing template **//
			tmpl, err := template.ParseFiles("template/index.gohtml")
			checkErr(err)
			tmpl, err = tmpl.ParseFiles("wpage/login.gohtml")
			checkErr(err)
			tmpl.Execute(w, data)
			//** process ends: executing template **//
		}
	} else {
		//** process starts: getting form data **//
		username := strings.TrimSpace(r.FormValue("username"))
		password := r.FormValue("password")
		//** process ends: getting form data **//

		check := doLogin(username, password) // calling a function which will try to login

		if check == "Done" { // if login done
			//** process ends: storing session values **//
			userData := make(map[string]string)
			userData["username"] = username

			sessionStr := mtool.MapToString(userData)
			sessionValue := mtool.EncodeStr(sessionStr, PassWordEncryptionDecryption)

			c := &http.Cookie{
				Name:   "log_session",
				Value:  sessionValue,
				MaxAge: 86400,
			}

			http.SetCookie(w, c)
			//** process ends: storing session values **//

			//checking access type
			accType := getAccType(username)

			if accType == "ADMIN" {
				fmt.Fprintln(w, "Login ADMIN")
			} else {
				fmt.Fprintln(w, "Login Done")
			}
		} else if check == "username" {
			fmt.Fprintln(w, "Username not found!")
		} else {
			fmt.Fprintln(w, "Wrong password!")
		}
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	//** process starts: deleting session values **//
	userData := make(map[string]string)
	userData["username"] = ""

	sessionStr := mtool.MapToString(userData)
	sessionValue := mtool.EncodeStr(sessionStr, PassWordEncryptionDecryption)

	c := &http.Cookie{
		Name:   "log_session",
		Value:  sessionValue,
		MaxAge: -1,
	}

	http.SetCookie(w, c)
	//** process ends: deleting session values **//

	http.Redirect(w, r, "/", http.StatusSeeOther) //after logout, redirecting to homepage
}

func forgotpassword(w http.ResponseWriter, r *http.Request) {
	//managing coockie
	cMap := cookieCheck(w, r)
	sessionUser := cMap["username"]
	//preparing data for sending frontend
	data := struct {
		Title    string
		Username string
	}{
		Title:    "Request for password reset | MASTER-ACADEMY",
		Username: sessionUser,
	}

	tmpl, err := template.ParseFiles("template/index.gohtml")
	checkErr(err)
	tmpl, err = tmpl.ParseFiles("wpage/forgot_password.gohtml")
	checkErr(err)
	tmpl.Execute(w, data)
}
