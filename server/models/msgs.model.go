package models

import (
	"time"
)

type MsgType string

const (
	MsgTypeText  MsgType = "text"
	MsgTypeImage MsgType = "image"
	MsgTypeVideo MsgType = "video"
	MsgTypeAudio MsgType = "audio"
	MsgTypeFile  MsgType = "file"
	MsgTypeDir   MsgType = "dir"
)

type Msg struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	MsgValue   string    `gorm:"column:msg_value" json:"msgValue"`
	SenderID   uint      `gorm:"sender_id" json:"senderId"`
	ReceiverID uint      `gorm:"receiver_id" json:"receiverId"`
	MsgType    MsgType   `gorm:"column:msg_type" json:"msgType"`
	Timestamp  time.Time `gorm:"timestamp" json:"timestamp"`
	SeenAt     time.Time `gorm:"seen_at" json:"seenAt"`
}

type Chats struct {
	Msg
	SenderUsername   string `gorm:"sender_username" json:"senderUsername"`
	ReceiverUsername string `gorm:"receiver_username" json:"receiverUsername"`
}

func (Msg) TableName() string {
	return "msgs"
}

func (m *Msg) Save() error {
	m.Timestamp = time.Now()
	r := db.Create(m)
	return r.Error
}

func MarkAsSeen(msgId uint, seenTimestamp time.Time) error {
	return db.Model(&Msg{}).Where("timestamp < ?", db.Model(&Msg{}).Where("id = ?", msgId).Where("seen_at = NULL").Select("timestamp")).Update("seen_at", seenTimestamp).Error
}

func GetUserConvo(userId1 uint, userId2 uint, limit int, offset uint) ([]Msg, error) {
	var convoMsgs []Msg
	if err := db.Where("(receiver_id = ? AND sender_id = ?) OR (receiver_id = ? AND sender_id = ?)", userId1, userId2, userId2, userId1).
		Order("timestamp asc").Offset(int(offset)).
		Limit(limit).Find(&convoMsgs).Error; err != nil {
		return nil, err
	}
	return convoMsgs, nil
}

func GetUserChats(userId uint) ([]Chats, error) {
	var chats []Chats

	subquery := db.Model(&Msg{}).
		Select("DISTINCT ON (GREATEST(sender_id, receiver_id), LEAST(sender_id, receiver_id)) *").
		Where("receiver_id = ?", userId).
		Or("sender_id = ?", userId).
		Order("GREATEST(sender_id, receiver_id), LEAST(sender_id, receiver_id), timestamp DESC")

	if err := db.Table("(?) as m", subquery).
		Select("m.id, msg_value, msg_type, sender_id, receiver_id, timestamp, seen_at, s.username as sender_username, r.username as receiver_username").
		Joins("inner join users as s on s.id = m.sender_id").
		Joins("inner join users as r on r.id = m.receiver_id").
		Order("GREATEST(m.sender_id, m.receiver_id), LEAST(m.sender_id, m.receiver_id), m.timestamp DESC").
		Find(&chats).Error; err != nil {
		return nil, err
	}
	return chats, nil
}
