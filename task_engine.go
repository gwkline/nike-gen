package main

import (
	"context"
	"errors"
	"time"

	"github.com/chromedp/cdproto/fetch"
	"github.com/chromedp/chromedp"
)

func runTasks(task *Task) error {

	if task.Attempt > 2 {
		return errors.New("runTasks: Task attempt limit reached")
	}

	log(task, "Initializing Browser")

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
							Username: task.Proxy.User,
							Password: task.Proxy.Pass,
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
		googleTask(ctx, task),
	)
	if err != nil {

		cancel()
		time.Sleep(5 * time.Second)
		task.Attempt++
		return errors.New("googleTask: Error navigating to random product page")

	}

	//Login Tasks
	err = chromedp.Run(ctx,
		fetch.Enable().WithHandleAuthRequests(true),
		nikeSignupTask(ctx, task),
	)
	if err != nil {
		cancel()
		time.Sleep(5 * time.Second)
		task.Attempt++
		return errors.New("nikeSignupTask: Error creating Nike account")
	}

	//Login Tasks
	err = chromedp.Run(ctx,
		fetch.Enable().WithHandleAuthRequests(true),
		nikeGoToPhoneNumber(ctx, task),
	)
	if err != nil {
		cancel()
		time.Sleep(5 * time.Second)
		task.Attempt++
		return errors.New("nikeGoToPhoneNumber: Error navigating to profile SMS page")
	}

	//SMS Tasks (Init)
	token, err := GetSMSToken(task)
	if err != nil {
		cancel()
		time.Sleep(5 * time.Second)
		task.Attempt++
		return err
	}

	order, err := OrderNewNumber(token, task)
	if err != nil {
		cancel()
		time.Sleep(5 * time.Second)
		task.Attempt++
		return err
	}

	err = chromedp.Run(ctx,
		fetch.Enable().WithHandleAuthRequests(true),
		nikeInputPhoneNumber(string(order.Number), ctx, task),
	)
	if err != nil {
		cancel()
		time.Sleep(5 * time.Second)
		task.Attempt++
		return errors.New("nikeInputPhoneNumber: Error inputting phone number")
	}

	//SMS Tasks (Confirm)
	code, err := CheckExistingNumber(token, order, task, 1)
	if err != nil {
		cancel()
		time.Sleep(5 * time.Second)
		task.Attempt++
		return err
	}

	err = chromedp.Run(ctx,
		fetch.Enable().WithHandleAuthRequests(true),
		nikeConfirmPhoneNumber(code, ctx, task))
	if err != nil {
		cancel()
		time.Sleep(5 * time.Second)
		task.Attempt++
		return errors.New("nikeConfirmPhoneNumber: Error confirming phone number")
	}

	return nil
}
