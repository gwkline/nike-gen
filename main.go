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

		rand.Seed(time.Now().UnixNano())
		randIdx := rand.Intn(len(proxies) - 1)

		randProx := proxies[randIdx]
		email := emails[0]
		//TODO: FIX OUT OF INDEX ERROR
		//proxies = append(proxies[:randIdx], proxies[randIdx+1:]...)
		//emails = append(emails[:i], emails[i+1:]...)

		if DEBUG {
			fmt.Println("")
			fmt.Println("Task ID: " + fmt.Sprint(i) + " | Starting")
			fmt.Println("Task ID: " + fmt.Sprint(i) + " | Proxy Being Used: " + randProx.IP)
			fmt.Println("Task ID: " + fmt.Sprint(i) + " | Email Being Used: " + email)
		}

		runTasks(randProx, email, fmt.Sprint(i), 1)
	}
}
