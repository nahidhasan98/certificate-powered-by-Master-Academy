package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func api(w http.ResponseWriter, r *http.Request) {
	vMap := mux.Vars(r)
	request := vMap["request"]
	//fmt.Println(request)

	//managing coockie
	cMap := cookieCheck(w, r)
	sessionUser := cMap["username"]

	//checking access type
	accType := getAccType(sessionUser)

	if sessionUser != "" && accType == "ADMIN" {
		if request == "update" {
			firstName := strings.TrimSpace(r.FormValue("firstName"))
			lastName := strings.TrimSpace(r.FormValue("lastName"))
			username := strings.TrimSpace(r.FormValue("username"))
			mobile := strings.TrimSpace(r.FormValue("mobile"))
			city := strings.TrimSpace(r.FormValue("city"))
			status := strings.TrimSpace(r.FormValue("status"))
			sts, _ := strconv.Atoi(status)

			res := updateStudentInfo(firstName, lastName, username, mobile, city, sts)
			fmt.Fprintln(w, res)
			return
		}

		if request == "register" {
			formData := make(map[string]string)
			formData["firstName"] = strings.TrimSpace(r.FormValue("firstName"))
			formData["lastName"] = strings.TrimSpace(r.FormValue("lastName"))
			formData["username"] = strings.TrimSpace(r.FormValue("username"))
			formData["email"] = strings.TrimSpace(r.FormValue("email"))
			formData["password"] = r.FormValue("password")
			formData["remarks"] = strings.TrimSpace(r.FormValue("remarks"))

			doRegistration(formData, w, r)
			return
		}

		if request == "email" {
			sendCertEmail(w, r)
			return
		}

		index := strings.Index(request, "student_info")
		if index != -1 { //url is like this "/api/student_info/....."
			processStudentInfo(request, w)
			return
		}

		index = strings.Index(request, "pdf")
		if index != -1 { //url is like this "/api/student_info/....."
			processPDF(request, w)
			return
		}
	} else {
		http.Redirect(w, r, "/error", http.StatusSeeOther)
	}
	//errorPage(w, http.StatusBadRequest) //http.StatusBadRequest = 400
}

// processStudentInfo() function returns a single student information to frontend
func processStudentInfo(path string, w http.ResponseWriter) {
	username := strings.TrimPrefix(path, "student_info-")
	info := getSingleStudent(username)

	var res []byte
	var err error
	w.Header().Set("Content-Type", "application/json")

	if len(info) > 0 {
		res, err = json.Marshal(info[0])
		checkErr(err)
	}
	w.Write(res)
}

func processPDF(path string, w http.ResponseWriter) {
	email := strings.TrimPrefix(path, "pdf-")
	info := getCertReqInfo(email)
	//fmt.Println(email, info[0]["account_id"])

	pdfFileName := fmt.Sprintf("%s.pdf", info[0]["code"])
	wDir, _ := os.Getwd()
	pdfFilePath := filepath.Join(wDir, "print", "certificate", pdfFileName)
	//fmt.Println(pdfFilePath)

	_, err := os.Stat(pdfFilePath)

	if os.IsNotExist(err) { // path/to/whatever does *not* exist
		//fmt.Println("Not Exsist")

		qrFileName := fmt.Sprintf("%s.png", info[0]["code"])
		filePath := filepath.Join(wDir, "data", "account", qrFileName)
		qrImagePath := template.URL(filePath)

		data := pdfData{
			StudentsName: fmt.Sprintf("%s", info[0]["students_name"]),
			CourseName:   fmt.Sprintf("%s", info[0]["course_name"]),
			Duration:     fmt.Sprintf("%s", info[0]["duration"]),
			QRImagePath:  qrImagePath,
			QRText:       fmt.Sprintf("MA-%s", info[0]["code"]),
			LedgerCode:   fmt.Sprintf("%s", info[0]["code"]),
		}
		generatePDF(data)
	}

	pdfLink := fmt.Sprintf("/view/%s", info[0]["code"])
	w.Write([]byte(pdfLink))
}

func sendCertEmail(w http.ResponseWriter, r *http.Request) {
	//r.ParseForm()
	list := r.FormValue("list")
	list = list[1 : len(list)-1]
	emails := strings.Split(list, ",")

	data := struct {
		Success []string
		Failed  []string
	}{}

	for key := range emails {
		email := emails[key][1 : len(emails[key])-1]
		ledger := getLedgerCode(email)

		//preparing path for store pdf
		wDir, _ := os.Getwd()
		certificatePath := filepath.Join(wDir, "print", "certificate", fmt.Sprintf("%s.pdf", ledger))

		from := "fakenahid@gmail.com"
		to := email
		subject := "Master Academy Certificate"
		body := ""
		attachmentFilePath := certificatePath

		sRes := sendEmail(from, to, subject, body, attachmentFilePath)
		if !sRes {
			data.Failed = append(data.Failed, emails[key])
			continue
		}
		//mail sent

		//updateing certificate table
		qRes := updateDeliveryStatus(email)
		if qRes != "failed" {
			data.Success = append(data.Success, emails[key])
		}
	}

	w.Header().Set("Content-Type", "application/json")
	b, _ := json.Marshal(data)
	w.Write(b)
}
