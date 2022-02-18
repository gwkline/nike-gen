package main

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
)

//////////////////////////////////////////////////////////////////////////////////////////
//									GOOGLE TASKS										//
//////////////////////////////////////////////////////////////////////////////////////////
func googleTask(ctx context.Context, tid string) chromedp.Tasks {
	return chromedp.Tasks{
		print("Task ID: " + tid + " - Beginning Google"),
		chromedp.Navigate("https://www.google.com"),
		chromedp.Sleep(time.Duration(3 * time.Second)),
		chromedp.WaitVisible(`#Mses6b`),
		chromedp.Click(googleSearchButton, chromedp.ByQuery),
		chromedp.SendKeys(googleSearchButton, "site:nike.com/t/ "+randomProd(tid)+kb.Enter, chromedp.ByQuery),
		chromedp.Sleep(time.Duration(3 * time.Second)),
		linkSelector(tid, ctx),
	}
}

func randomProd(tid string) string {

	choices := []string{"shirt", "shorts", "accessories", "hats"}
	rand.Seed(time.Now().UnixNano())
	num := rand.Intn(len(choices))

	fmt.Println("Task ID: " + tid + " - Search String Shosen: " + choices[num])

	return choices[num]
}

func linkSelector(tid string, ctx context.Context) chromedp.Tasks {

	var nodes []*cdp.Node
	var productLinks []string

	return chromedp.Tasks{

		chromedp.Nodes("a", &nodes),
		chromedp.ActionFunc(func(ctx context.Context) error {

			for _, n := range nodes {
				if strings.Contains(n.AttributeValue("href"), "https://www.nike.com/t/") {
					productLinks = append(productLinks, n.AttributeValue("href"))
				}
			}

			if len(productLinks) == 0 {
				//fmt.Println("Task ID: " + tid + " - Error Scraping Links - Restarting")
				return cdp.Error("Error Scraping links")

				//os.Exit(999)
			}

			//Random Link Selector
			rand.Seed(time.Now().UnixNano())
			randIdx := rand.Intn(len(productLinks))
			randURL := productLinks[randIdx]
			searchString := `[href="` + randURL + `"]`
			fmt.Println("Task ID: " + tid + " - Product Chosen (" + randURL + ")")
			chromedp.Action(chromedp.Click(searchString, chromedp.ByQuery)).Do(ctx)
			return nil
		}),
	}
}
