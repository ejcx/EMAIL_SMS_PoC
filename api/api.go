package api

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"smsmailer"
	"encoding/json"
	"math/rand"
	"model"
	"net/http"
	"regexp"
	"strconv"
	"time"
	"fmt"
)

type Error struct {
	Error   string
	Success int64
}
type RegisterResponse struct {
	RegisterResponse string
}
type LoginResponse struct {
	LoginResponse string
}
/* Not for production..... */
func gensalt() string {
	rand.Seed(time.Now().UTC().UnixNano())
	salt := ""
	for len(salt) < 50 {
		salt += strconv.Itoa(rand.Intn(65535))
	}
	return salt
}
func shatwo(usermaterial string) string {
	var h = sha256.New()
	h.Write([]byte(usermaterial))
	pwhash := hex.EncodeToString(h.Sum(nil))
	return pwhash
}
func getnewapikey(usermaterial string) string {
	pt1 := gensalt()
	pt2 := shatwo(usermaterial)
	apikeybig := shatwo(pt1 + pt2)
	//b64 it
	apikeyb64 := base64.StdEncoding.EncodeToString([]byte(apikeybig))
	apikeyb64 = apikeyb64[0:24]
	return "moo_" + apikeyb64

}

func registerout(w http.ResponseWriter, resp_msg []byte) {
	w.Header().Set("Content-Type", "application/json")
	jsonresp := RegisterResponse{RegisterResponse: string(resp_msg)}
	jsoncatted, err := json.Marshal(jsonresp)
	if err == nil {
		w.Write([]byte(jsoncatted))
	}

	return
}

func loginout(w http.ResponseWriter, resp_msg []byte) {
	w.Header().Set("Content-Type", "application/json")
	jsonresp := LoginResponse{LoginResponse: string(resp_msg)}
	jsoncatted, err := json.Marshal(jsonresp)
	if err == nil {
		w.Write([]byte(jsoncatted))
	}

	return
}
func Errorout(w http.ResponseWriter, err_type int, err_msg string) {
	w.Header().Set("Content-Type", "application/json")
	if err_type == 1 {
		jsonerror := Error{Error: err_msg, Success: 0}
		jsoncatted, err := json.Marshal(jsonerror)
		if err == nil {
			w.Write([]byte(jsoncatted))
		}
	}
	return
}

/* Create a new account */
func RegAccount(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	if r.PostForm["email"] == nil || r.PostForm["num"] == nil {
		Errorout(w, 1, "Must not have empty username or number")
		return
	}
	email := r.PostForm["email"][0]
	num := r.PostForm["num"][0]
	if len(email) == 0 || len(num) == 0 {
		Errorout(w, 1, "Must not have empty username or password")
		return
	}
	salt := gensalt()
	apikey := getnewapikey(salt)
	/* create user in db */
	rb := model.CreateNewUser(email, num, apikey, "free", r, w)
	if false == rb {
		Errorout(w, 1, "Account Creation Error")
	} else {
		registerout(w, []byte("Account created successfully"))
	}
}

func CodeVerify(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	if r.PostForm["code"] == nil || r.PostForm["email"] == nil {
		Errorout(w, 1, "Must specify email and code paramter")
		return
	}
	code := r.PostForm["code"][0]
	email := r.PostForm["email"][0]
	if len(code) == 0 || len(email) == 0{
		Errorout(w, 1, "Must not have an empty code or email parameter")
		return
	}
	login_status := model.AuthCode(email, code, r, w)
	if login_status {
		loginout(w, []byte("Successfully logged in"))
	} else {
		loginout(w, []byte("Could not log in"))
	}
	return
}

/* Returns false if s contains any non-alphanumeric characters */
/* or if the len  of the string is 0 */
func filterin(s string) bool {
	m := regexp.MustCompile(`^(\w|\d)*$`)
	all := m.FindString(s)
	if len(all) == 0 {
		return false
	}
	return true
}

func LogAccount(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	if r.PostForm["email"] == nil {
		Errorout(w, 1, "Must specify email and password paramter")
		return
	}
	email := r.PostForm["email"][0]
	if len(email) == 0 {
		Errorout(w, 1, "Must not have empty username or password")
		return
	}
	salt := gensalt()[:10]
	msg := fmt.Sprintf("Please enter the following code: %s", salt)
	num, err:= model.AuthUser(email, salt, r, w)

	if err == nil {
		smsmailer.Send_sms(msg, num)
		loginout(w, []byte("Code Sent"))
	} else {
		fmt.Println(err)
	}
}
