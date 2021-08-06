package main

import (
	"encoding/csv"
	"log"
	"os"
)

//Open `dbFile` and look for first row (account) with given `email` or `token`
//If account was found return that and true
//In other case return empty user and false
func FindByEmailOrToken(email, token string) (User, bool) {
	usersDb, err := os.OpenFile(dbFile, os.O_RDONLY, 0644)
	if err != nil {
		log.Fatalln(err)
	}

	defer usersDb.Close()

	csvLines, err := csv.NewReader(usersDb).ReadAll()
	if err != nil {
		log.Fatal(err)
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
	return err
}
