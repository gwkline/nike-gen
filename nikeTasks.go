package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/chromedp/chromedp"
	"github.com/sethvargo/go-password/password"
)

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
