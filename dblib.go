package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	models "github.com/mateors/master-erp"
	"github.com/mateors/mcb"
	"github.com/mateors/mtool"
	"golang.org/x/crypto/bcrypt"
)

var db *mcb.DB

const (
	BucketName = "master_erp" //
	cid        = "company::1" //company id

	AccessTable             = "access"               //
	AccountTable            = "account"              //
	AccountHeadTable        = "achead"               //account_head ledger head
	ActivityLogTable        = "activity_log"         //
	AddressTable            = "address"              // contact address
	CategoryTable           = "category"             //**ItemCategory
	CompanyTable            = "company"              //
	ContactTable            = "contact"              //
	DepartmentTable         = "department"           //owner = item | item_line | office
	DeviceLogTable          = "device_log"           //
	DocKeeperTable          = "doc_keeper"           //
	DocPayShipTable         = "docpayship"           //doc_payship_info
	FileStoreTable          = "file_store"           //
	ItemTable               = "item"                 //**
	ItemAttributeTable      = "item_attribute"       //
	ItemAttributeValueTable = "item_attribute_value" //
	ItemRewardTable         = "item_reward"          // item purchase point for loyalty
	LedgerTransactionTable  = "ledger_transaction"   //
	LoginTable              = "login"                //
	LoginSessionTable       = "login_session"        //
	MessageTable            = "message"              //
	PaymentOptionTable      = "payment_option"       //
	Rateplan                = "rateplan"             //
	SettingsTable           = "settings"             //
	ShippingAddressTable    = "shipping_address"     //*
	ShippingOptionTable     = "shipping_option"      //
	StockMovementTable      = "stock_movement"       //
	TaxTable                = "tax"                  //tax | vat table
	TransactionRecordTable  = "transaction_record"   //sales|purchase
	UOMTable                = "uom"                  //
	VerificationTable       = "verification"         //
	VisitorSessionTable     = "visitor_session"      //
	WarehouseTable          = "warehouse"            //
	//AccountGroupTable       = "account_group"        //removed account group
)

// init() function establish Database connection
func init() {
	//couchbase connection block
	db = mcb.Connect("localhost", "root", "bootcamp")
	res, err := db.Ping()
	if err != nil {
		fmt.Println(res)
		os.Exit(1)
	}
}

