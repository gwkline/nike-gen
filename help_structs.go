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
	Proxy      Proxy     `json:"proxy"`
	Email      string    `json:"email"`
	Task_ID    string    `json:"tid"`
	Status     string    `json:"status"`
	First_Name string    `json:"fname"`
	Last_Name  string    `json:"lname"`
	DOB        [3]string `json:"dob"`
	Gender     int       `json:"gender"`
	Attempts   Attempts  `json:"attempts"`
}

type Attempts struct {
	Google     int `json:"googleatt"`
	Signup     int `json:"signupatt"`
	Navigate   int `json:"navatt"`
	OrderNum   int `json:"orderatt"`
	CheckNum   int `json:"checkatt"`
	ConfirmNum int `json:"confirmatt"`
	InputNum   int `json:"inputatt"`
	SMSAuth    int `json:"authatt"`
	Nike       int `json:"nikeatt"`
}

func maxAttempts(att Task) int {
	Google := att.Attempts.Google
	Signup := att.Attempts.Signup
	Navigate := att.Attempts.Navigate
	OrderNum := att.Attempts.OrderNum
	CheckNum := att.Attempts.CheckNum
	ConfirmNum := att.Attempts.ConfirmNum
	InputNum := att.Attempts.InputNum
	SMSAuth := att.Attempts.SMSAuth
	Nike := att.Attempts.Nike

	x := []int{Google, Signup, Navigate, OrderNum, CheckNum, ConfirmNum, InputNum, SMSAuth, Nike}

	biggest := x[0]
	for _, v := range x {
		if v > biggest {
			biggest = v
		}

	}

	return biggest

}
