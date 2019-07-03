package models

import (
	"fmt"
	"github.com/alexkupriyanov/KODE-Blog/api/utils"
	"github.com/badoux/goscraper"
	"github.com/joho/godotenv"
	"io/ioutil"
	"math"
	"mime"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Message struct {
	Id       uint `gorm:"primary_key"`
	AuthorID uint
	Author   User
	Text     string
	Media    []Media
	Link     Link
	Likes    int
}

type MessageViewListModel struct {
	Id     uint              `json:"id,omitempty"`
	Author string            `json:"author,omitempty"`
	Text   string            `json:"text,omitempty"`
	Media  uint              `json:"mediaCount,omitempty"`
	Link   LinkViewListModel `json:"link,omitempty"`
	Likes  int               `json:"likes,omitempty"`
}

type MessageViewDetailsModel struct {
	Id     uint                 `json:"id,omitempty"`
	Author string               `json:"author,omitempty"`
	Text   string               `json:"text,omitempty"`
	Media  []MediaOutput        `json:"media,omitempty"`
	Link   LinkViewDetailsModel `json:"link,omitempty"`
	Likes  int                  `json:"likes"`
}

type Like struct {
	MessageId uint
	UserId    uint
}

func (m *Message) Create(r *http.Request) map[string]interface{} {
	if m.Link.Link == "" && m.Text == "" && len(m.Media) == 0 {
		return utils.Message(false, "Empty post")
	}
	resp := m.Author.Get()
	if resp["status"] == false {
		return utils.Message(false, "Token malformed")
	}
	m.AuthorID = m.Author.ID
	m.Likes = 0
	GetDB().Create(&m)
	if m.Link.Link != "" {
		link := Link{MessageID: m.Id, Link: m.Link.Link}
		s, err := goscraper.Scrape(link.Link, 5)
		if err != nil {
			utils.Message(false, "Can't generate preview")
		}
		link.Description = s.Preview.Description
		if link.Description == "" {
			link.Description = s.Preview.Title
		}
		link.Preview = s.Preview.Images[0]
		GetDB().Create(&link)
		m.Link = link
	}
	m.Media = AddFiles(r)
	for _, v := range m.Media {
		v.MessageID = m.Id
		GetDB().Create(&v)
	}
	return nil
}

func GetMessageList(page uint) []MessageViewListModel {
	var messages []Message
	var result []MessageViewListModel
	if page <= 0 {
		return result
	}
	e := godotenv.Load()
	if e != nil {
		fmt.Print(e)
	}
	messageCount, _ := strconv.Atoi(os.Getenv("message_per_list"))
	GetDB().Preload("Media").Preload("Link").Limit(messageCount).Offset(uint(math.Max(float64(page-1), 0)) * uint(messageCount)).Find(&messages)
	for _, v := range messages {
		result = append(result, v.ToListModel())
	}
	return result
}

func (m *Message) Details() map[string]interface{} {
	GetDB().Preload("Link").Preload("Media").Preload("Author").First(&m)
	if m.AuthorID <= 0 {
		return utils.Message(false, "Message not found")
	}
	return nil
}

func (m *Message) Delete(token string) map[string]interface{} {
	user := User{Token: token}
	resp := user.Get()
	if resp["status"] == false {
		return utils.Message(false, "Malformed token")
	}
	GetDB().Preload("Link").Preload("Media").Preload("Author").First(&m)
	if m.AuthorID <= 0 {
		return utils.Message(false, "Message not found")
	}
	if m.AuthorID != user.ID {
		return utils.Message(false, "You haven't permission for this action")
	}
	for _, v := range m.Media {
		GetDB().Delete(&v)
	}
	GetDB().Delete(&m.Link)
	var c int
	GetDB().Model(&Like{}).Where("Message_id = ? AND User_id = ?", m.Id, user.ID).Count(&c)
	if c != 0 {
		GetDB().Delete(&Like{MessageId: m.Id, UserId: user.ID})
	}
	GetDB().Delete(&m)
	return nil
}

func (m *Message) ToListModel() MessageViewListModel {
	var output MessageViewListModel
	output.Id = m.Id
	output.Likes = m.Likes
	output.Link = m.Link.ToListModel()
	output.Author = m.Author.Username
	output.Text = m.Text
	output.Media = uint(len(m.Media))
	return output
}

func (m *Message) ToDetailsModel() MessageViewDetailsModel {
	var output MessageViewDetailsModel
	output.Id = m.Id
	output.Likes = m.Likes
	output.Link = m.Link.ToDetailsModel()
	output.Author = m.Author.Username
	output.Text = m.Text
	for _, v := range m.Media {
		output.Media = append(output.Media, v.ToOutputModel())
	}
	return output
}

func (m *Media) ToOutputModel() MediaOutput {
	var output MediaOutput
	output.Id = m.Id
	output.Link = os.Getenv("address") + ":" + os.Getenv("port") + "/" + m.Link
	return output
}

func (l *Link) ToListModel() LinkViewListModel {
	var output LinkViewListModel
	output.Link = l.Link
	return output
}

func (l *Link) ToDetailsModel() LinkViewDetailsModel {
	var output LinkViewDetailsModel
	output.Link = l.Link
	output.Preview = l.Preview
	output.Description = l.Description
	return output
}

func (m *Message) Like(token string) map[string]interface{} {
	user := User{Token: token}
	GetDB().First(&user)
	if user.ID <= 0 {
		return utils.Message(false, "Malformed token")
	}
	GetDB().Preload("Link").Preload("Media").Preload("Author").First(&m)
	if m.AuthorID == 0 {
		return utils.Message(false, "Message not found")
	}
	var c int
	GetDB().Model(&Like{}).Where("Message_id = ? AND User_id = ?", m.Id, user.ID).Count(&c)
	if c == 0 {
		GetDB().Create(&Like{MessageId: m.Id, UserId: user.ID})
		m.Likes++
	} else {
		GetDB().Delete(&Like{MessageId: m.Id, UserId: user.ID})
		m.Likes--
	}
	GetDB().Save(&m)
	return nil
}

func AddFiles(r *http.Request) []Media {
	var result []Media
	files := r.MultipartForm.File["file"]
	for _, file := range files {
		f, _ := file.Open()
		ext, _ := mime.ExtensionsByType(file.Header.Get("Content-Type"))
		tempFile, err := ioutil.TempFile("files", "file_*"+ext[0])
		if err != nil {
			fmt.Println(err)
			return nil
		}
		defer tempFile.Close()

		fileBytes, err := ioutil.ReadAll(f)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		_, _ = tempFile.Write(fileBytes)
		result = append(result, Media{Link: strings.Replace(tempFile.Name(), "\\", "/", -1)})
	}
	return result
}
