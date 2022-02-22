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
func nikeSignupTask(ctx context.Context, task *Task) chromedp.Tasks {
	//TODO: Add scrolling/mouse movement, better email system
	rand.Seed(time.Now().UnixNano())
	task.Gender = rand.Intn(1)
	task.First_Name = randomdata.FirstName(task.Gender)
	task.Last_Name = randomdata.LastName()
	genderBtn := "test"
	month, err := monthToDigit(randomdata.Month())
	if err != nil {
		error.Error(err)
	}

	task.DOB[0] = month
	task.DOB[1] = fmt.Sprint(randomdata.Number(9))
	task.DOB[2] = fmt.Sprint(rand.Intn(45) + 1960)
	fmt.Println(fmt.Sprintf(task.DOB[2]))

	rand.Seed(time.Now().UnixNano())
	num := rand.Intn(2)
	if !USE_EMAIL_LIST {
		switch num {
		case 0:
			task.Email = task.First_Name + task.Last_Name + fmt.Sprint(rand.Intn(999)) + "@gmail.com"
		case 1:
			task.Email = task.First_Name + task.Last_Name + fmt.Sprint(rand.Intn(99)) + "@gmail.com"
		case 2:
			task.Email = task.First_Name + task.Last_Name + fmt.Sprint(rand.Intn(9)) + "@gmail.com"
		}
	}

	switch task.Gender {
	case 0:
		genderBtn = genderBtnMale
	case 1:
		genderBtn = genderBtnFemale
	}

	return chromedp.Tasks{

		logTask(task, "Beginning Sign Up"),
		logTask(task, "Waiting For Page Load"),
		chromedp.WaitVisible(waitPageNike1, chromedp.ByID),
		logTask(task, "Page Loaded"),
		//move mouse around
		//slowly scroll to bottom
		//mouse move
		randDelay(ctx),

		//move mouse
		//scroll halfway up
		//mouse move to Target
		logTask(task, "Random Scrolling"),
		chromedp.ScrollIntoView(detailsClick, chromedp.ByID),
		randDelay(ctx),
		chromedp.Click(detailsClick, chromedp.ByQuery),
		//or scroll to/click product details
		randDelay(ctx),

		//Click Favorite Button
		logTask(task, "Clicking Favorite Button"),
		chromedp.ScrollIntoView(favoriteBtn, chromedp.ByQuery),
		randDelay(ctx),
		//move mouse
		chromedp.Click(favoriteBtn, chromedp.ByQuery),
		randDelay(ctx),

		//Clicks Join Us Button
		logTask(task, "Clicking Join Us Button"),
		chromedp.WaitVisible(waitPageNike2, chromedp.ByQuery),
		//move mouse
		chromedp.Click(joinUsBtn, chromedp.ByQuery),
		randDelay(ctx),

		//Enters Form Data
		//random slight scroll

		//mouse move + click
		//delay
		logTask(task, "Entering Email"),
		typeWord(emailInput, task.Email, chromedp.BySearch, ctx),
		randDelay(ctx),

		//tab (+ med delay) OR mouse move + click (+ high delay)
		logTask(task, "Entering Password"),
		typeWord(passwordInput, password.MustGenerate(15, 3, 3, false, false), chromedp.BySearch, ctx),
		randDelay(ctx),

		//tab (+ med delay) OR mouse move + click (+ high delay)
		logTask(task, "Entering First Name"),
		typeWord(firstNameInput, task.First_Name, chromedp.BySearch, ctx),
		randDelay(ctx),

		//tab (+ med delay) OR mouse move + click (+ high delay)
		logTask(task, "Entering Last Name"),
		typeWord(lastNameInput, task.Last_Name, chromedp.BySearch, ctx),
		randDelay(ctx),

		//tab (+ med delay) OR mouse move + click (+ high delay)
		//TODO: FIX typeWord for DOB, make higher than avg delays, different backspace mechanics
		logTask(task, "Entering DOB"),

		chromedp.SendKeys(DOBInput, task.DOB[0], chromedp.BySearch),
		randDelay(ctx),

		chromedp.SendKeys(DOBInput, task.DOB[1], chromedp.BySearch),
		randDelay(ctx),

		chromedp.SendKeys(DOBInput, task.DOB[2], chromedp.BySearch),
		randDelay(ctx),

		//Clicks Gender Button
		//delay + mouse move
		logTask(task, "Choosing Gender"),
		chromedp.Click(genderBtn, chromedp.ByQuery),
		randDelay(ctx),

		//delay + mouse move
		logTask(task, "Clicking Sign Up"),
		chromedp.Click(signUpButton, chromedp.BySearch),

		logTask(task, "Signup Complete"),
	}
}

func nikeGoToPhoneNumber(ctx context.Context, task *Task) chromedp.Tasks {
	return chromedp.Tasks{

		logTask(task, "Navigating To Settings"),
		logTask(task, "Waiting For Page Load"),
		chromedp.WaitVisible(waitPageNike1, chromedp.ByID),
		logTask(task, "Page Loaded"),

		//CLICK ACCOUNT SETTINGS
		logTask(task, "Clicking Account Settings"),
		chromedp.Click(settingsButton, chromedp.ByQuery),

		//CLICK ADD PHONE NUMBER
		logTask(task, "Clicking Add Phone Number"),
		chromedp.WaitVisible(addPhoneBtn, chromedp.ByQuery),
		chromedp.Click(addPhoneBtn, chromedp.ByQuery),
	}
}

func nikeInputPhoneNumber(number string, ctx context.Context, task *Task) chromedp.Tasks {
	return chromedp.Tasks{

		//SEND CODE
		logTask(task, "Inputting Phone Number"),
		typeWord(phoneNumberInput, number, chromedp.ByQuery, ctx),

		//CLICK CODE BOX
		logTask(task, "Requesting Code"),
		chromedp.Click(codeBoxButton, chromedp.ByQuery),
	}
}

func nikeConfirmPhoneNumber(code string, ctx context.Context, task *Task) chromedp.Tasks {
	return chromedp.Tasks{

		//CLICK CODE BOX
		logTask(task, "Clicking Code Entry Box"),
		chromedp.WaitVisible(enterTheValueBtn, chromedp.BySearch),
		chromedp.Click(enterTheValueBtn, chromedp.BySearch),

		//SEND CODE
		logTask(task, "Inputting Code"),
		typeWord(enterTheValueBtn, code, chromedp.BySearch, ctx),

		//CLICK CONTINUE BOX
		logTask(task, "Saving"),
		chromedp.Click(continueButtonBox, chromedp.ByQuery),
	}
}
