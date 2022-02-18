package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/chromedp/cdproto/fetch"
	"github.com/chromedp/chromedp"
)

func runTasks(proxy Proxy, email string, tid string, attempt int) {

	if attempt > 2 {
		return
	}

	fmt.Println("Task ID: " + tid + " | Initializing Browser")

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

	//TODO: TRUE RANDOM PRODUCT SEARCH
	//Google Tasks
	err := chromedp.Run(ctx,
		fetch.Enable().WithHandleAuthRequests(true),
		googleTask(ctx, tid),
	)
	if err != nil {
		// fmt.Println(err)
		// os.Exit(145)
		cancel()
		time.Sleep(5 * time.Second)
		runTasks(proxy, email, tid, attempt+1)

	}

	//Login Tasks
	err = chromedp.Run(ctx,
		fetch.Enable().WithHandleAuthRequests(true),
		nikeSignupTask(ctx, tid),
		nikeGoToPhoneNumber(ctx, tid),
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(146)
	}

	//SMS Tasks (Init)
	token := GetSMSToken(tid)
	order := OrderNewNumber(token, tid)
	err = chromedp.Run(ctx,
		fetch.Enable().WithHandleAuthRequests(true),
		nikeInputPhoneNumber(string(order.Number), ctx, tid),
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(147)
	}

	//SMS Tasks (Confirm)
	code := CheckExistingNumber(token, order, tid, 1)
	err = chromedp.Run(ctx,
		fetch.Enable().WithHandleAuthRequests(true),
		nikeConfirmPhoneNumber(code, ctx, tid))
	if err != nil {
		panic(err)
	}
}

func logEngine()
