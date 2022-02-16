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
//							CUSTOM OBJECTS AND STRUCTURES								//
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
//						HELPER FUNCTIONS AND OTHER UTILITIES							//
//////////////////////////////////////////////////////////////////////////////////////////
func typeWord(sel interface{}, word string, opts func(*chromedp.Selector), ctx context.Context) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {

			runeList := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
			fmt.Println(word)
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

func randomProd() string {
	rand.Seed(time.Now().UnixNano())
	num := rand.Intn(9)
	switch num {
	case 0:
		fmt.Println("shirt")
		return "shirt"
	case 1:
		fmt.Println("shoes")
		return "shoes"
	case 3:
		fmt.Println("pants")
		return "pants"
	case 4:
		fmt.Println("socks")
		return "socks"
	case 5:
		fmt.Println("jacket | DIDNT WORK")
		return "jacket"
	case 6:
		fmt.Println("fleece")
		return "fleece"
	case 7:
		fmt.Println("sweatshirt")
		return "sweatshirt"
	case 8:
		fmt.Println("sweat")
		return "sweat"
	case 9:
		fmt.Println("dri fit")
		return "dri fit"
	}
	fmt.Println("clothes | DIDNT WORK")
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

func loadEmails() []string {

	var emails []string
	// open file
	f, err := os.Open("emails.txt")
	if err != nil {
		panic(err)
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
		panic(err)
	}

	return emails
}

//////////////////////////////////////////////////////////////////////////////////////////
//									GOOGLE TASKS										//
//////////////////////////////////////////////////////////////////////////////////////////
func googleTask(ctx context.Context) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate("https://www.google.com"),
		chromedp.Sleep(time.Duration(3 * time.Second)),
		chromedp.WaitVisible(`#Mses6b`),
		chromedp.Click(`body > div.L3eUgb > div.o3j99.ikrT4e.om7nvf > form > div:nth-child(1) > div.A8SBwf > div.RNNXgb > div > div.a4bIc > input`, chromedp.ByQuery),
		//TODO: Random nike product selector
		typeWord(`body > div.L3eUgb > div.o3j99.ikrT4e.om7nvf > form > div:nth-child(1) > div.A8SBwf > div.RNNXgb > div > div.a4bIc > input`, "site:nike.com/t/ "+randomProd()+kb.Enter, chromedp.ByQuery, ctx),
		chromedp.Sleep(time.Duration(3 * time.Second)),
		chromedp.Click(`h3`, chromedp.ByQuery),
		chromedp.Sleep(time.Duration(rand.Intn(6)) * time.Second),
	}
}

//////////////////////////////////////////////////////////////////////////////////////////
//									LOGIN TASKS											//
//////////////////////////////////////////////////////////////////////////////////////////
func nikeSignupTask(favoriteBtn string, joinUsBtn string, genderBtn string, ctx context.Context) chromedp.Tasks {
	//TODO: Add scrolling/mouse movement

	gender := randomdata.RandomGender
	fmt.Print(gender)
	firstName := randomdata.FirstName(gender)
	lastName := randomdata.LastName()

	rand.Seed(time.Now().UnixNano())
	num := rand.Intn(2)
	email := "test"
	switch num {
	case 0:
		email = firstName + lastName + fmt.Sprint(rand.Intn(999)) + "@gmail.com"
	case 1:
		email = firstName + lastName + fmt.Sprint(rand.Intn(99)) + "@gmail.com"
	case 2:
		email = firstName + lastName + fmt.Sprint(rand.Intn(9)) + "@gmail.com"
	}

	switch gender {
	case 0:
		genderBtn = `li:nth-child(1) > input[type="button"]`
	case 1:
		genderBtn = `li:nth-child(2) > input[type="button"]`
	}

	return chromedp.Tasks{

		chromedp.WaitVisible(`#hf_header_label_copyright`, chromedp.ByID),
		//move mouse around
		//slowly scroll to bottom
		//mouse move
		randDelay(10, ctx),

		//move mouse
		//scroll halfway up
		//mouse move to Target
		chromedp.ScrollIntoView(`#RightRail > div > div.prl6-sm.prl0-lg > div > span > details`, chromedp.ByID),
		chromedp.Click(`#RightRail > div > div.prl6-sm.prl0-lg > div > span > details`, chromedp.ByQuery),
		//or scroll to/click product details
		//LONGER pause

		//Click Favorite Button
		chromedp.ScrollIntoView(`#floating-atc-wrapper > div > button.wishlist-btn.ncss-btn-secondary-dark.btn-lg.mt3-sm`, chromedp.ByQuery),
		//move mouse
		chromedp.Click(`#floating-atc-wrapper > div > button.wishlist-btn.ncss-btn-secondary-dark.btn-lg.mt3-sm`, chromedp.ByQuery),

		//Clicks Join Us Button
		chromedp.WaitVisible(`#nike-unite-login-view > header > div.nike-unite-swoosh`, chromedp.ByQuery),
		//move mouse
		chromedp.Click(`.loginJoinLink.current-member-signin > a`, chromedp.ByQuery),

		//Enters Form Data
		//random slight scroll

		//mouse move + click
		//delay
		typeWord(`//*[@placeholder="Email address"]`, email, chromedp.BySearch, ctx),

		//tab (+ med delay) OR mouse move + click (+ high delay)
		typeWord(`[placeholder="Password"]`, password.MustGenerate(15, 3, 3, false, false), chromedp.BySearch, ctx),

		//tab (+ med delay) OR mouse move + click (+ high delay)
		typeWord(`[placeholder="First Name"]`, firstName, chromedp.BySearch, ctx),

		//tab (+ med delay) OR mouse move + click (+ high delay)
		typeWord(`[placeholder="Last Name"]`, lastName, chromedp.BySearch, ctx),

		//tab (+ med delay) OR mouse move + click (+ high delay)
		//TODO: FIX typeWord for DOB, make higher than avg delays, different backspace mechanics
		chromedp.SendKeys(`[placeholder="Date of Birth"]`, monthToDigit(randomdata.Month())+fmt.Sprint(randomdata.Number(2))+fmt.Sprint(randomdata.Number(9))+fmt.Sprint(rand.Intn(45)+1960), chromedp.BySearch),

		//Clicks Gender Button
		//delay + mouse move
		chromedp.Click(genderBtn, chromedp.ByQuery),

		//delay + mouse move
		chromedp.Click(`[value="JOIN US"]`, chromedp.BySearch),
		chromedp.Sleep(1000 * time.Second),
	}
}

