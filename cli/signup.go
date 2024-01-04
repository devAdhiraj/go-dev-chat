package main

import (
	"net/http"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const SignupEndpoint = "/signup"

func (a *App) initSignupView() {
	a.signupView = tview.NewFlex().SetDirection(tview.FlexRow).AddItem(
		tview.NewForm().
			AddInputField("Username:", "", 30, nil, nil).
			AddPasswordField("Password:", "", 30, '*', nil).
			AddPasswordField("Confirm Password:", "", 30, '*', nil).
			AddButton("Signup", a.handleSignup).
			AddButton("Cancel", a.showStartView),
		0, 3, true,
	).AddItem(
		tview.NewTextView().SetTextColor(tcell.ColorRed),
		0, 1, false,
	)
	a.pages.AddPage("signup", a.signupView, true, false)
}

func (a *App) showSignupView() {
	signupForm := a.signupView.GetItem(0).(*tview.Form)
	signupForm.GetFormItem(0).(*tview.InputField).SetText("")
	signupForm.GetFormItem(1).(*tview.InputField).SetText("")
	signupForm.GetFormItem(2).(*tview.InputField).SetText("")
	a.showView(signup)
}

func (a *App) setSignupError(errorMsg string) {
	a.signupView.GetItem(1).(*tview.TextView).SetText("Error: " + errorMsg)
}

func (a *App) handleSignup() {
	signupForm := a.signupView.GetItem(0).(*tview.Form)
	username := signupForm.GetFormItem(0).(*tview.InputField).GetText()
	password := signupForm.GetFormItem(1).(*tview.InputField).GetText()
	confirm := signupForm.GetFormItem(2).(*tview.InputField).GetText()
	if username == "" {
		a.setSignupError("Empty username")
		return
	}
	if password != confirm {
		a.setSignupError("Passwords don't match")
		signupForm.GetFormItem(1).(*tview.InputField).SetText("")
		signupForm.GetFormItem(2).(*tview.InputField).SetText("")
		return
	}
	signupRequest := &AuthRequest{Username: username, Password: password}
	var resp AuthResponse
	err := a.makeAPIRequest(http.MethodPost, SignupEndpoint, false, signupRequest, &resp)
	if err != nil {
		a.setSignupError(err.Error())
		a.setStartViewText("signup successful!", tcell.ColorGreen)
		a.showSignupView()
		return
	}
	a.showStartView()
}
