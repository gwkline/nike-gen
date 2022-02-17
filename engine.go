package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/fetch"
	"github.com/chromedp/chromedp"
)

func runTasks(proxy Proxy, email string, tid string) {

	fmt.Println("Task ID: " + tid + " - Initializing Browser")

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag(`headless`, false),
		chromedp.Flag(`incognito`, true),
		chromedp.DisableGPU,
		chromedp.Flag(`disable-extensions`, false),
		chromedp.Flag(`enable-automation`, false),
		chromedp.WindowSize(1876, 896), //any width under 1024 will have hamburger menu bar
		//chromedp.UserAgent(uarand.GetRandom()),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.82 Safari/537.36"),
		chromedp.ProxyServer(""),
		//chromedp.ProxyServer("http://"+proxy.IP+":"+proxy.Port),
	)
	parentCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	ctx, cancel := chromedp.NewContext(parentCtx)
	// ctx, cancel = chromedp.NewContext(ctx, chromedp.WithDebugf(log.Printf))
	defer cancel()

	//create a global task timeout
	//TODO: ADD ADJUSTABILITY
	ctx, cancel = context.WithTimeout(ctx, 600*time.Second)
	defer cancel()
	lctx, lcancel := context.WithCancel(ctx)
	chromedp.ListenTarget(lctx, func(ev interface{}) {
		switch ev := ev.(type) {

		case *fetch.EventRequestPaused:
			go func() {
				_ = chromedp.Run(ctx, fetch.ContinueRequest(ev.RequestID))
			}()

		case *fetch.EventAuthRequired:
			if ev.AuthChallenge.Source == fetch.AuthChallengeSourceProxy {
				go func() {
					_ = chromedp.Run(ctx,
						fetch.ContinueWithAuth(ev.RequestID, &fetch.AuthChallengeResponse{
							Response: fetch.AuthChallengeResponseResponseProvideCredentials,
							Username: proxy.User,
							Password: proxy.Pass,
						}),
						fetch.Disable(),
					)
					lcancel()
				}()
			}
		}
	})

	//googles a random nike product
	//TODO: TRUE RANDOM PRODUCT SEARCH
	var nodes []*cdp.Node
	var productLinks []string

	//Google Tasks (And Debug)
	fmt.Println("Task ID: " + tid + " - Launching Browser")
	err := chromedp.Run(ctx,
		fetch.Enable().WithHandleAuthRequests(true),
		//testTask(tid),
		//chromedp.Sleep(10*time.Second),
		googleTask(ctx, tid),
		chromedp.Nodes("a", &nodes),
	)
	if err != nil {
		panic(err)
	}

	//Product Link Collector
	for _, n := range nodes {
		if strings.Contains(n.AttributeValue("href"), "https://www.nike.com/t/") {
			productLinks = append(productLinks, n.AttributeValue("href"))
		}
	}
	if len(productLinks) == 0 {
		os.Exit(129)
	}

	//Random Link Selector
	rand.Seed(time.Now().UnixNano())
	randIdx := rand.Intn(len(productLinks))
	randURL := productLinks[randIdx]
	//searchString := `[href="` + randURL + `"]`
	fmt.Println("Task ID: " + tid + " - Product Chosen (" + randURL + ")")

	//Login Tasks

	err = chromedp.Run(ctx,
		fetch.Enable().WithHandleAuthRequests(true),
		randDelay(10, ctx),
		//chromedp.Click(searchString, chromedp.BySearch),
		print("Task ID: "+tid+" - Beginning Sign Up"),
		nikeSignupTask(ctx, tid),
		print("Task ID: "+tid+" - Signup Complete"),
		print("Task ID: "+tid+" - Navigating To Settings"),
		nikeGoToPhoneNumber(ctx, tid),
	)
	if err != nil {
		panic(err)
	}

	//SMS Tasks (Init)
	fmt.Println("Task ID: " + tid + " - Getting SMS Token")
	token := GetSMSToken(tid)
	fmt.Println("Task ID: " + tid + " - Ordering New Number")
	order := OrderNewNumber(token, tid)
	err = chromedp.Run(ctx,
		fetch.Enable().WithHandleAuthRequests(true),
		print("Task ID: "+tid+" - Beginning Number Input Process"),
		nikeInputPhoneNumber(string(order.Number), ctx, tid))
	if err != nil {
		panic(err)
	}

	//SMS Tasks (Confirm)
	fmt.Println("Task ID: " + tid + " - Waiting For Code")
	code := CheckExistingNumber(token, order, tid, 1)
	err = chromedp.Run(ctx,
		fetch.Enable().WithHandleAuthRequests(true),
		nikeConfirmPhoneNumber(code, ctx, tid))
	if err != nil {
		panic(err)
	}
}
