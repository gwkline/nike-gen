package main

import (
	"bufio"
	"context"
	"encoding/json"
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

func randDelay(length int, ctx context.Context) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			rand.Seed(time.Now().UnixNano())
			num := rand.Intn(150)
			switch num {
			case 0:
				chromedp.Sleep(time.Duration(num+0) * time.Millisecond).Do(ctx)
			case 1:
				chromedp.Sleep(time.Duration(num+10) * time.Millisecond).Do(ctx)
			case 2:
				chromedp.Sleep(time.Duration(num+20) * time.Millisecond).Do(ctx)
			case 3:
				chromedp.Sleep(time.Duration(num+30) * time.Millisecond).Do(ctx)
			case 4:
				chromedp.Sleep(time.Duration(num+50) * time.Millisecond).Do(ctx)
			case 5:
				chromedp.Sleep(time.Duration(num+100) * time.Millisecond).Do(ctx)
			case 6:
				chromedp.Sleep(time.Duration(num+175) * time.Millisecond).Do(ctx)
			case 7:
				chromedp.Sleep(time.Duration(num+250) * time.Millisecond).Do(ctx)
			case 8:
				chromedp.Sleep(time.Duration(num+500) * time.Millisecond).Do(ctx)
			case 9:
				chromedp.Sleep(time.Duration(num+750) * time.Millisecond).Do(ctx)
			case 10:
				chromedp.Sleep(time.Duration(num+1000) * time.Millisecond).Do(ctx)
			}
			return nil
		}),
	}
}

func print(message string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			fmt.Println(message)
			return nil
		}),
	}
}

func monthToDigit(month string) string {
	switch month {
	case "January":
		return "01"
	case "February":
		return "02"
	case "March":
		return "03"
	case "April":
		return "04"
	case "May":
		return "05"
	case "June":
		return "06"
	case "July":
		return "07"
	case "August":
		return "08"
	case "September":
		return "09"
	case "October":
		return "10"
	case "November":
		return "11"
	case "December":
		return "12"

	}
	return "ERROR"
}

func PrettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

func testTask(tid string) chromedp.Tasks {
	return chromedp.Tasks{
		print("Task ID: " + tid + " - Beginning Debug"),
		chromedp.Navigate("https://detect.azerpas.com"),
		print("Task ID: " + tid + " - Waiting: 10s"),
		chromedp.Sleep(time.Duration(10 * time.Second)),
	}
}

func loadProxies() []Proxy {

	var proxies []Proxy

	// open file
	f, err := os.Open("proxies.txt")
	if err != nil {
		os.Exit(133)
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
		os.Exit(130)
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(proxies), func(i, j int) { proxies[i], proxies[j] = proxies[j], proxies[i] })
	return proxies
}

func loadEmails() []string {

	var emails []string
	// open file
	f, err := os.Open("emails.txt")
	if err != nil {
		os.Exit(131)
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
		os.Exit(132)
	}

	return emails
}