func doRegistration(formData map[string]string, w http.ResponseWriter, r *http.Request) {
	//** process starts: checking into DB **//
	check := isExsist(formData) //checking if username or email already exist or not
	if check != "" {
		fmt.Fprintln(w, check)
		return
	}
	//** process ends: checking into DB **//

	//** process starts: AccountTable **//
	accID, accSerial := findAailabledocID(AccountTable)
	ledgerCode := makeLedgerCode(accSerial)
	fullName := fmt.Sprintf("%s %s", strings.TrimSpace(formData["firstName"]), strings.TrimSpace(formData["lastName"]))
	createDate := time.Now().String()[:19]

	r.Form.Set("bucket", BucketName)                                   //
	r.Form.Set("aid", accID)                                           //
	r.Form.Set("type", AccountTable)                                   // account
	r.Form.Set("cid", cid)                                             // foreign key
	r.Form.Set("serial", strconv.FormatInt(accSerial, 10))             // company wise increase
	r.Form.Set("photo", "")                                            // account owner photo
	r.Form.Set("account_type", "STUDENT")                              // vendor,goods_supplier,customer,consumer,payment_provider,shipping_provider
	r.Form.Set("account_name", fullName)                               // supplier business name or username
	r.Form.Set("customid", "")                                         // unique customer IDENTIFICATION
	r.Form.Set("code", ledgerCode)                                     // supplier or customer code or ledgercode
	r.Form.Set("login_id", "")                                         // foreign key
	r.Form.Set("first_name", strings.TrimSpace(formData["firstName"])) //
	r.Form.Set("last_name", strings.TrimSpace(formData["lastName"]))   //
	r.Form.Set("dob", "")                                              //
	r.Form.Set("gender", "")                                           // female,male,other
	r.Form.Set("mobile", "")                                           // phone
	r.Form.Set("email", formData["email"])                             //
	r.Form.Set("remarks", formData["remarks"])                         // avg salary
	r.Form.Set("create_date", createDate)                              //
	r.Form.Set("update_date", createDate)                              //
	r.Form.Set("status", "1")                                          //

	var account models.Account
	insert1 := db.Insert(r.Form, &account) //insert query
	//** process ends: AccountTable **//

	//** process starts: LoginTable **//
	loginID, loginSerial := findAailabledocID(LoginTable)

	r.Form.Set("aid", loginID)                               //
	r.Form.Set("type", LoginTable)                           //
	r.Form.Set("cid", cid)                                   // foreign key
	r.Form.Set("serial", strconv.FormatInt(loginSerial, 10)) // company wise increasing
	r.Form.Set("account_id", accID)                          // foreign key
	r.Form.Set("access_name", "STUDENT")                     // customer type
	r.Form.Set("username", formData["username"])             // email or mobile as username
	r.Form.Set("passw", makeHash(formData["password"]))      //
	r.Form.Set("create_date", createDate)                    //
	r.Form.Set("last_login", "")                             //
	r.Form.Set("status", "1")                                //

	var login models.Login
	insert2 := db.Insert(r.Form, &login)
	//** process ends: LoginTable **//

	//** process starts: AddressTable **//
	addID, addSerial := findAailabledocID(AddressTable)

	r.Form.Set("aid", addID)                               //
	r.Form.Set("type", AddressTable)                       //
	r.Form.Set("cid", cid)                                 // foreign key
	r.Form.Set("serial", strconv.FormatInt(addSerial, 10)) //
	r.Form.Set("account_id", accID)                        // foreign key
	r.Form.Set("address_type", "Billing")                  // billing,shipping
	r.Form.Set("country", "")                              //
	r.Form.Set("state", "")                                //
	r.Form.Set("city", "")                                 //
	r.Form.Set("address1", "")                             //
	r.Form.Set("address2", "")                             //
	r.Form.Set("zip", "")                                  //
	r.Form.Set("status", "1")                              //

	var address models.Address
	insert3 := db.Insert(r.Form, &address)
	//** process ends: AddressTable **//

	//** process starts: AccountHeadTable **//
	fVal := getFieldVals("aid,account_type", AccountHeadTable, fmt.Sprintf(`cid="%s" AND name="%s" AND account_group="group"`, cid, "Account Receivable"))
	parentID := fVal["aid"]
	accHeadID, accHeadSerial := findAailabledocID(AccountHeadTable)

	r.Form.Set("aid", accHeadID)                               // unique id
	r.Form.Set("type", AccountHeadTable)                       // table
	r.Form.Set("cid", cid)                                     // foreign key
	r.Form.Set("serial", strconv.FormatInt(accHeadSerial, 10)) // company wise increase
	r.Form.Set("account_group", fVal["account_group"])         // AccountGroup= Asset|Liability|Equity|Revenue|Expense
	r.Form.Set("account_type", "STUDENT")                      // group|head
	r.Form.Set("name", fmt.Sprintf("%s (STUDENT)", fullName))  // ledger name
	r.Form.Set("description", "")                              // ledger_details
	r.Form.Set("identifier", "receivable")                     // for ensuring no ledgers are duplicate
	r.Form.Set("code", ledgerCode)                             // ledger number or ledgercode
	r.Form.Set("parent_id", parentID)                          // parent account
	r.Form.Set("balance", fmt.Sprintf("%f", 0.0))              // ledger balance
	r.Form.Set("baltype", "Eq")                                // ledger balance type Dr or Cr
	r.Form.Set("restricted", "0")                              // 1=Yes, No=0
	r.Form.Set("cost_center", "0")                             // 1=Yes, No=0
	r.Form.Set("remarks", "")                                  //
	r.Form.Set("create_date", createDate)                      // insert date
	r.Form.Set("status", "1")                                  // 0=Inactive, 1=Active, 9=Deleted

	var accountHead models.AccountHead
	insert4 := db.Insert(r.Form, &accountHead)
	//** process ends: AccountHeadTable **//

	//** process starts: updating AccountTable with login ID **//
	qs := `UPDATE master_erp SET login_id = "%s" WHERE type = "account" AND aid = "%s" RETURNING login_id`
	sql := fmt.Sprintf(qs, loginID, accID)

	pRes := db.Query(sql)
	rows := pRes.GetRows()
	updateRes := rows[0]["login_id"].(string)
	//** process ends: updating AccountTable **//

	if insert1.Status == "success" && insert2.Status == "success" && insert3.Status == "success" && insert4.Status == "success" && updateRes == loginID {
		//sending mail verification link to the user mail
		//sendMail(r.FormValue("email"), r.FormValue("username"), "http://localhost:9000/verifyemail/"+token)

		fmt.Fprintln(w, "Registration Done")
	} else {
		fmt.Fprintln(w, "Registration Error")
	}
}

