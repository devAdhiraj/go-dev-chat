package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/fatih/color"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/gorilla/websocket"
)

var myApp *App

func (a *App) initStartView() {
	a.startView = tview.NewFlex().SetDirection(tview.FlexRow).AddItem(
		tview.NewForm().
			AddButton("Signup", a.showSignupView).
			AddButton("Login", a.showLoginView).
			AddButton("Quit", a.exit),
		0, 1, true,
	).AddItem(tview.NewTextView(), 0, 1, false)
	a.pages.AddPage("start", a.startView, true, false)
}

func (a *App) setStartViewText(msg string, color tcell.Color) {
	a.startView.GetItem(1).(*tview.TextView).SetTextColor(color).SetText(msg)
}

func (a *App) showStartView() {
	a.showView(start)
}

func (a *App) exit() {
	a.app.Stop()
}

func initApp() {
	myApp = &App{
		app:        tview.NewApplication(),
		pages:      tview.NewPages(),
		httpClient: &http.Client{},
	}
	myApp.initStartView()
	myApp.initLoginView()
	myApp.initSignupView()
	myApp.initChatsView()
	myApp.initChatView()
	myApp.initLoadingView()
	myApp.initNewChatView()
	myApp.outgoingMsgs = make(chan Msg, 10)
	myApp.quit = make(chan bool)
	myApp.app.SetRoot(myApp.pages, true)
}

func (a *App) startSocketConn() error {
	header := http.Header{}
	header.Add("Authorization", "Bearer "+myApp.authInfo.token)
	var err error
	myApp.conn, _, err = websocket.DefaultDialer.Dial("ws://localhost:8085/ws", header)
	if err != nil {
		a.setStartViewText("Error connecting. Can't create socket connection: "+err.Error(), tcell.ColorRed)
		return err
	}
	return nil
}

func (a *App) msgSender() {

	for {
		select {
		case msg := <-a.outgoingMsgs:
			if err := a.conn.WriteJSON(msg); err != nil {
				fmt.Fprintln(os.Stderr, "Error sending msg -", err)
			}
		case <-a.quit:
			return
		}
	}
}

func (a *App) msgReceiver() {
	for {
		messageType, p, err := a.conn.ReadMessage()
		if err != nil {
			break
		}
		if messageType != websocket.TextMessage && messageType != websocket.BinaryMessage {
			continue
		}
		var msg Msg
		if err := json.Unmarshal(p, &msg); err != nil {
			continue
		}
		if msg.ReceiverID == 0 || msg.MsgValue == "" || msg.MsgType == "" {
			continue
		}
		a.handleNewMsg(msg)
	}
	// TODO: handle disconnect?
}

func (a *App) handleNewMsg(msg Msg) {
	currentPage, _ := a.pages.GetFrontPage()
	if msg.MsgType != MsgTypeText {
		return
	}
	switch currentPage {
	case string(chats):
		found := false
		a.mu.Lock()
		for i, chat := range a.chatList {
			if (msg.SenderID == chat.SenderID && msg.ReceiverID == chat.ReceiverID) || (msg.ReceiverID == chat.SenderID && msg.SenderID == chat.ReceiverID) {
				a.chatList[i].Msg = msg
				a.updateChatsView()
				found = true
				break
			}
		}
		a.mu.Unlock()
		if found {
			break
		}
		var usr User
		if err := a.makeAPIRequest(http.MethodGet, fmt.Sprintf("/user/%d", msg.SenderID), true, nil, &usr); err != nil {
			return
		}
		newChat := Chats{Msg: msg, ReceiverUsername: a.authInfo.currentUser.Username, SenderUsername: usr.Username}
		a.mu.Lock()
		a.chatList = append(a.chatList, newChat)
		a.updateChatsView()
		a.mu.Unlock()
	case string(chat):
		a.mu.Lock()
		if a.currentChat.userId == msg.SenderID {
			found := false
			for i, m := range a.currentChat.Messages {
				if m.ID == msg.ID {
					a.currentChat.Messages[i] = msg
					found = true
					break
				}
			}
			if !found {
				a.currentChat.Messages = append(a.currentChat.Messages, msg)
			}
		}
		a.updateChatView()
		a.mu.Unlock()
	}
	a.app.Draw()
}

func main() {
	initApp()
	myApp.showStartView()

	if err := myApp.app.Run(); err != nil {
		fmt.Println(err)
	}
	color.New(color.FgBlue).Println("Exited go cli app")
}
