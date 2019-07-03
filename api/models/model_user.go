package models

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"os"
	"time"
)

// User object
type User struct {
	ID       uint   `gorm:"primary_key"`
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

func (u *User) Create() error {
	if u.Username == "" || u.Password == "" {
		return errors.New("Username or password are empty")
	}
	var user User
	GetDB().First(&user, User{Username: u.Username})
	if user.ID > 0 {
		return errors.New("User already exist")
	}
	u.Password = getMD5Hash(u.Password)
	GetDB().Create(u)
	if u.ID <= 0 {
		return errors.New("Failed to create account, connection error")
	}
	return nil
}

func (u *User) Login() (string,error) {
	e := godotenv.Load()
	if e != nil {
		fmt.Print(e)
	}
	var user User
	GetDB().First(&user, User{Username: u.Username, Password: getMD5Hash(u.Password)})
	if user.ID <= 0 {
		return "", errors.New("User not found")
	}
	tokenLifeTime := time.Now().Add(time.Hour * 24).Unix()
	tk := jwt.StandardClaims{
		ExpiresAt: tokenLifeTime,
	}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	user.Token = tokenString //Store the token in the response
	GetDB().Save(&user)
	return user.Token, nil
}

func (u *User) Logout() error {
	err := u.Get()
	if err != nil{
		return err
	}
	GetDB().Model(&u).Update("token", gorm.Expr("NULL"))
	return nil
}

func (u *User) Get() error {
	GetDB().First(&u, User{Token: u.Token})
	if u.ID <= 0 {
		return errors.New("User not found")
	}
	return nil
}

func getMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
