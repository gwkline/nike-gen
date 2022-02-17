package main

import (
	"context"
	"math/rand"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
)

//////////////////////////////////////////////////////////////////////////////////////////
//									GOOGLE TASKS										//
//////////////////////////////////////////////////////////////////////////////////////////
func googleTask(ctx context.Context, tid string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate("https://www.google.com"),
		chromedp.Sleep(time.Duration(3 * time.Second)),
		chromedp.WaitVisible(`#Mses6b`),
		chromedp.Click(`body > div.L3eUgb > div.o3j99.ikrT4e.om7nvf > form > div:nth-child(1) > div.A8SBwf > div.RNNXgb > div > div.a4bIc > input`, chromedp.ByQuery),
		//TODO: Random nike product selector
		typeWord(`body > div.L3eUgb > div.o3j99.ikrT4e.om7nvf > form > div:nth-child(1) > div.A8SBwf > div.RNNXgb > div > div.a4bIc > input`, "site:nike.com/t/ "+randomProd(tid)+kb.Enter, chromedp.ByQuery, ctx),
		chromedp.Sleep(time.Duration(3 * time.Second)),
		chromedp.Click(`h3`, chromedp.ByQuery),
		chromedp.Sleep(time.Duration(rand.Intn(6)) * time.Second),
	}
}
