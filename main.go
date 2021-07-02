package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

var PassWordEncryptionDecryption = "master_certificate"

func checkErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func main() {
	//Just a message for ensuring the local server is running
	fmt.Println("Local server is listening on port 9001...")

	//(instead of default 'http' router) using Gorilla mux router
	rMux := mux.NewRouter()

	rMux.HandleFunc("/", home)
	rMux.HandleFunc("/register", register)
	rMux.HandleFunc("/login", login)
	rMux.HandleFunc("/logout", logout)
	rMux.HandleFunc("/forgot_password", forgotpassword)
	rMux.HandleFunc("/dashboard", dashboard)
	rMux.HandleFunc("/student_list", studentList)
	rMux.HandleFunc("/student_search", studentSearch)
	rMux.HandleFunc("/about", about)
	rMux.HandleFunc("/contact", contact)

	rMux.HandleFunc("/api/{request}", api)

	rMux.HandleFunc("/cert_req_list", certRequestList)
	rMux.HandleFunc("/apply", apply)
	rMux.HandleFunc("/cert/{certCode}", cert)
	rMux.HandleFunc("/view/{certID}", view)

	rMux.HandleFunc("/error", http.NotFound)

	//serving file from server to client
	rMux.PathPrefix("/resources/").Handler(http.StripPrefix("/resources/", http.FileServer(http.Dir("assets"))))

	//localhost running on port 9001
	http.ListenAndServe(":9001", rMux)
}

func home(w http.ResponseWriter, r *http.Request) {
	//managing coockie
	cMap := cookieCheck(w, r)
	sessionUser := cMap["username"]

	//** process starts: preparing data for sending to frontend **//
	// using struct literal
	data := struct {
		Title    string
		Username string
	}{
		Title:    "Certificate Generator | MASTER-ACADEMY",
		Username: sessionUser,
	}
	//** process ends: preparing data for sending to frontend **//

	//** process starts: executing template **//
	tmpl, err := template.ParseFiles("template/index.gohtml")
	checkErr(err)
	tmpl.Execute(w, data)
	//** process ends: executing template **//
}

func about(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "This is about page")
}

func contact(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "This is contact page")
}