// isExist() function checks the existance of username and email in DB
func isExsist(formData map[string]string) string {
	//checking: username
	qs := `SELECT username FROM master_erp WHERE type="login" AND username="%s"`
	sql := fmt.Sprintf(qs, formData["username"])

	pRes := db.Query(sql)
	rows := pRes.GetRows()
	if len(rows) != 0 {
		return "username"
	}

	//checking: email
	qs = `SELECT email FROM master_erp WHERE type="account" AND email="%s"`
	sql = fmt.Sprintf(qs, formData["email"])

	pRes = db.Query(sql)
	rows = pRes.GetRows()
	if len(rows) != 0 {
		return "email"
	}
	return ""
}

// getLedgerCode() function creates a ledger code
func makeLedgerCode(accSerial int64) string {
	serial := strconv.FormatInt(accSerial, 10)
	ledgerCode := "12"
	for i := 0; i < 8-len(serial); i++ {
		ledgerCode += "0"
	}
	ledgerCode += serial
	return ledgerCode
}

// makeHash() function encrypts the password
func makeHash(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes)
}

// findAailabledocID() function finds the next serial number and aid for any table [this function Got from Mostain Sir]
func findAailabledocID(table string) (string, int64) {
	var docID string
	var scount int64 = int64(maxDoc(BucketName, table, db))
	i := scount + 1

	for {
		count := i
		docID = fmt.Sprintf(`%s::%v`, table, count)
		//sql := fmt.Sprintf(`SELECT COUNT(*)as cnt FROM %s WHERE META().id="%s";`, BucketName, docID)
		sql := fmt.Sprintf(`SELECT META(d).id, true AS docexist FROM %s AS d USE KEYS ["%s"];`, BucketName, docID)
		//fmt.Println(sql)
		pRes := db.Query(sql)
		rows := pRes.GetRows()

		if len(rows) == 0 {
			return docID, count
		}
		exist := rows[0]["docexist"].(bool)
		fmt.Println("ALREADY EXIST", i, exist, sql)
		i++
	}
}

// maxDoc() function is called from findAailabledocID() function [this function Got from Mostain Sir]
func maxDoc(bucketName, tableName string, db *mcb.DB) (count float64) {
	sql := fmt.Sprintf(`SELECT NVL(MAX(serial),0)as cnt FROM %s WHERE type="%s";`, bucketName, tableName)
	pResponse := db.Query(sql)
	//fmt.Println("CountDoc>", sql, pResponse.Result)
	for _, val := range pResponse.Result {
		//fmt.Printf("%v %T\n", val, val)
		//rMap:=val.(map)
		cMap := val.(map[string]interface{})
		if c, ok := cMap["cnt"]; ok {
			count = c.(float64)
		}
	}
	//fmt.Println("CountTable>", pResponse.Result)
	return
}

// getFieldVals() function returns the field values [this function Got from Mostain Sir]
func getFieldVals(fieldNames, table, where string) (fVals map[string]string) {
	fields := strings.Split(fieldNames, ",")
	fVals = make(map[string]string)

	sql := fmt.Sprintf(`SELECT %v FROM %v WHERE type="%s" AND %v;`, fieldNames, BucketName, table, where)
	pRes := db.Query(sql)
	//fmt.Println("Field >>", sql)

	rows := pRes.GetRows()
	for _, row := range rows {

		for key := range row {

			if mtool.ArrayValueExist(fields, key) {
				fieldValue := fmt.Sprintf("%v", row[key])
				fVals[key] = fieldValue
			}
		}
	}
	return
}

// doLogin() function tried to login to user account
func doLogin(username, password string) string {
	qs := `SELECT passw FROM master_erp WHERE type="login" AND username="%s"`
	sql := fmt.Sprintf(qs, username)

	pRes := db.Query(sql)
	rows := pRes.GetRows()
	if len(rows) == 0 {
		return "username" // username not found
	}
	originalPassword := rows[0]["passw"].(string) //password saved in DB
	if checkHash(password, originalPassword) {    //checking if original password is matched with provided password or not
		return "Done" // login successful
	}
	return ""
}

// checkHash() function checks the provided password with original hashed password
func checkHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// getStudentList() function returns the registered student list with detail information
func getStudentList() []map[string]interface{} {
	qs := `SELECT ac.first_name, ac.last_name, ac.email, ac.mobile, ac.create_date, ac.status, l.username as username FROM master_erp ac
	LEFT JOIN master_erp as a ON a.account_id=META(ac).id AND a.type="address"
	LEFT JOIN master_erp as l ON l.account_id=META(ac).id AND l.type="login"
	WHERE ac.cid = "%s" AND ac.type = "account" AND ac.account_type = "STUDENT" AND ac.status IN [0,1] ORDER BY ac.create_date DESC;`
	sql := fmt.Sprintf(qs, cid)

	pRes := db.Query(sql)
	rows := pRes.GetRows()
	//fmt.Println(rows)

	return rows
}

