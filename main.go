package main

import (
	"fmt"
	"math/rand"
	"time"
)

//////////////////////////////////////////////////////////////////////////////////////////
//								TASK ENGINE AND MAIN									//
//////////////////////////////////////////////////////////////////////////////////////////
//TODO:
// - MAKE THREADED
// - BEGIN BROWSER CONFIGURATION (TLS?)
// -
//
// FURTHER TODO'S LISTED BELOW

func main() {

	proxies := loadProxies()
	emails := loadEmails()

	for i := range emails {

		var task Task

		rand.Seed(time.Now().UnixNano())
		randIdx := rand.Intn(len(proxies) - 1)

		task.Proxy = proxies[randIdx]
		task.Email = emails[0]
		task.Attempt = 1
		task.Task_ID = fmt.Sprint(i)

		//TODO: FIX OUT OF INDEX ERROR
		//proxies = append(proxies[:randIdx], proxies[randIdx+1:]...)
		//emails = append(emails[:i], emails[i+1:]...)

		fmt.Println("")
		log(&task, "Starting")
		log(&task, "Proxy Being Used: "+task.Proxy.IP)
		log(&task, "Email Being Used: "+task.Email)

		runTasks(&task)
	}
}
