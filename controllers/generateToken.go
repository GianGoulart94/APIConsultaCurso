package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-oci8"
	_ "github.com/qodrorid/godaemon"
	"os"
	"crypto/md5"
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type Token struct{
	Token 	string `json:"token"`
}

func GenerateToken(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r);
	var token Token
	login := vars["app"];
	
	hash, err := bcrypt.GenerateFromPassword([]byte(login), bcrypt.DefaultCost)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Hash to store:", string(hash))

    hasher := md5.New()
	hasher.Write(hash)
	token.Token =fmt.Sprintf("%s %s", "Basic",hex.EncodeToString(hasher.Sum(nil)))
	
	
	
	w.Header().Set("Content-Type", "application/json")

	os.Setenv("TOKEN_API",token.Token)


	json.NewEncoder(w).Encode(token)
}