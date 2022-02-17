package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

//////////////////////////////////////////////////////////////////////////////////////////
//								SMS VERIFICATION TASKS									//
//////////////////////////////////////////////////////////////////////////////////////////
func GetSMSToken() BearerResponse {
	tvURL := "https://www.textverified.com/api/"
	client := http.Client{}
	req, err := http.NewRequest("POST", tvURL+"SimpleAuthentication", nil)
	if err != nil {
		panic(err)
	}

	req.Header = http.Header{
		"X-SIMPLE-API-ACCESS-TOKEN": []string{"1_l8T2o8bPsP252roHDXhtO-tRzX3tqROlPzzam8kTaj7YPRlnccMekpuAmRsDh4r9_H13sjOr"},
	}

	res, err := client.Do(req)
	if err != nil {
		panic(err)
	} else {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}
		//bodyString := string(bodyBytes)
		var bear BearerResponse
		if err := json.Unmarshal(bodyBytes, &bear); err != nil { // Parse []byte to go struct pointer
			fmt.Println("Can not unmarshal JSON")
		}
		return bear
	}
}
func OrderNewNumber(bear BearerResponse) VerificationObject {

	url := "https://www.textverified.com/api/Verifications"
	method := "POST"

	payload := strings.NewReader(`{` + "\n" + `"id": 53` + "\n" + `}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		panic(err)
	}
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", "Bearer "+bear.BearerToken)

	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	var numRes VerificationObject
	if err := json.Unmarshal(body, &numRes); err != nil { // Parse []byte to go struct pointer
		fmt.Println("Can not unmarshal JSON")
	}

	return numRes
}
func CheckExistingNumber(bear BearerResponse, num VerificationObject) string {

	url := "https://www.textverified.com/api/Verifications/" + num.ID
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		panic(err)
	}
	req.Header.Add("Authorization", "Bearer "+bear.BearerToken)

	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	var numRes VerificationObject
	if err := json.Unmarshal(body, &numRes); err != nil { // Parse []byte to go struct pointer
		fmt.Println("Can not unmarshal JSON")
	}

	if numRes.Code != "null" {
		return numRes.Code
	} else {
		time.Sleep(time.Duration(rand.Intn(120)) * time.Second)
		return CheckExistingNumber(bear, num)
	}
}
