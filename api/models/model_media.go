package models

type MediaOutput struct {
	Id   uint   `json:"Id"`
	Link string `json:"Link"`
}

type Media struct {
	Id        uint `gorm:"primary_key"`
	Link      string
	MessageID uint
	Message   *Message
}
