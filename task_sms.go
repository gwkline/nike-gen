package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

//////////////////////////////////////////////////////////////////////////////////////////
//								SMS VERIFICATION TASKS									//
//////////////////////////////////////////////////////////////////////////////////////////
func GetSMSToken(task *Task) BearerResponse {

	//Request Setup
	tvURL := "https://www.textverified.com/api/"
	client := http.Client{}
	req, err := http.NewRequest("POST", tvURL+"SimpleAuthentication", nil)
	if err != nil {
		os.Exit(134)
	}
	req.Header = http.Header{
		"X-SIMPLE-API-ACCESS-TOKEN": []string{"1_l8T2o8bPsP252roHDXhtO-tRzX3tqROlPzzam8kTaj7YPRlnccMekpuAmRsDh4r9_H13sjOr"},
	}

	//Performing Request
	log(task, "Getting Bearer Token")
	res, err := client.Do(req)
	if err != nil {

		nilbear := BearerResponse{"Auth Not Found", time.Now(), 0}
		os.Exit(135)
		return nilbear

	} else {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			os.Exit(136)
		}
		//bodyString := string(bodyBytes)
		var bear BearerResponse
		if err := json.Unmarshal(bodyBytes, &bear); err != nil { // Parse []byte to go struct pointer
			fmt.Println("Can not unmarshal JSON")
		}
		return bear
	}
}

func OrderNewNumber(bear BearerResponse, task *Task) VerificationObject {

	url := "https://www.textverified.com/api/Verifications"
	method := "POST"
	payload := strings.NewReader(`{` + "\n" + `"id": 53` + "\n" + `}`)
	client := &http.Client{}

	log(task, "Ordering New Number")
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		os.Exit(137)
	}
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", "Bearer "+bear.BearerToken)

	res, err := client.Do(req)
	if err != nil {
		os.Exit(138)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		os.Exit(139)
	}
	var numRes VerificationObject
	if err := json.Unmarshal(body, &numRes); err != nil { // Parse []byte to go struct pointer
		fmt.Println("Can not unmarshal JSON")
	}
	log(task, "Number Recieved")
	return numRes
}

func CheckExistingNumber(bear BearerResponse, num VerificationObject, task *Task, iter int) string {
	log(task, "Waiting For Code - Attempt: "+fmt.Sprint(iter))

	//Request Setup
	url := "https://www.textverified.com/api/Verifications/" + num.ID
	method := "GET"
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		os.Exit(140)
	}
	req.Header.Add("Authorization", "Bearer "+bear.BearerToken)

	//Performing Request
	res, err := client.Do(req)
	if err != nil {
		os.Exit(141)
	}
	defer res.Body.Close()

	//Reading Request
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		os.Exit(142)
	}

	//Parsing Response
	var numRes VerificationObject
	if err := json.Unmarshal(body, &numRes); err != nil { // Parse []byte to go struct pointer
		fmt.Println("Can not unmarshal JSON")
	}

	if numRes.Code != "null" {
		return numRes.Code
	} else {
		time.Sleep(time.Duration(rand.Intn(120)) * time.Second)
		return CheckExistingNumber(bear, num, task, iter+1)
	}
}
