package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
)

//////////////////////////////////////////////////////////////////////////////////////////
//						HELPER FUNCTIONS AND OTHER UTILITIES							//
//////////////////////////////////////////////////////////////////////////////////////////
func typeWord(sel interface{}, word string, opts func(*chromedp.Selector), ctx context.Context) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {

			runeList := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
			for _, v := range word {
				rand.Seed(time.Now().UnixNano())
				randomInt := rand.Intn(999)
				randDel := (rand.Intn(150) + 75)

				if randomInt < 939 {
					chromedp.SendKeys(sel, string(v), opts).Do(ctx)
					chromedp.Sleep(time.Duration(randDel) * time.Millisecond).Do(ctx)
				} else if randomInt < 965 {
					chromedp.SendKeys(sel, string(byte(runeList[rand.Intn(len(runeList))])), opts).Do(ctx)
					chromedp.Sleep(time.Duration(randDel) * time.Millisecond).Do(ctx)
					chromedp.SendKeys(sel, string(byte(runeList[rand.Intn(len(runeList))])), opts).Do(ctx)
					chromedp.Sleep(time.Duration(randDel) * time.Millisecond).Do(ctx)
					chromedp.SendKeys(sel, kb.Backspace, opts).Do(ctx)
					chromedp.Sleep(time.Duration(randDel) * time.Millisecond).Do(ctx)
					chromedp.SendKeys(sel, kb.Backspace, opts).Do(ctx)
					chromedp.Sleep(time.Duration(randDel) * time.Millisecond).Do(ctx)
					chromedp.SendKeys(sel, string(v), opts).Do(ctx)
				} else {
					chromedp.SendKeys(sel, string(byte(runeList[rand.Intn(len(runeList))])), opts).Do(ctx)
					chromedp.Sleep(time.Duration(randDel) * time.Millisecond).Do(ctx)
					chromedp.SendKeys(sel, kb.Backspace, opts).Do(ctx)
					chromedp.Sleep(time.Duration(randDel) * time.Millisecond).Do(ctx)
					chromedp.SendKeys(sel, string(v), opts).Do(ctx)
				}
			}

			return nil
		}),
	}
}

func randDelay(ctx context.Context) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			rand.Seed(time.Now().UnixNano())
			num := rand.Intn(MAX_DELAY-MIN_DELAY) + MIN_DELAY
			chromedp.Sleep(time.Duration(num) * time.Millisecond).Do(ctx)
			return nil
		}),
	}
}

func log(task *Task, status string) {
	task.Status = status
	fmt.Printf("Task ID: %s | %s\n", task.Task_ID, task.Status)
}

func logTask(task *Task, status string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			task.Status = status
			fmt.Printf("Task ID: %s | %s\n", task.Task_ID, task.Status)
			return nil
		}),
	}
}

func monthToDigit(month string) (string, error) {
	switch month {
	case "January":
		return "01", nil
	case "February":
		return "02", nil
	case "March":
		return "03", nil
	case "April":
		return "04", nil
	case "May":
		return "05", nil
	case "June":
		return "06", nil
	case "July":
		return "07", nil
	case "August":
		return "08", nil
	case "September":
		return "09", nil
	case "October":
		return "10", nil
	case "November":
		return "11", nil
	case "December":
		return "12", nil

	}
	return "", errors.New("monthToDigit: Error with given input: " + month)
}

func PrettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

// func testTask(task *Task) chromedp.Tasks {
// 	return chromedp.Tasks{
// 		logTask(task, "Beginning Debug"),
// 		chromedp.Navigate("https://detect.azerpas.com"),
// 		logTask(task, "Waiting 10s"),
// 		chromedp.Sleep(time.Duration(10 * time.Second)),
// 	}
// }

func loadProxies() ([]Proxy, error) {

	var proxies []Proxy

	// open file
	f, err := os.Open("proxies.txt")
	if err != nil {
		return nil, err
	}
	// remember to close the file at the end of the program
	defer f.Close()

	// read the file line by line using scanner
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		if scanner.Text() != "" {
			split := strings.Split(scanner.Text(), ":")
			var prox Proxy
			prox.IP = split[0]
			prox.Port = split[1]
			prox.User = split[2]
			prox.Pass = split[3]
			proxies = append(proxies, prox)

		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(proxies), func(i, j int) { proxies[i], proxies[j] = proxies[j], proxies[i] })
	return proxies, nil
}

func loadEmails() ([]string, error) {

	var emails []string
	// open file
	f, err := os.Open("emails.txt")
	if err != nil {
		return nil, err
	}
	// remember to close the file at the end of the program
	defer f.Close()

	// read the file line by line using scanner
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		if scanner.Text() != "" {
			emails = append(emails, scanner.Text())

		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return emails, nil
}
