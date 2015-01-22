package smsmailer 

import (
	"net/url"
	"net/http"
	"fmt"
	"os"
	"strings"
	"io/ioutil"
)
func Send_sms(content string, number string) string{
	v := url.Values{}
	v.Set("To", number)
	v.Set("From", "+15406841667")
	v.Set("Body", content)
	//Contains User key (public)
	req, err := http.NewRequest("POST", "https://api.twilio.com/2010-04-01/Accounts/AC3931fc9e83f3d38a4d29ac8efcd19915/Messages.json", strings.NewReader(v.Encode()))
	req.SetBasicAuth(os.Getenv("TWILIO_USER"), os.Getenv("TWILIO_PW"))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
	return string(body)
}
