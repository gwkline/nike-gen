package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/fetch"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
	"github.com/sethvargo/go-password/password"
)

//////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////
type Proxy struct {
	User string `json:"user"`
	Pass string `json:"pass"`
	IP   string `json:"ip"`
	Port string `json:"port"`
}

type BearerResponse struct {
	BearerToken string    `json:"bearer_token"`
	Expiration  time.Time `json:"expiration"`
	Ticks       int64     `json:"ticks"`
}

type VerificationObject struct {
	ID              string      `json:"id"`
	Cost            float64     `json:"cost"`
	TargetName      string      `json:"target_name"`
	Number          string      `json:"number"`
	TimeRemaining   string      `json:"time_remaining"`
	ReuseWindow     interface{} `json:"reuse_window"`
	Status          string      `json:"status"`
	Sms             interface{} `json:"sms"`
	Code            string      `json:"code"`
	VerificationURI string      `json:"verification_uri"`
	CancelURI       string      `json:"cancel_uri"`
	ReportURI       string      `json:"report_uri"`
	ReuseURI        string      `json:"reuse_uri"`
}

//////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////
func typeWord(sel interface{}, word string, opts func(*chromedp.Selector)) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.ActionFunc(func(context.Context) error {

			chromedp.Navigate("https://nike.com")
			runeList := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

			for _, v := range word {

				randomInt := rand.Intn(999)

				if randomInt < 920 {

					chromedp.SendKeys(sel, string(v), opts)
					chromedp.Sleep(time.Duration(randomInt) * time.Millisecond)
				} else {
					chromedp.SendKeys(sel, string(byte(runeList[rand.Intn(len(runeList))])), opts)
					chromedp.Sleep(time.Duration(randomInt) * time.Millisecond)
					chromedp.SendKeys(sel, kb.Backspace, opts)
					chromedp.Sleep(time.Duration(randomInt) * time.Millisecond)
					chromedp.SendKeys(sel, string(v), opts)
				}
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

func randomProd() string {
	rand.Seed(time.Now().UnixNano())
	num := rand.Intn(9)
	switch num {
	case 0:
		return "shirt"
	case 1:
		return "shoes"
	case 3:
		return "pants"
	case 4:
		return "socks"
	case 5:
		return "jacket"
	case 6:
		return "fleece"
	case 7:
		return "sweatshirt"
	case 8:
		return "sweat"
	case 9:
		return "dri fit"
	}
	return "clothes"
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

//////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////
func googleTask() chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate("https://www.google.com"),
		chromedp.Sleep(time.Duration(10 * time.Second)),
		chromedp.WaitVisible(`#Mses6b`),
		chromedp.Click(`body > div.L3eUgb > div.o3j99.ikrT4e.om7nvf > form > div:nth-child(1) > div.A8SBwf > div.RNNXgb > div > div.a4bIc > input`, chromedp.ByQuery),
		//TODO: Random nike product selector
		chromedp.SendKeys(`body > div.L3eUgb > div.o3j99.ikrT4e.om7nvf > form > div:nth-child(1) > div.A8SBwf > div.RNNXgb > div > div.a4bIc > input`, "site:nike.com/t/ "+randomProd()+kb.Enter, chromedp.ByQuery),
		chromedp.Sleep(time.Duration(3 * time.Second)),
		chromedp.Click(`h3`, chromedp.ByQuery),
		chromedp.Sleep(time.Duration(rand.Intn(6)) * time.Second),
	}
}

//////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////
func nikeSignupTask(favoriteBtn string, joinUsBtn string, genderBtn string) chromedp.Tasks {
	//TODO: Pull emails from .txt file
	//TODO: Gender Selector
	//TODO: Add scrolling/mouse movement
	// chromedp.ActionFunc(func(ctx context.Context) error {
	// 	_, exp, err := runtime.Evaluate(`window.getElementByID("shippingPickup").scrollIntoView({behavior: smooth});`).Do(ctx)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	if exp != nil {
	// 		return exp
	// 	}
	// 	return nil
	// }),
	// chromedp.ActionFunc(func(ctx context.Context) error {
	// 	_, exp, err := runtime.Evaluate(`window.scroll({top: -200, behavior:'smooth'});`).Do(ctx)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	if exp != nil {
	// 		return exp
	// 	}
	// 	return nil
	// }),
	return chromedp.Tasks{

		//Wait for product page to become visible
		chromedp.Navigate("https://www.nike.com/t/metcon-7-x-training-shoes-0l6Psg"),
		chromedp.WaitVisible(`#hf_header_label_copyright`, chromedp.ByID),
		//Click Favorite Button
		chromedp.Click(`#floating-atc-wrapper > div > button.wishlist-btn.ncss-btn-secondary-dark.btn-lg.mt3-sm`, chromedp.ByQuery),
		chromedp.WaitVisible(`#nike-unite-login-view > header > div.nike-unite-swoosh`, chromedp.ByQuery),

		//Enters Email in Field
		print("Inputting email"),
		chromedp.SendKeys(`//*[@placeholder="Email address"]`, randomdata.Email(), chromedp.BySearch),

		//Clicks Join Us Button
		print("Clicking 'Join Us' Button"),
		chromedp.Click(`.loginJoinLink.current-member-signin > a`, chromedp.ByQuery),

		//Enters Form Data
		print("Inputting password"),
		chromedp.SendKeys(`[placeholder="Password"]`, password.MustGenerate(15, 3, 3, false, false), chromedp.BySearch),
		print("Inputting First Name"),
		chromedp.SendKeys(`[placeholder="First Name"]`, randomdata.FirstName(randomdata.RandomGender), chromedp.BySearch),
		print("Inputting Last Name"),
		chromedp.SendKeys(`[placeholder="Last Name"]`, randomdata.LastName(), chromedp.BySearch),
		print("Inputting DOB"),
		chromedp.SendKeys(`[placeholder="Date of Birth"]`, monthToDigit(randomdata.Month())+fmt.Sprint(randomdata.Number(2))+fmt.Sprint(randomdata.Number(9))+fmt.Sprint(rand.Intn(45)+1960), chromedp.BySearch),

		//Clicks Male Gender Button
		print("Clicking Gender"),
		chromedp.Click(`li:nth-child(1) > input[type="button"]`, chromedp.ByQuery),

		print("Clicking 'Create Account'"),
		chromedp.Click(`[value="JOIN US"]`, chromedp.BySearch),
		chromedp.Sleep(1000 * time.Second),
	}
}

func nikeGoToPhoneNumber() chromedp.Tasks {
	return chromedp.Tasks{

		chromedp.WaitVisible(`#hf_header_label_copyright`, chromedp.ByID),

		//CLICK ACCOUNT SETTINGS
		chromedp.Click(`#root > div > div > div.main-layout > div > header > div.d-sm-h.d-lg-b > section > div > ul > li.member-nav-item.d-sm-ib.va-sm-m > div > div > ul > li:nth-child(1)`, chromedp.ByQuery),

		//CLICK ADD PHONE NUMBER
		chromedp.WaitVisible(`#mobile-container > div > div > form > div.account-form > div.mex-mobile-input-wrapper.ncss-col-sm-12.ncss-col-md-12.pl0-sm.pr0-sm.pb3-sm > div > div > div > div.ncss-col-sm-6.ta-sm-r.va-sm-m.flx-jc-sm-fe.d-sm-iflx > button`, chromedp.ByQuery),
		chromedp.Click(`#mobile-container > div > div > form > div.account-form > div.mex-mobile-input-wrapper.ncss-col-sm-12.ncss-col-md-12.pl0-sm.pr0-sm.pb3-sm > div > div > div > div.ncss-col-sm-6.ta-sm-r.va-sm-m.flx-jc-sm-fe.d-sm-iflx > button`, chromedp.ByQuery),

		//CLICK ADD PHONE NUMBER
		chromedp.Click(`div.sendCode > div.mobileNumber-div > input`, chromedp.ByQuery),
		chromedp.SendKeys(`div.sendCode > div.mobileNumber-div > input`, randomdata.Email(), chromedp.ByQuery),
	}
}

func nikeInputPhoneNumber(number string) chromedp.Tasks {
	return chromedp.Tasks{

		//SEND CODE
		chromedp.SendKeys(`div.sendCode > div.mobileNumber-div > input`, number, chromedp.ByQuery),

		//CLICK CODE BOX
		chromedp.Click(`#nike-unite-progressiveForm > div > div > input[type="button"]`, chromedp.ByQuery),
	}
}

func nikeConfirmPhoneNumber(code string) chromedp.Tasks {
	return chromedp.Tasks{

		//CLICK CODE BOX
		chromedp.WaitVisible(`input[type="number"]`, chromedp.BySearch),
		chromedp.Click(`input[type="number"]`, chromedp.BySearch),

		//SEND CODE
		chromedp.SendKeys(`#input[type="number"]`, code, chromedp.BySearch),

		//CLICK CONTINUE BOX
		chromedp.Click(`#nike-unite-progressiveForm > div > input[type="button"]`, chromedp.ByQuery),
	}
}

//////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////
func GetSMSToken() BearerResponse {
	tvURL := "https://www.textverified.com/api/"
	client := http.Client{}
	req, err := http.NewRequest("POST", tvURL+"SimpleAuthentication", nil)
	if err != nil {
		panic(err)
	}

	req.Header = http.Header{
		"X-SIMPLE-API-ACCESS-TOKEN": []string{"1_l8T2o8bPsP252roHDXhtO-tRzX3tqROlPzzam8kTaj7YPRlnccMekpuAmRsDh4r9_H13sjOr"},
	}

	res, err := client.Do(req)
	if err != nil {
		panic(err)
	} else {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}
		//bodyString := string(bodyBytes)
		var bear BearerResponse
		if err := json.Unmarshal(bodyBytes, &bear); err != nil { // Parse []byte to go struct pointer
			fmt.Println("Can not unmarshal JSON")
		}
		return bear
	}
}
func OrderNewNumber(bear BearerResponse) VerificationObject {

	url := "https://www.textverified.com/api/Verifications"
	method := "POST"

	payload := strings.NewReader(`{` + "\n" + `"id": 53` + "\n" + `}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		panic(err)
	}
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", "Bearer "+bear.BearerToken)

	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	var numRes VerificationObject
	if err := json.Unmarshal(body, &numRes); err != nil { // Parse []byte to go struct pointer
		fmt.Println("Can not unmarshal JSON")
	}

	return numRes
}
func CheckExistingNumber(bear BearerResponse, num VerificationObject) string {

	url := "https://www.textverified.com/api/Verifications/" + num.ID
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		panic(err)
	}
	req.Header.Add("Authorization", "Bearer "+bear.BearerToken)

	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	var numRes VerificationObject
	if err := json.Unmarshal(body, &numRes); err != nil { // Parse []byte to go struct pointer
		fmt.Println("Can not unmarshal JSON")
	}

	if numRes.Code != "null" {
		return numRes.Code
	} else {
		time.Sleep(time.Duration(rand.Intn(120)) * time.Second)
		return CheckExistingNumber(bear, num)
	}
}

//////////////////////////////////////////////////////////////////////////////////////////
func testTask() chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate("https://detect.azerpas.com"),
		chromedp.Sleep(time.Duration(10 * time.Second)),
	}
}
func loadProxies() []Proxy {

	var proxies []Proxy

	// open file
	f, err := os.Open("proxies.txt")
	if err != nil {
		panic(err)
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
		panic(err)
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(proxies), func(i, j int) { proxies[i], proxies[j] = proxies[j], proxies[i] })
	return proxies
}

//TODO:
// - MAKE THREADED
// - BEGIN BROWSER CONFIGURATION (TLS?)
// - FIX TypeWord FUNCTION
//
//
// FURTHER TODO'S LISTED BELOW

func runTasks(proxy Proxy) {

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

	//TODO: FIX ALL THIS SHIT
	fmt.Println("Setting up browser")
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

	// create chromedp's context
	fmt.Println("Creating context")
	parentCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	ctx, cancel := chromedp.NewContext(parentCtx)
	// ctx, cancel = chromedp.NewContext(ctx, chromedp.WithDebugf(log.Printf))
	defer cancel()

	//create a timeout
	//TODO: ADD ADJUSTABILITY
	fmt.Println("Creating timeout")
	ctx, cancel = context.WithTimeout(ctx, 150*time.Second)
	defer cancel()

	fmt.Println("Authenticating proxy")
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
	fmt.Println("Beginning New Task: Google")
	var nodes []*cdp.Node
	//var productLinks []cdp.NodeID

	err := chromedp.Run(ctx,
		fetch.Enable().WithHandleAuthRequests(true),
		//testTask(),
		googleTask(),
		print("Google - Complete"),
		chromedp.Nodes("a", &nodes),
	)
	if err != nil {
		panic(err)
	}

	// for _, n := range nodes {
	// 	if strings.Contains(n.AttributeValue("href"), "https://www.nike.com/t/") {
	// 		fmt.Println(n.AttributeValue("href"))
	// 		productLinks = append(productLinks, n.NodeID)
	// 	}
	// }

	// rand.Seed(time.Now().UnixNano())
	// randIdx := rand.Intn(len(productLinks))
	// randNodeID := productLinks[randIdx]

	fmt.Println("Beginning New Task: Nike")
	err = chromedp.Run(ctx,
		fetch.Enable().WithHandleAuthRequests(true),
		nikeSignupTask(favoriteBtn, registerBtn, genderBtn),
		print("Signup Complete"),
		nikeGoToPhoneNumber(),
	)
	if err != nil {
		panic(err)
	}

	token := GetSMSToken()
	order := OrderNewNumber(token)
	err = chromedp.Run(ctx,
		fetch.Enable().WithHandleAuthRequests(true),
		nikeInputPhoneNumber(string(order.Number)))
	if err != nil {
		panic(err)
	}

	code := CheckExistingNumber(token, order)
	err = chromedp.Run(ctx,
		fetch.Enable().WithHandleAuthRequests(true),
		nikeConfirmPhoneNumber(code))
	if err != nil {
		panic(err)
	}
}

func main() {

	proxies := loadProxies()

	//for i in len(proxies):

	rand.Seed(time.Now().UnixNano())
	randIdx := rand.Intn(len(proxies))

	randProx := proxies[randIdx]
	//proxies = append(proxies[:randIdx], proxies[randIdx+1:]...)

	fmt.Println("Proxy being used: " + randProx.IP)

	runTasks(randProx)
}
