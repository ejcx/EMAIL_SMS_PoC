package model

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"errors"
	"time"
	"net/http"
	//"database/sql"
	"github.com/gorilla/sessions"
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native" // Native engine
	// _ "github.com/ziutek/mymysql/thrsafe" // Thread safe engine
)

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_AES_KEY")))
func getdb() mysql.Conn {
  db := mysql.New("tcp", "", "127.0.0.1:3306", "proddbuser", os.Getenv("MYSQL_PW"), "minibase");
  return db
}

func shatwo(usermaterial string) string {
	var h = sha256.New()
	h.Write([]byte(usermaterial))
	pwhash := hex.EncodeToString(h.Sum(nil))
	return pwhash
}
func CreateNewUser(email string, num string, apikey string, level string, r *http.Request, w http.ResponseWriter) bool {
	db := getdb()
	err := db.Connect()
	if err != nil {
		panic(err)
		return false
	}

        exists, _ := Check_exists(email)
	if exists {
		return false
	}
	ins, err := db.Prepare("INSERT INTO user VALUES (?,?,?,?,?,?)")

	if err != nil {
		panic(err)
		return false
	}
	ins.Bind([]byte(nil), email, num, level, apikey, []byte(nil))
	_, err = ins.Run()

	/*setup the user's session, this is icing */
	/* so don't return false if it fails */
	if err == nil {
		session, err := store.Get(r, "SESSIONID")
		if err != nil {

		}
		session.Values["lvl"] = level
		session.Values["apikey"] = apikey
		session.Values["email"] = email
		session.Values["num"] = num
		session.Save(r, w)
		return true
	} else {
		return false
	}
	if err != nil {
		return false
	}
	return true
}

func AuthCode(email string, code string, r *http.Request, w http.ResponseWriter) bool {
	db := getdb()
	err := db.Connect()
	if err != nil {
		panic(err)
		return false
	}

	/*check if account exists already*/
	sel, err := db.Prepare("SELECT a.timestamp, u.num, u.acctlevel, u.apikey from authcodes as a INNER JOIN user as u on u.email=a.email where a.email=? and a.code=?")
	if err != nil {
		fmt.Println(err)
	}
	sel.Bind(email, code)
	res, err := sel.Run()
	if err != nil {
		fmt.Println(err)
	}
	rows, _ := res.GetRows()

	if len(rows) != 1 {
		return false
	}
	ts := int32(rows[0][0].(int32))
	timestamp := int32(time.Now().Unix())

	if timestamp > (ts+500){
		fmt.Println("Expired")
		return false
	} else {
		NUM := string(rows[0][1].([]byte))
		LVL := string(rows[0][2].([]byte))
		APIKEY := string(rows[0][3].([]byte))
		session, err := store.Get(r, "SESSIONID")
		if err == nil {
			session.Values["lvl"] = LVL
			session.Values["apikey"] = APIKEY
			session.Values["email"] = email
			session.Values["num"] = NUM
			session.Save(r, w)
		}
		return true
	}
}
func AuthUser(email string,salt string, r *http.Request, w http.ResponseWriter) (string, error){
	db := getdb()
	err := db.Connect()
	if err != nil {
		return "", err
	}
	/*check if account exists already*/
	sel, err := db.Prepare("SELECT * from user where email=?")
	if err != nil {
		fmt.Println(err)
	}
	sel.Bind(email)
	res, _ := sel.Run()
	row, _ := res.GetFirstRow()
	//+--------------+-------------+
	//| Field        | Type        |
	//+--------------+-------------+
	//| uid          | int(11)     |
	//| email        | varchar(64) |
	//| num          | varchar(64) |
	//| acctlevel    | varchar(64) |
	//| apikey       | varchar(64) |
	//| stripecharge | blob        |
	//+--------------+-------------+
	if len(row) == 0 {
		return "", errors.New("User Does Not Exist")
	}
	NUM:= string(row[2].([]byte))

	timestamp := int32(time.Now().Unix())
	code := salt
	ins, err := db.Prepare("INSERT into authcodes VALUES (?,?,?,?) ON DUPLICATE KEY UPDATE code=?, timestamp=?, attempts=?")
	if err != nil {
		fmt.Println(err)
	}
	ins.Bind(email , code, timestamp, 0, code, timestamp, 0)
	_, err = ins.Run()
	if err != nil {
		fmt.Println(err)
		return NUM, err
	}
	return NUM, nil
}

/* Returns true if the user exists, and their phone number */
func Check_exists(or string) (bool, string){
	db := getdb()
	err := db.Connect()
        if err != nil {
                panic(err)
                return false, ""
        }
	up, err := db.Prepare("SELECT num from user where email=?")
	if err != nil {
		panic(err)
		return false, ""
	}
	up.Bind(or)
	res, err := up.Run()
	rows, _ := res.GetRows()
	if len(rows) == 0 {
		return false, ""
	}
	if len(rows) > 1 {
		fmt.Printf("REGISTERED TWO USERS %s", or)
		return false, ""
	}
	return true, string(rows[0][0].([]byte))
}
