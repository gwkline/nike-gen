package main

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/fetch"
	"github.com/chromedp/chromedp"
)

func runTasks(proxy Proxy, email string, tid string) {

	//GET DOM TRAVERSAL VALUES
	favoriteBtn := `#floating-atc-wrapper > div > button.wishlist-btn.ncss-btn-secondary-dark.btn-lg.mt3-sm`
	registerBtn := `.loginJoinLink.current-member-signin > a`
	genderBtn := `li:nth-child(1) > input[type="button"]`
	// phoneBtn := `div.sendCode > div.mobileNumber-div > input`
	// phoneloginBtn := `#root > div > div > div.main-layout > div > header > div.d-sm-h.d-lg-b > section > div > ul > li.member-nav-item.d-sm-ib.va-sm-m > div > div > button`
	// settingBtn := `#root > div > div > div.main-layout > div > header > div.d-sm-h.d-lg-b > section > div > ul > li.member-nav-item.d-sm-ib.va-sm-m > div > div > ul > li:nth-child(1)`
	// addBtn := `#mobile-container > div > div > form > div.account-form > div.mex-mobile-input-wrapper.ncss-col-sm-12.ncss-col-md-12.pl0-sm.pr0-sm.pb3-sm > div > div > div > div.ncss-col-sm-6.ta-sm-r.va-sm-m.flx-jc-sm-fe.d-sm-iflx > button`
	// numcountryBtn := `select[class="country"]`
	// sendNumBtn := `#nike-unite-progressiveForm > div > div > input[type="button"]`
	// enterTheValueBtn := `input[type="number"]`
	// storedSubmitBtn := `#nike-unite-progressiveForm > div > input[type="button"]`
	// acceptCookies := `#cookie-settings-layout > div > div > div > div.ncss-row.mt5-sm.mb7-sm > div:nth-child(2) > button`
	// loginBtn := `li.member-nav-item.d-sm-ib.va-sm-m > button`

	fmt.Println("Task ID: " + tid + " - Initializing")

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

	fmt.Println("Task ID: " + tid + " - Beginning Google")
	err := chromedp.Run(ctx,
		fetch.Enable().WithHandleAuthRequests(true),
		testTask(),
		chromedp.Sleep(10*time.Second),
		googleTask(ctx, tid),
		chromedp.Nodes("a", &nodes),
	)
	if err != nil {
		panic(err)
	}

	//Product link collector
	for _, n := range nodes {
		if strings.Contains(n.AttributeValue("href"), "https://www.nike.com/t/") {
			fmt.Println(n.AttributeValue("href"))
			productLinks = append(productLinks, n.AttributeValue("href"))
		}
	}

	//Random link selector
	rand.Seed(time.Now().UnixNano())
	randIdx := rand.Intn(len(productLinks))
	randURL := productLinks[randIdx]
	searchString := `[href="` + randURL + `"]`
	fmt.Println("Task ID: " + tid + " - Product Chosen (" + randURL + ")")

	fmt.Println("Task ID: " + tid + " - Beginning Login")
	err = chromedp.Run(ctx,
		fetch.Enable().WithHandleAuthRequests(true),
		randDelay(10, ctx),
		chromedp.Click(searchString, chromedp.BySearch),

		nikeSignupTask(favoriteBtn, registerBtn, genderBtn, ctx),
		print("Signup Complete"),
		nikeGoToPhoneNumber(ctx),
	)
	if err != nil {
		panic(err)
	}

	token := GetSMSToken()
	order := OrderNewNumber(token)
	err = chromedp.Run(ctx,
		fetch.Enable().WithHandleAuthRequests(true),
		nikeInputPhoneNumber(string(order.Number), ctx))
	if err != nil {
		panic(err)
	}

	code := CheckExistingNumber(token, order)
	err = chromedp.Run(ctx,
		fetch.Enable().WithHandleAuthRequests(true),
		nikeConfirmPhoneNumber(code, ctx))
	if err != nil {
		panic(err)
	}
}
