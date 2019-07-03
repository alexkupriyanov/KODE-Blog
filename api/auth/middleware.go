package auth

import (
	"github.com/alexkupriyanov/KODE-Blog/api/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"net/http"
	"os"
)

func CheckToken(inputToken string, r *http.Request) error {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(inputToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("token_password")), nil
	})
	if err != nil {
		return err
	}
	var count int
	models.GetDB().Find(&models.User{Token:inputToken}).Count(&count)
	if !token.Valid || count == 0{
		return errors.New("Token not valid")
	}
	return nil
}