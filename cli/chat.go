package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (a *App) initChatView() {
	chatTextView := tview.NewTextView().
		SetTextAlign(tview.AlignLeft).
		SetDynamicColors(true).
		SetText("")

	chatTextView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			a.app.SetFocus(a.pages)
		}
		return event
	})

	var inputField *tview.InputField
	inputField = tview.NewInputField().
		SetLabel("Type your message: ").
		SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEnter {
				message := inputField.GetText()
				if message != "" {
					msg := Msg{MsgValue: message, MsgType: MsgTypeText, ReceiverID: a.currentChat.userId, SenderID: a.authInfo.currentUser.ID}
					a.mu.Lock()
					a.currentChat.Messages = append(a.currentChat.Messages, msg)
					a.addMessage("You", message, "blue", time.Now().Format(time.DateTime))
					a.mu.Unlock()
					a.outgoingMsgs <- msg
					inputField.SetText("")
				}
			}
			if key == tcell.KeyTab || key == tcell.KeyUp {
				a.app.SetFocus(a.chatView.GetItem(2))
			}
		})

	backButton := tview.NewButton("Back")
	backButton.SetInputCapture(
		func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyEnter {
				a.showChatsView()
			}
			if event.Key() == tcell.KeyTab {
				a.app.SetFocus(a.chatView.GetItem(0))
			}
			return event
		})

	chatTextView.SetScrollable(true)
	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(chatTextView, 0, 1, false).
		AddItem(inputField, 2, 1, true).
		AddItem(tview.NewFlex().AddItem(nil, 0, 1, false).AddItem(backButton, 10, 1, true).AddItem(nil, 2, 1, false), 1, 1, false)

	a.chatView = flex
	flex.SetBorder(true).SetBorderColor(tcell.ColorWhite)
	a.pages.AddPage(string(chat), a.chatView, true, false)
}

func (a *App) addMessage(username string, message string, color string, timestamp string) {
	chatTextView := a.chatView.GetItem(0).(*tview.TextView)
	formattedUsername := tview.Escape(fmt.Sprintf("[%s]", username))
	chatTextView.SetText(fmt.Sprintf("%s\n[%s]%s %s:[-] %s", chatTextView.GetText(false), color, timestamp, formattedUsername, message))
	chatTextView.ScrollToEnd()
}

type MsgResponse struct {
	Messages []Msg  `json:"messages"`
	Message  string `json:"message"`
}

func (m *MsgResponse) getError() string {
	return m.Message
}

func (a *App) showChatView() {
	var msgResp MsgResponse
	if err := a.makeAPIRequest(http.MethodGet, fmt.Sprintf("/msgs/%d", a.currentChat.userId), true, nil, &msgResp); err != nil {
		a.setStartViewText("Error: "+err.Error(), tcell.ColorRed)
		a.showStartView()
		return
	}

	a.chatView.SetTitle("Chat with " + a.currentChat.username)

	a.mu.Lock()
	a.currentChat.Messages = msgResp.Messages
	a.updateChatView()
	a.mu.Unlock()
}

func (a *App) updateChatView() {
	a.chatView.GetItem(0).(*tview.TextView).SetText("")
	for _, msg := range a.currentChat.Messages {
		if a.authInfo.currentUser.ID == msg.SenderID {
			a.addMessage("You", msg.MsgValue, "blue", msg.Timestamp.Format(time.DateTime))
		} else {
			a.addMessage(a.currentChat.username, msg.MsgValue, "green", msg.Timestamp.Format(time.DateTime))
		}
	}
	a.showView(chat)
}
