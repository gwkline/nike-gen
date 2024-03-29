package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

//////////////////////////////////////////////////////////////////////////////////////////
//								TASK ENGINE AND MAIN									//
//////////////////////////////////////////////////////////////////////////////////////////
//TODO:
// - MAKE INTO GOROUTINE
// - BEGIN BROWSER CONFIGURATION (TLS?)
// - INDIVIDUAL STEP TIMEOUTS
// - CLEAN TASK KILL/SCRIPT KILL UPON BROWSER CLOSE
// - RETRY LOGIC
// - ACCOUNT OUTPUT FILE
// - INDIVIDUAL CHROMEDP TASK ERROR HANDLING (IF NEEDED)
// - IMPROVE randomProd
// - IMRPROVE HUMAN EMULATION / ADD SCROLLING
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
		task.Task_ID = fmt.Sprint(i)

		//TODO: FIX OUT OF INDEX ERROR
		//proxies = append(proxies[:randIdx], proxies[randIdx+1:]...)
		//emails = append(emails[:i], emails[i+1:]...)

		fmt.Println("")
		log(&task, "Starting")
		log(&task, "Proxy Being Used: "+task.Proxy.IP)
		log(&task, "Email Being Used: "+task.Email)

		err = runTasks(&task)
		if err != nil && maxAttempts(task) <= RETRY_LIMIT {
			fmt.Printf("Task ID: %s | ERROR - %+v\n", task.Task_ID, err.Error())
			fmt.Printf("Task ID: %s | Attempting retry\n", task.Task_ID)

			//TODO: RETRY LOGIC
			//err = runTasks(&task)
			runTasks(&task)

		} else if err != nil && maxAttempts(task) > RETRY_LIMIT {
			fmt.Printf("Task ID: %s | ERROR - %+v\n", task.Task_ID, err.Error())
			fmt.Printf("Task ID: %s | Task retry attempts exhausted\n", task.Task_ID)
			os.Exit(999)
		} else {
			log(&task, "Creation Complete")

			//TODO: OUTPUT ACCOUNT
		}
	}
}
