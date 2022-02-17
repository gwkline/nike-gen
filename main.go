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

	for i := range proxies {

		rand.Seed(time.Now().UnixNano())
		randIdx := rand.Intn(len(proxies))

		randProx := proxies[randIdx]
		email := emails[0]
		proxies = append(proxies[:randIdx], proxies[randIdx+1:]...)
		emails = append(emails[:i], emails[i+1:]...)

		fmt.Println("")
		fmt.Println("Task: " + fmt.Sprint(i) + " - Starting")
		fmt.Println("Proxy being used: " + randProx.IP)
		fmt.Println("Email being used: " + email)
		fmt.Println("")

		runTasks(randProx, email, fmt.Sprint(i))
	}
}
