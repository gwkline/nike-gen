package main

import (
	"encoding/json"
	"errors"
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
func GetSMSToken(task *Task) (BearerResponse, error) {

	//Request Setup
	tvURL := "https://www.textverified.com/api/"
	client := http.Client{}
	req, err := http.NewRequest("POST", tvURL+"SimpleAuthentication", nil)
	if err != nil {
		nilbear := BearerResponse{"Auth Not Found", time.Now(), 0}
		return nilbear, errors.New("GetSMSToken: Error creating authorization request")
	}
	req.Header = http.Header{
		"X-SIMPLE-API-ACCESS-TOKEN": []string{"1_l8T2o8bPsP252roHDXhtO-tRzX3tqROlPzzam8kTaj7YPRlnccMekpuAmRsDh4r9_H13sjOr"},
	}

	//Performing Request
	log(task, "Getting Bearer Token")
	res, err := client.Do(req)
	if err != nil {
		nilbear := BearerResponse{"Auth Not Found", time.Now(), 0}
		return nilbear, errors.New("GetSMSToken: Cannot find bearer token")

	} else {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			nilbear := BearerResponse{"Auth Not Found", time.Now(), 0}
			return nilbear, errors.New("GetSMSToken: Error reading response body")
		}
		//bodyString := string(bodyBytes)
		var bear BearerResponse
		if err := json.Unmarshal(bodyBytes, &bear); err != nil { // Parse []byte to go struct pointer
			nilbear := BearerResponse{"Auth Not Found", time.Now(), 0}
			return nilbear, errors.New("GetSMSToken: Cannot unmarshal JSON response")
		}
		return bear, nil
	}
}

func OrderNewNumber(bear BearerResponse, task *Task) (VerificationObject, error) {

	url := "https://www.textverified.com/api/Verifications"
	method := "POST"
	payload := strings.NewReader(`{` + "\n" + `"id": 53` + "\n" + `}`)
	client := &http.Client{}

	log(task, "Ordering New Number")
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		nilVerOb := VerificationObject{}
		return nilVerOb, errors.New("OrderNewNumber: Error creating order request")
	}
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", "Bearer "+bear.BearerToken)

	res, err := client.Do(req)
	if err != nil {
		nilVerOb := VerificationObject{}
		return nilVerOb, errors.New("OrderNewNumber: Error performing order request")
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		nilVerOb := VerificationObject{}
		return nilVerOb, errors.New("OrderNewNumber: Error reading response body")
	}
	var numRes VerificationObject
	if err := json.Unmarshal(body, &numRes); err != nil { // Parse []byte to go struct pointer
		nilVerOb := VerificationObject{}
		return nilVerOb, errors.New("OrderNewNumber: Cannot unmarshal JSON response")
	}
	log(task, "Number Recieved")
	return numRes, nil
}

func CheckExistingNumber(bear BearerResponse, num VerificationObject, task *Task, iter int) (string, error) {
	log(task, "Waiting For Code - Attempt: "+fmt.Sprint(iter))

	//Request Setup
	url := "https://www.textverified.com/api/Verifications/" + num.ID
	method := "GET"
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return "", errors.New("CheckExistingNumber: Error creating SMS update request")
	}
	req.Header.Add("Authorization", "Bearer "+bear.BearerToken)

	//Performing Request
	res, err := client.Do(req)
	if err != nil {
		return "", errors.New("CheckExistingNumber: Error performing SMS update request")
	}
	defer res.Body.Close()

	//Reading Request
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", errors.New("CheckExistingNumber: Error reading SMS update response body")
	}

	//Parsing Response
	var numRes VerificationObject
	if err := json.Unmarshal(body, &numRes); err != nil { // Parse []byte to go struct pointer
		return "", errors.New("CheckExistingNumber: Cannot unmarshal JSON response")
	}

	if numRes.Code != "null" {
		return numRes.Code, nil
	} else {
		time.Sleep(time.Duration(rand.Intn(120)) * time.Second)
		return CheckExistingNumber(bear, num, task, iter+1)
	}
}
