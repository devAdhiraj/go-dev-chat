package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rivo/tview"
)

const ApiBaseUrl = "http://localhost:8085"

type MsgType string

const (
	MsgTypeText  MsgType = "text"
	MsgTypeImage MsgType = "image"
	MsgTypeVideo MsgType = "video"
	MsgTypeAudio MsgType = "audio"
	MsgTypeFile  MsgType = "file"
	MsgTypeDir   MsgType = "dir"
	MsgTypeCtrl  MsgType = "ctrl"
)

type View string

const (
	start   View = "start"
	signup  View = "signup"
	login   View = "login"
	chats   View = "chats"
	newChat View = "new-chat"
	chat    View = "chat"
	loading View = "loading"
)

type Msg struct {
	ID         uint      `json:"id"`
	MsgValue   string    `json:"msgValue"`
	SenderID   uint      `json:"senderId"`
	ReceiverID uint      `json:"receiverId"`
	MsgType    MsgType   `json:"msgType"`
	Timestamp  time.Time `json:"timestamp"`
	SeenAt     time.Time `json:"seenAt"`
}

type Chats struct {
	Msg
	SenderUsername   string `json:"senderUsername"`
	ReceiverUsername string `json:"receiverUsername"`
}

type ChatInfo struct {
	Messages []Msg
	username string
	userId   uint
}

type App struct {
	httpClient   *http.Client
	app          *tview.Application
	conn         *websocket.Conn
	pages        *tview.Pages
	startView    *tview.Flex
	loginView    *tview.Flex
	signupView   *tview.Flex
	chatsView    *tview.List
	chatView     *tview.Flex
	loadingView  *tview.TextView
	newChatView  *tview.Flex
	authInfo     AuthInfo
	chatList     []Chats
	currentChat  ChatInfo
	mu           sync.Mutex
	outgoingMsgs chan Msg
	quit         chan bool
}

type User struct {
	ID       uint   `json:"ID"`
	Username string `json:"username"`
}
type AuthInfo struct {
	currentUser User
	token       string
}

type ErrorResponse interface {
	getError() string
}

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token   string `json:"token"`
	Message string `json:"message"`
	UserID  uint   `json:"userId"`
}

var views = []View{
	start,
	login,
	signup,
	chats,
	newChat,
	chat,
	loading,
}

func (a *AuthResponse) getError() string {
	return a.Message
}

func (a *App) makeAPIRequest(method string, endpoint string, useAuth bool, body any, response any) error {
	requestBody, err := json.Marshal(body)
	if err != nil {
		return errors.New("error encoding data: " + err.Error())
	}
	req, err := http.NewRequest(method, ApiBaseUrl+endpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		return errors.New("Error creating request: " + err.Error())
	}
	if useAuth {
		req.Header.Add("Authorization", "Bearer "+a.authInfo.token)
	}
	resp, err := a.httpClient.Do(req)
	if err != nil {
		return errors.New("error making request")
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return errors.New("server error: invalid response from server when calling " + endpoint + "\nError: " + err.Error())
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return errors.New(resp.Status + " - " + response.(ErrorResponse).getError())
	}
	return nil
}

func (a *App) hideAllViews() {
	for _, v := range views {
		a.pages.HidePage(string(v))
	}
}

func (a *App) showView(view View) {
	v := string(view)
	a.hideAllViews()
	a.pages.SendToFront(v)
	a.pages.ShowPage(v)
}
