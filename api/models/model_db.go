package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"os"
	"time"
)

var db *gorm.DB

func init() {
	username := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbName := os.Getenv("db_name")
	dbHost := os.Getenv("db_host")
	dbType := os.Getenv("db_type")

	dbUri := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true", username, password, dbHost, dbName)
	conn, err := gorm.Open(dbType, dbUri)
	for i := 0; err != nil && i < 5; i++ {
		time.Sleep(1 * time.Second)
		conn, err = gorm.Open(dbType, dbUri)
	}
	if err != nil {
		fmt.Print(err)
		time.Sleep(2 * time.Second)
		os.Exit(4060)
	}

	db = conn
	err = db.DB().Ping()
	if err != nil {
		fmt.Println(err)
		time.Sleep(2 * time.Second)
		os.Exit(4060)
	}
	db.Debug().AutoMigrate(&User{}, &Media{}, &Message{}, &Link{}, &Like{})
	var c int
	db.Model(&User{}).Count(&c)
	if c == 0 {
		db.Create(&User{Username: "admin", Password: getMD5Hash("admin")})
	}
}

func GetDB() *gorm.DB {
	return db
}
