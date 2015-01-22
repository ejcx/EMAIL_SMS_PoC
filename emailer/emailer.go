package main
import (
	"flag"
	"os"
	"fmt"
	"model"
	"smsmailer"
	"strings"
	"io/ioutil"
)

func main(){
	origemail := flag.String("c", "", "The origin recipient")
        flag.Parse()

	if len(*origemail) == 0 {
		fmt.Println("-c is required. It is username")
		return
        } 
	stdinbytes, _:= ioutil.ReadAll(os.Stdin)
	var emailbody []byte
	var lastbyte byte
	writing := false
	
	// This is super hacky since I have no idea 
	// SMTP message data correctly	
	for i:=range stdinbytes {
		if lastbyte == 10  && stdinbytes[i] == 10 {
			if strings.Contains(string(emailbody),"Content-Type") {
				emailbody = []byte(nil)
			} else if len(emailbody)>1 {
				writing = false	
			} else {
				writing = true
			}
		} else {
			lastbyte = byte(stdinbytes[i])
		}

		if writing==true {
			emailbody = append(emailbody, stdinbytes[i])
		}
	}
	origparsed := strings.Replace(*origemail, "@twinsen.barnbase.com","", 1)
        tosend, number := model.Check_exists(origparsed) 
	if tosend {
		smsmailer.Send_sms(string(emailbody), number)
	}
}
