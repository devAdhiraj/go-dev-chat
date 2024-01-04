package main

import (
	"net/http"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const LoginEndpoint = "/login"

func (a *App) setLoginError(errorMsg string) {
	a.loginView.GetItem(1).(*tview.TextView).SetText("Error: " + errorMsg)
}

func (a *App) initLoginView() {
	a.loginView =
		tview.NewFlex().SetDirection(tview.FlexRow).AddItem(
			tview.NewForm().
				AddInputField("Username:", "", 30, nil, nil).
				AddPasswordField("Password:", "", 30, '*', nil).
				AddButton("Login", a.handleLogin).
				AddButton("Cancel", a.showStartView),
			0, 3, true,
		).AddItem(
			tview.NewTextView().SetTextColor(tcell.ColorRed),
			0, 1, false,
		)
	a.pages.AddPage("login", a.loginView, true, false)
}

func (a *App) handleLogin() {
	loginForm := a.loginView.GetItem(0).(*tview.Form)
	loginReq := &AuthRequest{
		Username: loginForm.GetFormItem(0).(*tview.InputField).GetText(),
		Password: loginForm.GetFormItem(1).(*tview.InputField).GetText(),
	}
	if loginReq.Password == "" {
		a.setLoginError("missing password")
		return
	}
	a.showLoadingView("logging you in...")
	var r AuthResponse
	err := a.makeAPIRequest(http.MethodPost, LoginEndpoint, false, loginReq, &r)
	if err != nil {
		a.setLoginError(err.Error())
		a.showLoginView()
		return
	}
	a.authInfo.currentUser = User{Username: loginReq.Username, ID: r.UserID}
	a.authInfo.token = r.Token
	if a.conn != nil {
		a.conn.Close()
	}
	if err := a.startSocketConn(); err != nil {
		a.setLoginError((err.Error()))
		a.showLoginView()
		return
	}
	go a.msgReceiver()
	go a.msgSender()

	a.showChatsView()
}

func (a *App) showLoginView() {
	loginForm := a.loginView.GetItem(0).(*tview.Form)
	loginForm.GetFormItem(0).(*tview.InputField).SetText("")
	loginForm.GetFormItem(1).(*tview.InputField).SetText("")
	a.showView(login)
}
