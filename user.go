package main

import (
	"crypto/sha512"
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"time"
)

type User struct {
	Email        string
	PasswordHash string
	Token        string //Random string generated at registation
}

//Hashing with salt
//Return 44-symbols string (base64)
func PasswordHash(password string) string {
	salt := "Yeeh_zMVk3"
	hasher := sha512.New512_256()
	hasher.Write([]byte(password + salt))
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

//Open `dbFile` and look for first row (account) with given `email` or `token`
//If account was found return that and true
//In other case return empty user and false
func FindByEmailOrToken(email, token string) (User, bool) {
	usersDb, err := os.OpenFile(dbFile, os.O_RDONLY, 0644)
	if err != nil {
		return User{}, false
	}

	defer usersDb.Close()

	csvLines, err := csv.NewReader(usersDb).ReadAll()
	if err != nil {
		fmt.Println(err)
	}

	//Read csv file line-by-line
	for _, line := range csvLines {
		u := User{
			Email:        line[0],
			PasswordHash: line[1],
			Token:        line[2],
		}
		if u.Email == email || u.Token == token {
			return u, true
		}
	}
	return User{}, false
}

func FindByEmail(email string) (User, bool) {
	return FindByEmailOrToken(email, "")
}

func FindByToken(token string) (User, bool) {
	return FindByEmailOrToken("", token)
}

//Add `u` to other users in `dbFile`
func AppendUser(u User) error {
	usersDb, err := os.OpenFile(dbFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer usersDb.Close()

	_, err = usersDb.Write([]byte(u.Email + "," + u.PasswordHash + "," + u.Token + "\n"))
	if err != nil {
		return err
	}
	return nil
}

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func IsEmailValid(email string) bool {
	return len(email) >= 3 && len(email) <= 255 && emailRegex.MatchString(email)
}

//Add new user to database
//If successfull return (200, nil)
func UserRegister(email, pass string) (httpStatus int, err error) {
	if email == "" {
		return http.StatusBadRequest, fmt.Errorf("Missed email")
	}

	if !IsEmailValid(email) {
		return http.StatusBadRequest, fmt.Errorf("Incorrect email")
	}

	if pass == "" {
		return http.StatusBadRequest, fmt.Errorf("Password can't be empty")
	}

	if _, has := FindByEmail(email); has {
		return http.StatusNotAcceptable, fmt.Errorf("Email already used")
	}

	var u User
	u.Email = email
	u.PasswordHash = PasswordHash(pass)
	u.Token = RandomString(12)

	if err := AppendUser(u); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

//Find token for registered user in database
func UserLogin(email, pass string) (token string, httpStatus int, err error) {
	if email == "" {
		return "", http.StatusBadRequest, fmt.Errorf("Incorrect email")
	}

	u, has := FindByEmail(email)
	if !has || u.PasswordHash != PasswordHash(pass) {
		return "", http.StatusUnauthorized, fmt.Errorf("Incorrect login")
	}

	return u.Token, http.StatusOK, nil
}

//Check an existing user with this token
func IsAvaiableToken(token string) bool {
	if token == "" {
		return false
	}

	_, has := FindByToken(token)
	return has
}
