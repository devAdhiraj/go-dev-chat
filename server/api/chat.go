package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/devAdhiraj/go-dev-chat/server/models"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	TxtMsg = iota
	ImgMsg
	VidMsg
	AudioMsg
	FileMsg
	DirMsg
)

type WsConn struct {
	conn *websocket.Conn
	mq   chan models.Msg
}

var producerQueue chan models.Msg = make(chan models.Msg, 100)
var connUsers map[uint][]WsConn = make(map[uint][]WsConn)
var connUsersLock sync.Mutex

func ConsumerRoutine() {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "PLAINTEXT://broker:29092",
		"group.id":          "ws-server-consumer",
		"auto.offset.reset": "smallest",
	})
	if err != nil {
		fmt.Println("Error creating consumer -", err)
		return
	}
	defer consumer.Close()
	err = consumer.Subscribe("chat-messages", nil)
	if err != nil {
		fmt.Println("Error subscribing to topic -", err)
		return
	}
	for {
		ev := consumer.Poll(100)
		switch e := ev.(type) {
		case *kafka.Message:
			fmt.Println("event key value -", string(e.Key), string(e.Value))
			var msg models.Msg
			if err := json.Unmarshal(e.Value, &msg); err != nil {
				fmt.Println("Error decoding msg from mq -", err, e.Value)
				continue
			}
			SendUserMsg(msg)

		case kafka.Error:
			fmt.Println("Consumer Error: ", e)
			// default:
			// 	fmt.Println("timeout/default? -", e, ev)
		}
	}
}

func ProducerRoutine() {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "PLAINTEXT://broker:29092",
		"client.id":         "http://localhost:8085",
		"acks":              "all"})
	if err != nil {
		fmt.Println("Error creating producer -", err)
		return
	}
	deliveryChan := make(chan kafka.Event, 10000)
	topic := "chat-messages"
	for {
		msg := <-producerQueue
		jsonMsg, err := json.Marshal(msg)
		if err != nil {
			fmt.Println("error encoding msg to json -", err, msg)
			continue
		}
		if err := producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          jsonMsg, Key: []byte(fmt.Sprint("%u", msg.ID))},
			deliveryChan,
		); err != nil {
			fmt.Println("Error producing -", err)
		}
	}
}

func SendUserMsg(msg models.Msg) {
	connUsersLock.Lock()
	defer connUsersLock.Unlock()
	if connUsers[msg.ReceiverID] == nil {
		return
	}
	for _, connection := range connUsers[msg.ReceiverID] {
		connection.mq <- msg
	}
}

func addUserConn(userId uint, c *websocket.Conn) error {
	connUsersLock.Lock()
	defer connUsersLock.Unlock()
	if connUsr := connUsers[userId]; connUsr == nil {
		connUsers[userId] = make([]WsConn, 0)
	}
	connUsers[userId] = append(connUsers[userId], WsConn{conn: c, mq: make(chan models.Msg, 10)})
	return nil
}

func removeUserConn(userId uint, c *websocket.Conn) {
	connUsersLock.Lock()
	defer connUsersLock.Unlock()
	if connUsr := connUsers[userId]; connUsr == nil {
		return
	}
	if len(connUsers[userId]) == 1 {
		delete(connUsers, userId)
		return
	}
	idx := -1
	for i, wc := range connUsers[userId] {
		if wc.conn == c {
			idx = i
		}
	}
	if idx != -1 {
		connUsers[userId] = append(connUsers[userId][:idx], connUsers[userId][idx+1:]...)
	}
}

func WsHandler(c *gin.Context) {
	uid, exists := c.Get("user_id")
	if _, ok := uid.(uint); !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthenticated",
		})
		return
	}
	userId := uint(uid.(uint))

	upgrader := websocket.Upgrader{
		ReadBufferSize:  4096,
		WriteBufferSize: 4096,
	}

	upgrader.CheckOrigin = func(r *http.Request) bool { return true } // TODO: remove

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cannot create ws connection",
		})
		return
	}
	defer conn.Close()

	if err := addUserConn(userId, conn); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error storing connection",
		})
		return
	}
	defer removeUserConn(userId, conn)

	fmt.Println("socket connected from ", userId)
	quit := make(chan bool, 1)
	go ConnMsgSender(userId, conn, &quit)
	defer func(q *chan bool) {
		*q <- true
	}(&quit)
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil || (messageType != websocket.TextMessage && messageType != websocket.BinaryMessage) {
			fmt.Println("Error reading msg -", err, messageType)
			break
		}

		var msg models.Msg
		if err := json.Unmarshal(p, &msg); err != nil {
			fmt.Println("error decoding json -", err)
			break
		}
		msg.SenderID = userId
		if msg.ReceiverID == 0 || msg.MsgValue == "" || msg.MsgType == "" {
			e := conn.WriteJSON(models.Msg{SenderID: 0, ReceiverID: userId, MsgValue: "invalid msg payload", MsgType: "ctrl"})
			if e != nil {
				break
			}
			continue
		}
		if err := msg.Save(); err != nil {
			fmt.Println("Error saving msg -", err)
			e := conn.WriteJSON(models.Msg{SenderID: 0, ReceiverID: userId, MsgValue: "error saving", MsgType: "ctrl"})
			if e != nil {
				break
			}
			continue
		}
		producerQueue <- msg
	}
}

func ConnMsgSender(userId uint, conn *websocket.Conn, q *chan bool) {
	if connUsers[userId] == nil {
		return
	}
	var mq chan models.Msg = nil
	for _, c := range connUsers[userId] {
		if c.conn == conn {
			mq = c.mq
		}
	}
	if mq == nil {
		return
	}
	for {
		select {
		case msg := <-mq:
			if err := conn.WriteJSON(msg); err != nil {
				fmt.Println("Error sending msg -", err)
			}
		case <-*q:
			return
		}
	}
}