func nikeGoToPhoneNumber(ctx context.Context) chromedp.Tasks {
	return chromedp.Tasks{

		chromedp.WaitVisible(`#hf_header_label_copyright`, chromedp.ByID),

		//CLICK ACCOUNT SETTINGS
		chromedp.Click(`#root > div > div > div.main-layout > div > header > div.d-sm-h.d-lg-b > section > div > ul > li.member-nav-item.d-sm-ib.va-sm-m > div > div > ul > li:nth-child(1)`, chromedp.ByQuery),

		//CLICK ADD PHONE NUMBER
		chromedp.WaitVisible(`#mobile-container > div > div > form > div.account-form > div.mex-mobile-input-wrapper.ncss-col-sm-12.ncss-col-md-12.pl0-sm.pr0-sm.pb3-sm > div > div > div > div.ncss-col-sm-6.ta-sm-r.va-sm-m.flx-jc-sm-fe.d-sm-iflx > button`, chromedp.ByQuery),
		chromedp.Click(`#mobile-container > div > div > form > div.account-form > div.mex-mobile-input-wrapper.ncss-col-sm-12.ncss-col-md-12.pl0-sm.pr0-sm.pb3-sm > div > div > div > div.ncss-col-sm-6.ta-sm-r.va-sm-m.flx-jc-sm-fe.d-sm-iflx > button`, chromedp.ByQuery),

		//CLICK ADD PHONE NUMBER
		chromedp.Click(`div.sendCode > div.mobileNumber-div > input`, chromedp.ByQuery),
		typeWord(`div.sendCode > div.mobileNumber-div > input`, randomdata.Email(), chromedp.ByQuery, ctx),
	}
}

func nikeInputPhoneNumber(number string, ctx context.Context) chromedp.Tasks {
	return chromedp.Tasks{

		//SEND CODE
		typeWord(`div.sendCode > div.mobileNumber-div > input`, number, chromedp.ByQuery, ctx),

		//CLICK CODE BOX
		chromedp.Click(`#nike-unite-progressiveForm > div > div > input[type="button"]`, chromedp.ByQuery),
	}
}

func nikeConfirmPhoneNumber(code string, ctx context.Context) chromedp.Tasks {
	return chromedp.Tasks{

		//CLICK CODE BOX
		chromedp.WaitVisible(`input[type="number"]`, chromedp.BySearch),
		chromedp.Click(`input[type="number"]`, chromedp.BySearch),

		//SEND CODE
		typeWord(`#input[type="number"]`, code, chromedp.BySearch, ctx),

		//CLICK CONTINUE BOX
		chromedp.Click(`#nike-unite-progressiveForm > div > input[type="button"]`, chromedp.ByQuery),
	}
}

//////////////////////////////////////////////////////////////////////////////////////////
//								SMS VERIFICATION TASKS									//
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
//								TASK ENGINE AND MAIN									//
//////////////////////////////////////////////////////////////////////////////////////////
//TODO:
// - MAKE THREADED
// - BEGIN BROWSER CONFIGURATION (TLS?)
// -
//
// FURTHER TODO'S LISTED BELOW

func runTasks(proxy Proxy, email string) {

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

	//create a timeout
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

	err := chromedp.Run(ctx,
		fetch.Enable().WithHandleAuthRequests(true),
		testTask(),
		chromedp.Sleep(10*time.Second),
		googleTask(ctx),
		chromedp.Nodes("a", &nodes),
	)
	if err != nil {
		panic(err)
	}
	for _, n := range nodes {
		if strings.Contains(n.AttributeValue("href"), "https://www.nike.com/t/") {
			fmt.Println(n.AttributeValue("href"))
			productLinks = append(productLinks, n.AttributeValue("href"))
		}
	}

	rand.Seed(time.Now().UnixNano())
	randIdx := rand.Intn(len(productLinks))
	randURL := productLinks[randIdx]
	searchString := `[href="` + randURL + `"]`
	fmt.Println(searchString)

	fmt.Println("Beginning New Task: Nike")
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

func main() {

	proxies := loadProxies()
	emails := loadEmails()

	//for i in len(proxies):

	rand.Seed(time.Now().UnixNano())
	randIdx := rand.Intn(len(proxies))

	randProx := proxies[randIdx]
	email := emails[0]
	//proxies = append(proxies[:randIdx], proxies[randIdx+1:]...)
	//emails = append(emails[:i], emails[i+1:]...)

	fmt.Println("")
	fmt.Println("Proxy being used: " + randProx.IP)
	fmt.Println("Email being used: " + email)
	fmt.Println("")

	runTasks(randProx, email)
}
