package models

import (
	"KODE-Blog/api/utils"
	"crypto/md5"
	"encoding/hex"
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
	Role     string `json:"role,omitempty"`
	Token    string `json:"token"`
}

func (u *User) Create() map[string]interface{} {
	if u.Username == "" || u.Password == "" {
		return utils.Message(false, "Username or password are empty")
	}
	if u.Role != "admin" {
		u.Role = "user"
	}
	var user User
	GetDB().First(&user, User{Username: u.Username})
	if user.ID > 0 {
		return utils.Message(false, "User already exist")
	}
	u.Password = getMD5Hash(u.Password)
	GetDB().Create(u)
	if u.ID <= 0 {
		return utils.Message(false, "Failed to create account, connection error")
	}
	return utils.Message(true, "Account created")
}

func (u *User) Login() map[string]interface{} {
	e := godotenv.Load()
	if e != nil {
		fmt.Print(e)
	}
	var user User
	GetDB().First(&user, User{Username: u.Username, Password: getMD5Hash(u.Password)})
	if user.ID <= 0 {
		return utils.Message(false, "User not found")
	}
	tokenLifeTime := time.Now().Add(time.Hour * 24).Unix()
	tk := jwt.StandardClaims{
		ExpiresAt: tokenLifeTime,
	}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	user.Token = tokenString //Store the token in the response
	GetDB().Save(&user)
	return map[string]interface{}{"token": user.Token}
}

func (u *User) Logout() map[string]interface{} {
	resp := u.Get()
	if resp["status"] == false {
		return resp
	}
	GetDB().Model(&u).Update("token", gorm.Expr("NULL"))
	return map[string]interface{}{}
}

func (u *User) Get() map[string]interface{} {
	GetDB().First(&u, User{Token: u.Token})
	if u.ID <= 0 {
		return utils.Message(false, "User not found")
	}
	return map[string]interface{}{}
}

func getMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
