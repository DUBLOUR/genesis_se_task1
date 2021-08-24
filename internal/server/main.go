package server

import (
	"encoding/json"
	"github.com/dublour/genesis_se_task1/pkg/binance"
	"github.com/dublour/genesis_se_task1/pkg/model"
	"log"
	"net/http"
	"os"
)

func Respond(w http.ResponseWriter, r *http.Request, httpStatus int, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(data)

	//Log requests and response data
	log.Println("REQ:", r.URL.String(), "\nSTATUS:", httpStatus, "BODY:", data)
}

func main() {
	//Init log file
	logger, err := os.OpenFile(LogFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()
	log.SetOutput(logger)

	//Common case for incorrect endpoints
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		mess := map[string]interface{}{
			"status": "Fail",
		}
		Respond(w, r, http.StatusNotFound, mess)
	})

	//Registration of new account
	//Require `email` and `password` fields in URL
	//Append new user in `UsersFile` and return `{"status":"Ok"}` if successful
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
	//Return `{"status":"Ok","token":"############"}` if successful
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
	//Return `{"status":"Ok","BTCUAH":$$$$$$}` if successful
	http.HandleFunc("/btcRate", func(w http.ResponseWriter, r *http.Request) {
		var urlToken, headerToken, token string
		urlToken = r.URL.Query().Get("token")
		headerToken = r.Header.Get("X-API-Key")

		if headerToken != "" {
			token = headerToken
		} else {
			token = urlToken
		}

		if IsAvailableToken(token) { //Find user in database with that token
			cost, err := Cost("BTCUAH")

			if err == nil {
				mess := map[string]interface{}{
					"status": "Ok",
					"BTCUAH": cost,
				}
				Respond(w, r, http.StatusOK, mess)
			} else {
				mess := map[string]interface{}{
					"status": err.Error(),
				}
				Respond(w, r, http.StatusInternalServerError, mess)
			}
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

	http.ListenAndServe(ServerPort, nil)
}
