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
func nikeSignupTask(ctx context.Context, tid string) chromedp.Tasks {
	//TODO: Add scrolling/mouse movement, better email system
	rand.Seed(time.Now().UnixNano())
	gender := rand.Intn(1)
	firstName := randomdata.FirstName(gender)
	lastName := randomdata.LastName()
	genderBtn := "test"
	DOBString := monthToDigit(randomdata.Month()) + fmt.Sprint(randomdata.Number(2)) + fmt.Sprint(randomdata.Number(9)) + fmt.Sprint(rand.Intn(45)+1960)

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
		genderBtn = genderBtnMale
	case 1:
		genderBtn = genderBtnFemale
	}

	return chromedp.Tasks{
		//chromedp.Navigate("https://www.nike.com/t/metcon-7-x-training-shoes-0l6Psg/CZ8281-883"),
		print("Task ID: " + tid + " | Beginning Sign Up"),
		print("Task ID: " + tid + " | Waiting For Page Load"),
		chromedp.WaitVisible(waitPageNike1, chromedp.ByID),
		print("Task ID: " + tid + " | Page Loaded"),
		//move mouse around
		//slowly scroll to bottom
		//mouse move
		randDelay(ctx),

		//move mouse
		//scroll halfway up
		//mouse move to Target
		print("Task ID: " + tid + " | Random Scrolling"),
		chromedp.ScrollIntoView(detailsClick, chromedp.ByID),
		randDelay(ctx),
		chromedp.Click(detailsClick, chromedp.ByQuery),
		//or scroll to/click product details
		randDelay(ctx),

		//Click Favorite Button
		print("Task ID: " + tid + " | Clicking Favorite Button"),
		chromedp.ScrollIntoView(favoriteBtn, chromedp.ByQuery),
		randDelay(ctx),
		//move mouse
		chromedp.Click(favoriteBtn, chromedp.ByQuery),
		randDelay(ctx),

		//Clicks Join Us Button
		print("Task ID: " + tid + " | Clicking Join Us Button"),
		chromedp.WaitVisible(waitPageNike2, chromedp.ByQuery),
		//move mouse
		chromedp.Click(joinUsBtn, chromedp.ByQuery),
		randDelay(ctx),

		//Enters Form Data
		//random slight scroll

		//mouse move + click
		//delay
		print("Task ID: " + tid + " | Entering Email"),
		typeWord(emailInput, email, chromedp.BySearch, ctx),
		randDelay(ctx),

		//tab (+ med delay) OR mouse move + click (+ high delay)
		print("Task ID: " + tid + " | Entering Password"),
		typeWord(passwordInput, password.MustGenerate(15, 3, 3, false, false), chromedp.BySearch, ctx),
		randDelay(ctx),

		//tab (+ med delay) OR mouse move + click (+ high delay)
		print("Task ID: " + tid + " | Entering First Name"),
		typeWord(firstNameInput, firstName, chromedp.BySearch, ctx),
		randDelay(ctx),

		//tab (+ med delay) OR mouse move + click (+ high delay)
		print("Task ID: " + tid + " | Entering Last Name"),
		typeWord(lastNameInput, lastName, chromedp.BySearch, ctx),
		randDelay(ctx),

		//tab (+ med delay) OR mouse move + click (+ high delay)
		//TODO: FIX typeWord for DOB, make higher than avg delays, different backspace mechanics
		print("Task ID: " + tid + " | Entering DOB"),
		chromedp.SendKeys(DOBInput, DOBString, chromedp.BySearch),
		randDelay(ctx),

		//Clicks Gender Button
		//delay + mouse move
		print("Task ID: " + tid + " | Choosing Gender"),
		chromedp.Click(genderBtn, chromedp.ByQuery),
		randDelay(ctx),

		//delay + mouse move
		print("Task ID: " + tid + " | Clicking Sign Up"),
		chromedp.Click(signUpButton, chromedp.BySearch),

		print("Task ID: " + tid + " | Signup Complete"),
	}
}

func nikeGoToPhoneNumber(ctx context.Context, tid string) chromedp.Tasks {
	return chromedp.Tasks{

		print("Task ID: " + tid + " | Navigating To Settings"),
		print("Task ID: " + tid + " | Waiting For Page Load"),
		chromedp.WaitVisible(waitPageNike1, chromedp.ByID),
		print("Task ID: " + tid + " | Page Loaded"),

		//CLICK ACCOUNT SETTINGS
		print("Task ID: " + tid + " | Clicking Account Settings"),
		chromedp.Click(settingsButton, chromedp.ByQuery),

		//CLICK ADD PHONE NUMBER
		print("Task ID: " + tid + " | Clicking Add Phone Number"),
		chromedp.WaitVisible(addPhoneBtn, chromedp.ByQuery),
		chromedp.Click(addPhoneBtn, chromedp.ByQuery),
	}
}

func nikeInputPhoneNumber(number string, ctx context.Context, tid string) chromedp.Tasks {
	return chromedp.Tasks{

		//SEND CODE
		print("Task ID: " + tid + " | Inputting Phone Number"),
		typeWord(phoneNumberInput, number, chromedp.ByQuery, ctx),

		//CLICK CODE BOX
		print("Task ID: " + tid + " | Requesting Code"),
		chromedp.Click(codeBoxButton, chromedp.ByQuery),
	}
}

func nikeConfirmPhoneNumber(code string, ctx context.Context, tid string) chromedp.Tasks {
	return chromedp.Tasks{

		//CLICK CODE BOX
		print("Task ID: " + tid + " | Clicking Code Entry Box"),
		chromedp.WaitVisible(enterTheValueBtn, chromedp.BySearch),
		chromedp.Click(enterTheValueBtn, chromedp.BySearch),

		//SEND CODE
		print("Task ID: " + tid + " | Inputting Code"),
		typeWord(enterTheValueBtn, code, chromedp.BySearch, ctx),

		//CLICK CONTINUE BOX
		print("Task ID: " + tid + " | Saving"),
		chromedp.Click(continueButtonBox, chromedp.ByQuery),
	}
}
