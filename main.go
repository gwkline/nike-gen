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
// - MAKE INTO GOROUTINE
// - BEGIN BROWSER CONFIGURATION (TLS?)
//
//
// FURTHER TODO'S LISTED BELOW

func main() {

	proxies, err := loadProxies()
	if err != nil {
		fmt.Printf("\nProxy initialization error, %s", err.Error())
	}

	emails, err := loadEmails()
	if err != nil && USE_EMAIL_LIST {
		fmt.Printf("\n%s", err.Error())
	}

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

		err = runTasks(&task)
		if err != nil && task.Attempt < 3 {
			fmt.Printf("Task ID: %s | ERROR - %+v\n", task.Task_ID, err.Error())
			fmt.Printf("Task ID: %s | Attempting retry\n", task.Task_ID)
			task.Attempt++
			//err = runTasks(&task)
		} else if err != nil && task.Attempt > 2 {
			fmt.Printf("Task ID: %s | ERROR - %+v\n", task.Task_ID, err.Error())
			fmt.Printf("Task ID: %s | Task retry attempts exhausted\n", task.Task_ID)
			//break
		} else {
			log(&task, "Creation Complete")
		}
	}
}