// getSingleStudent() function returns the information of a single students
func getSingleStudent(username string) []map[string]interface{} {
	qs := `SELECT ac.first_name, ac.last_name, ac.email, ac.mobile, ac.status, l.username as username, ad.city as city FROM master_erp ac
	LEFT JOIN master_erp as a ON a.account_id=META(ac).id AND a.type="address"
	LEFT JOIN master_erp as l ON l.account_id=META(ac).id AND l.type="login"
	LEFT JOIN master_erp as ad ON ad.account_id=META(ac).id AND ad.type="address"
	WHERE l.username = "%s" AND ac.type = "account" AND ac.status IN [0,1];`
	sql := fmt.Sprintf(qs, username)

	pRes := db.Query(sql)
	rows := pRes.GetRows()

	return rows
}

// updateStudentInfo() function updates the provided field information of a single student
func updateStudentInfo(firstName, lastName, username, mobile, city string, status int) string {
	// getting account_id of the user by provided username
	qs := `SELECT account_id FROM master_erp WHERE type = "login" AND username = "%s"`
	sql := fmt.Sprintf(qs, username)

	pRes := db.Query(sql)
	rows := pRes.GetRows()
	accID := rows[0]["account_id"].(string)

	// updating account table
	currentDateTime := time.Now().String()[:19]
	qs = `UPDATE master_erp SET first_name = "%s", last_name = "%s", mobile = "%s", status = %d, update_date = "%s" WHERE type = "account" AND aid = "%s" RETURNING aid`
	sql = fmt.Sprintf(qs, firstName, lastName, mobile, status, currentDateTime, accID)

	pRes = db.Query(sql)
	rows = pRes.GetRows()
	update1 := rows[0]["aid"].(string)

	// updating address table
	qs = `UPDATE master_erp SET city = "%s" WHERE type = "address" AND account_id = "%s" RETURNING account_id`
	sql = fmt.Sprintf(qs, city, accID)

	pRes = db.Query(sql)
	rows = pRes.GetRows()
	update2 := rows[0]["account_id"].(string)

	if update1 == accID && update2 == accID { //if update success
		return "OK"
	}

	return "failed"
}

func doSearch(formData map[string]string) []map[string]interface{} {
	qs := `SELECT ac.first_name, ac.last_name, ac.email, ac.mobile, ac.create_date, ac.status, l.username as username FROM master_erp ac
	LEFT JOIN master_erp as a ON a.account_id=META(ac).id AND a.type="address"
	LEFT JOIN master_erp as l ON l.account_id=META(ac).id AND l.type="login"
	WHERE ac.cid = "` + cid + `" AND ac.type = "account"  AND ac.account_type = "STUDENT" AND ac.status IN [0,1]`
	for key, val := range formData {
		if val != "" {
			if key == "account_name" {
				qs += ` AND LOWER(ac.` + key + `) LIKE "%` + val + `%"`
			} else {
				qs += ` AND ac.` + key + `="` + val + `"`
			}
		}
	}
	qs += ` ORDER BY ac.create_date DESC;`

	pRes := db.Query(qs)
	rows := pRes.GetRows()
	//fmt.Println(rows)

	return rows
}

func getLedgerCode(email string) string {
	qs := `SELECT code FROM master_erp WHERE type="account" AND email="%s"`
	sql := fmt.Sprintf(qs, email)

	pRes := db.Query(sql)
	rows := pRes.GetRows()
	if len(rows) == 0 {
		return "Not Found"
	}

	return rows[0]["code"].(string)
}

func isAlreadyRequested(email, courseName, duration string) bool {
	qs := `SELECT email FROM master_erp WHERE email="%s" AND course_name ="%s" AND duration = "%s" AND type = "certificate" AND cid = "%s";`
	sql := fmt.Sprintf(qs, email, courseName, duration, cid)

	pRes := db.Query(sql)
	rows := pRes.GetRows()

	return len(rows) != 0
}

func updateWithQR(email, randomLedger string) string {
	qs := `UPDATE master_erp SET customid = "%s" WHERE type = "account" AND email = "%s" RETURNING aid, code`
	sql := fmt.Sprintf(qs, randomLedger, email)

	pRes := db.Query(sql)
	rows := pRes.GetRows()

	if len(rows) == 0 {
		return "failed"
	}
	return rows[0]["aid"].(string)
}

