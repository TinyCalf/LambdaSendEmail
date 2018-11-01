package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/smtp"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
)

// SQS sqs form
type SQS struct {
	Records []struct {
		Body string `json:"body"`
	} `json:"Records"`
}

// Params the params i need
type Params struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

func sendMail(username, useraddress, password, host, to, subject, body, mailtype string) error {
	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", useraddress, password, hp[0])
	var contentType string
	if mailtype == "html" {
		contentType = "Content-Type: text/" + mailtype + "; charset=UTF-8"
	} else {
		contentType = "Content-Type: text/plain" + "; charset=UTF-8"
	}
	msg := []byte("To: " + to + "\r\nFrom: " + username + "<" + useraddress + ">\r\nSubject: " + subject + "\r\n" + contentType + "\r\n\r\n" + body)
	sendTo := strings.Split(to, ";")
	err := smtp.SendMail(host, auth, useraddress, sendTo, msg)
	return err
}

var (
	host        = "smtp.163.com:25"
	username    = "TinyCalf"
	useraddress = "15061519070@163.com"
	password    = "77e7c96a"
)

// HandleRequest handler
func HandleRequest(ctx context.Context, name SQS) (string, error) {
	fmt.Printf("Received: %+v", name)
	records := name.Records
	for i := range records {
		b := records[i].Body
		params := &Params{}
		err := json.Unmarshal([]byte(b), &params)
		if err != nil {
			return "", err
		}

		to := params.Email
		code := params.Code
		subject := "Text email from AWS SQS lambda"
		body := fmt.Sprintf(`
			<html>
			<body>
			<h3>
			Your code is %v
			</h3>
			</body>
			</html>
			`, code)
		err = sendMail(username, useraddress, password, host, to, subject, body, "html")
		if err != nil {
			return "", err
		}
		return "发送成功", nil
	}
	return "", nil
}

func main() {
	lambda.Start(HandleRequest)
	// fmt.Println("sending...")
	// err := sendMail(username, useraddress, password, host, useraddress, "Title", "aha", "html")
	// if err != nil {
	// 	fmt.Println("发送失败")
	// 	return
	// }
	// fmt.Println("发送成功")
}
