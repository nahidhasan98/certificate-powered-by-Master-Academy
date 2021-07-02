package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// GenerateQrCode function takes a string
func GenerateQrCode(value string, fileName string) bool {
	apiURL := fmt.Sprintf("https://chart.googleapis.com/chart?chs=150x150&cht=qr&chl=%s", value)

	//setting up new request
	req, err := http.NewRequest("GET", apiURL, nil)
	checkErr(err)

	// executing request & receiving response
	client := &http.Client{}
	res, err := client.Do(req)
	checkErr(err)
	defer res.Body.Close()

	// taking response body
	body, err := ioutil.ReadAll(res.Body)
	checkErr(err)

	// saving qr image to directory
	fullPath := fmt.Sprintf("data/account/%s.png", fileName)
	err = ioutil.WriteFile(fullPath, body, 0644)
	checkErr(err)

	if body != nil && err == nil {
		return true
	}
	return false
}
