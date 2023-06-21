package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

type LoginCredentials struct {
	ID        int    `json:"ID"`
	GrantType string `json:"grant_type"`
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/token", CreateToken).Methods("POST")

	log.Fatal(http.ListenAndServe(":8000", router))
}

func CreateToken(w http.ResponseWriter, r *http.Request) {
	var credentials LoginCredentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tokenString, err := GenerateToken(strconv.Itoa(credentials.ID), credentials.GrantType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(tokenString))
}

func GenerateToken(id, grant_type string) (string, error) {
	claims := jwt.MapClaims{
		"id":         id,
		"grant_type": grant_type,
	}

	secretKey := "MySecretKeyIsSecretSoDoNotTell"

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