func insertCert(accID, ledgerCode, randomLedger, email, studentsName, courseName, duration string) string {
	cerID, cerSerial := findAailabledocID("certificate")
	requestDate := time.Now().String()[:19]

	qs := `INSERT INTO master_erp ( KEY, VALUE ) VALUES ( "%s", 
	{ "aid": "%s", "serial": %s, "type": "certificate", "cid": "%s",  "account_id":"%s", "email":"%s",
	"students_name":"%s","course_name":"%s","duration":"%s","delivery_status":"Pending", "request_date":"%s", "code":"%s", "customid":"%s"})
	RETURNING aid;`
	sql := fmt.Sprintf(qs, cerID, cerID, strconv.FormatInt(cerSerial, 10), cid, accID, email, studentsName, courseName, duration, requestDate, ledgerCode, randomLedger)

	pRes := db.Query(sql)
	rows := pRes.GetRows()

	if len(rows) == 0 {
		return "failed"
	}
	return rows[0]["aid"].(string)
}

// getSingleStudent() function returns the information of a single students
func getCertInfo(certCode string) []map[string]interface{} {
	qs := `SELECT ac.code, ac.customid, cer.aid, cer.account_id, cer.students_name, cer.email, cer.course_name, cer.duration, cer.request_date, cer.delivery_status FROM master_erp ac
	LEFT JOIN master_erp as cer ON cer.account_id=META(ac).id AND cer.type="certificate"
	WHERE ac.type = "account" AND ac.customid = "%s" AND ac.cid = "%s";`
	sql := fmt.Sprintf(qs, certCode, cid)

	pRes := db.Query(sql)
	rows := pRes.GetRows()

	return rows
}

// getStudentList() function returns the registered student list with detail information
func getCertRequestList() []map[string]interface{} {
	qs := `SELECT aid, account_id, students_name, email, course_name, duration, request_date, delivery_status FROM master_erp
	WHERE type = "certificate" AND cid = "%s" ORDER BY request_date DESC;`
	sql := fmt.Sprintf(qs, cid)

	pRes := db.Query(sql)
	rows := pRes.GetRows()
	//fmt.Println(rows)

	return rows
}
func getCertReqInfo(email string) []map[string]interface{} {
	qs := `SELECT ac.code, ac.customid, cer.aid, cer.account_id, cer.students_name, cer.email, cer.course_name, cer.duration, cer.request_date, cer.delivery_status FROM master_erp ac
	LEFT JOIN master_erp as cer ON cer.account_id=META(ac).id AND cer.type="certificate"
	WHERE ac.type = "account" AND cer.email = "%s" AND ac.cid = "%s";`
	sql := fmt.Sprintf(qs, email, cid)

	pRes := db.Query(sql)
	rows := pRes.GetRows()
	//fmt.Println(rows)

	return rows
}

func getAccType(username string) string {
	//checking: access name
	qs := `SELECT access_name FROM master_erp WHERE type="login" AND username="%s" AND cid="%s";`
	sql := fmt.Sprintf(qs, username, cid)

	pRes := db.Query(sql)
	rows := pRes.GetRows()
	if len(rows) == 0 {
		return ""
	}
	return rows[0]["access_name"].(string)
}

func getTotalStudentsNumber() int {
	qs := `SELECT aid FROM master_erp
	WHERE cid = "%s" AND type = "account" AND account_type = "STUDENT" AND status IN [0,1] ORDER BY create_date DESC;`
	sql := fmt.Sprintf(qs, cid)

	pRes := db.Query(sql)
	rows := pRes.GetRows()
	//fmt.Println(rows)

	return len(rows)
}
func getTotalCertReqNumber() (int, int) {
	qs := `SELECT delivery_status FROM master_erp
	WHERE type = "certificate" AND cid = "%s" ORDER BY request_date DESC;`
	sql := fmt.Sprintf(qs, cid)

	pRes := db.Query(sql)
	rows := pRes.GetRows()
	//fmt.Println(rows)

	pending := 0
	for _, val := range rows {
		if val["delivery_status"] == "Pending" {
			pending++
		}
	}

	return pending, len(rows) - pending
}

func updateDeliveryStatus(email string) string {
	qs := `UPDATE master_erp SET delivery_status = "Delivered" WHERE type = "certificate" AND email = "%s" RETURNING aid`
	sql := fmt.Sprintf(qs, email)

	pRes := db.Query(sql)
	rows := pRes.GetRows()

	if len(rows) == 0 {
		return "failed"
	}
	return rows[0]["aid"].(string)
}
