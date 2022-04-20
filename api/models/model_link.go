package models

type Link struct {
	Id          uint   `gorm:"primary_key"`
	Link        string `json:"link,omitempty"`
	Preview     string `json:"preview,omitempty"`
	Description string `json:"description,omitempty"`
	MessageID   uint
	Message     *Message
}

type LinkViewDetailsModel struct {
	Link        string `json:"link,omitempty"`
	Preview     string `json:"preview,omitempty"`
	Description string `json:"description,omitempty"`
}

type LinkViewListModel struct {
	Link string `json:"link"`
}
