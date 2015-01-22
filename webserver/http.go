package main

import (
	"api"
	"fmt"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
	"os"
)

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_AES_KEY")))

func main() {
	//API ACTIONS
	http.HandleFunc("/regaccount", api.RegAccount)
	http.HandleFunc("/logaccount", api.LogAccount)
	http.HandleFunc("/codeverify", api.CodeVerify)
	err := http.ListenAndServe(":8889", nil)
	if err != nil {
		fmt.Println("Server did not start")
		log.Println(err)
		os.Exit(1)
	}
}
