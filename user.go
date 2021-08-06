package main

import (
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/http"
	"regexp"
	"time"
)

//Hashing with saltFindByEmailOrToken
//Return 44-symbols string (base64)
func PasswordHash(password string) string {
	hasher := sha512.New512_256()
	hasher.Write([]byte(password + PasswordSalt))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

func RandomString(length int) string {
	var letters = []rune("0123456789abcdefghijklmnopqrstuvwxyz")

	rand.Seed(time.Now().UnixNano())
	s := make([]rune, length)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func CreateToken() string {
	return RandomString(TokenLength)
}

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func IsEmailValid(email string) bool {
	return len(email) >= 3 && len(email) <= 255 && emailRegex.MatchString(email)
}

//Add new user to database
//If successful return (200, nil)
func UserRegister(email, pass string) (httpStatus int, err error) {
	if email == "" {
		return http.StatusBadRequest, fmt.Errorf("missed email")
	}

	if !IsEmailValid(email) {
		return http.StatusBadRequest, fmt.Errorf("incorrect email")
	}

	if pass == "" {
		return http.StatusBadRequest, fmt.Errorf("password can't be empty")
	}

	if _, has := FindByEmail(email); has {
		return http.StatusNotAcceptable, fmt.Errorf("email already used")
	}

	u := User{
		Email:        email,
		PasswordHash: PasswordHash(pass),
		Token:        CreateToken(),
	}

	if err := AppendUser(u); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

//Request new token for registered user
func UserLogin(email, pass string) (token string, httpStatus int, err error) {
	if email == "" {
		return "", http.StatusBadRequest, fmt.Errorf("incorrect email")
	}

	u, has := FindByEmail(email)
	if !has || u.PasswordHash != PasswordHash(pass) {
		return "", http.StatusUnauthorized, fmt.Errorf("incorrect login")
	}

	newToken := CreateToken()
	err = RewriteToken(email, newToken)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}
	return newToken, http.StatusOK, nil
}

//Check an existing user with this token
func IsAvailableToken(token string) bool {
	if token == "" {
		return false
	}

	_, has := FindByToken(token)
	return has
}
