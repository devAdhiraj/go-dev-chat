package main

import (
	"net/http"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const ChatsEndpoint = "/chats"

type ChatsResponse struct {
	Chats   []Chats `json:"chats"`
	Message string  `json:"message"`
}

func (c *ChatsResponse) getError() string {
	return c.Message
}

func (a *App) initChatsView() {
	a.chatsView = tview.NewList().SetMainTextColor(tcell.ColorWhite).SetSecondaryTextColor(tcell.ColorWhite)
	a.pages.AddPage(string(chats), a.chatsView, true, false)
}

func (a *App) showChatsView() {
	var resp ChatsResponse
	err := a.makeAPIRequest(http.MethodGet, ChatsEndpoint, true, nil, &resp)
	if err != nil {
		a.setStartViewText("error fetching chats "+err.Error(), tcell.ColorRed)
		a.showStartView()
		return
	}
	a.mu.Lock()
	a.chatList = resp.Chats
	a.updateChatsView()
	a.mu.Unlock()
}

func (a *App) updateChatsView() {
	a.chatsView.Clear()
	a.chatsView.AddItem("[red]Quit[-]", "", 'q', func() {
		a.showStartView()
	})
	a.chatsView.AddItem("[green]+ New Chat[-]", "", 'n', func() {
		a.showNewChatView()
	})
	for i, chat := range a.chatList {
		username := chat.SenderUsername
		if username == a.authInfo.currentUser.Username {
			username = chat.ReceiverUsername
		}
		lastMsg := chat.MsgValue
		formattedUsername := username
		if chat.ReceiverID == a.authInfo.currentUser.ID && chat.SeenAt.IsZero() {
			formattedUsername = "[::b]" + username
			lastMsg = "[::b]" + lastMsg
		}
		userId := chat.SenderID
		if userId == a.authInfo.currentUser.ID {
			userId = chat.ReceiverID
		}
		a.chatsView.AddItem(formattedUsername, lastMsg, rune(i%26)+'a', func() {
			a.currentChat.username = username
			a.currentChat.userId = userId
			a.showChatView()
		})
	}
	a.showView(chats)
}
