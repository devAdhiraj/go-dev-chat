package main

import (
	"net/http"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (a *App) handleNewChatStart() {
	newChatForm := a.newChatView.GetItem(1).(*tview.Form)
	username := newChatForm.GetFormItem(0).(*tview.InputField).GetText()
	var usr User
	if err := a.makeAPIRequest(http.MethodGet, "/users/"+username, true, nil, &usr); err != nil {
		a.setStartViewText(err.Error(), tcell.ColorRed)
		a.showStartView()
		return
	}
	a.currentChat.username = usr.Username
	a.currentChat.userId = usr.ID
	a.showChatView()
}

func (a *App) initNewChatView() {
	a.newChatView =
		tview.NewFlex().SetDirection(tview.FlexRow).AddItem(
			tview.NewTextView().SetText("Enter user you want to start a chat with: "), 0, 1, false,
		).AddItem(
			tview.NewForm().
				AddInputField("Username:", "", 30, nil, nil).
				AddButton("Start chat", a.handleNewChatStart).
				AddButton("Cancel", a.showChatsView),
			0, 3, true,
		).AddItem(
			tview.NewTextView().SetTextColor(tcell.ColorRed),
			0, 1, false,
		)
	a.pages.AddPage("new-chat", a.newChatView, true, false)
}

func (a *App) showNewChatView() {
	newChatForm := a.newChatView.GetItem(1).(*tview.Form)
	newChatForm.GetFormItem(0).(*tview.InputField).SetText("")
	a.showView(newChat)
}
