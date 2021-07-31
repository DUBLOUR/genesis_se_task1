package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

const logFile string = "/var/log/btc-requester/responces.log"
const dbFile string = "/etc/btc-requester/users.csv"
const serverPort string = ":9990"

func Respond(w http.ResponseWriter, r *http.Request, httpStatus int, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(data)

	//Log requst and responce data
	log.Println("REQ:", r.URL.String(), "\nSTATUS:", httpStatus, "BODY:", data)
}

func main() {
	//Init log file
	loger, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer loger.Close()
	log.SetOutput(loger)

	//Common case for incorrect endpoints
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		mess := map[string]interface{}{
			"status": "Fail",
		}
		Respond(w, r, http.StatusNotFound, mess)
	})

	//Registration of new account
	//Require `email` and `password` fields in URL
	//Append new user in `dbFile` and return `{"status":"Ok"}` if successfull
	http.HandleFunc("/user/create", func(w http.ResponseWriter, r *http.Request) {
		email := r.URL.Query().Get("email")
		password := r.URL.Query().Get("password")

		httpStatus, err := UserRegister(email, password)

		status := "Ok"
		if err != nil {
			status = err.Error() //Detail about what went wrong
		}

		mess := map[string]interface{}{
			"status": status,
		}
		Respond(w, r, httpStatus, mess)
	})

	//Get an API-token for registered users
	//Require `email` and `password` fields in URL
	//Return `{"status":"Ok","token":"############"}` if successfull
	http.HandleFunc("/user/login", func(w http.ResponseWriter, r *http.Request) {
		email := r.URL.Query().Get("email")
		password := r.URL.Query().Get("password")

		token, httpStatus, err := UserLogin(email, password)

		mess := map[string]interface{}{}
		if err == nil {
			mess = map[string]interface{}{
				"status": "Ok",
				"token":  token,
			}
		} else {
			mess = map[string]interface{}{
				"status": err.Error(),
			}
		}
		Respond(w, r, httpStatus, mess)
	})

	//Get a current cost of BitCoin in UAH
	//Require `token` in URL on `X-API-Key` in Header (more priority)
	//Return `{"status":"Ok","BTCUAH":$$$$$$}` if successfull
	http.HandleFunc("/btcRate", func(w http.ResponseWriter, r *http.Request) {
		var urlToken, headerToken, token string
		urlToken = r.URL.Query().Get("token")
		headerToken = r.Header.Get("X-API-Key")

		if headerToken != "" {
			token = headerToken
		} else {
			token = urlToken
		}

		if IsAvaiableToken(token) { //Find user in database with that token
			cost, httpStatus, err := Cost("BTCUAH")

			mess := map[string]interface{}{}
			if err == nil {
				mess = map[string]interface{}{
					"status": "Ok",
					"BTCUAH": cost,
				}
			} else {
				mess = map[string]interface{}{
					"status": err.Error(),
				}
			}
			Respond(w, r, httpStatus, mess)
			return
		}

		if token == "" {
			mess := map[string]interface{}{
				"status": "Missing token",
			}
			Respond(w, r, http.StatusBadRequest, mess)
			return
		}

		mess := map[string]interface{}{
			"status": "Invalid token",
		}
		Respond(w, r, http.StatusForbidden, mess)

	})

	http.ListenAndServe(serverPort, nil)
}
