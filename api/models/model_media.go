package models

type MediaOutput struct {
	Id   uint   `json:"id"`
	Link string `json:"link"`
}

type Media struct {
	Id        uint `gorm:"primary_key"`
	Link      string
	MessageID uint
	Message   *Message
}
