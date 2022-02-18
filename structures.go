package main

import "time"

//////////////////////////////////////////////////////////////////////////////////////////
//							CUSTOM OBJECTS AND STRUCTURES								//
//////////////////////////////////////////////////////////////////////////////////////////
type Proxy struct {
	User string `json:"user"`
	Pass string `json:"pass"`
	IP   string `json:"ip"`
	Port string `json:"port"`
}

type BearerResponse struct {
	BearerToken string    `json:"bearer_token"`
	Expiration  time.Time `json:"expiration"`
	Ticks       int64     `json:"ticks"`
}

type VerificationObject struct {
	ID              string      `json:"id"`
	Cost            float64     `json:"cost"`
	TargetName      string      `json:"target_name"`
	Number          string      `json:"number"`
	TimeRemaining   string      `json:"time_remaining"`
	ReuseWindow     interface{} `json:"reuse_window"`
	Status          string      `json:"status"`
	Sms             interface{} `json:"sms"`
	Code            string      `json:"code"`
	VerificationURI string      `json:"verification_uri"`
	CancelURI       string      `json:"cancel_uri"`
	ReportURI       string      `json:"report_uri"`
	ReuseURI        string      `json:"reuse_uri"`
}

type Task struct {
	Proxy      Proxy  `json:"proxy"`
	Email      string `json:"email"`
	Task_ID    string `json:"tid"`
	Attempt    int    `json:"attempt"`
	Status     string `json:"status"`
	First_Name string `json:"fname"`
	Last_Name  string `json:"lname"`
	DOB        string `json:"dob"`
	Gender     int    `json:"gender"`
}
