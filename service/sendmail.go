package service

import (
	"fmt"
	"gomailer/utils"
	"net/smtp"
	"os"
	"strings"
)

var (
    smtpServerAddr = utils.GetEnvVariable("AWS_SMTP_SERVER_ADDR")
    smtpServerPort = utils.GetEnvVariable("ASW_SMTP_SERVER_PORT")
)

func SetupSMTPAuth(username string, password string, serverAddr string)smtp.Auth{
    if username == "" || password == "" || serverAddr == "" {
        fmt.Printf("Invalid values for authentication!\n")
        os.Exit(0)
        fmt.Printf("%s \n%s \n%s \n",username,password,serverAddr)
    }
    fmt.Println("AWS SES Authentication requested.")
    auth := smtp.PlainAuth("", username, password, serverAddr)
    fmt.Println("AWS SES Authenticated successfuly!")
    return auth
}

// SendMail sends an email using the provided parameters.
func SendMail(auth smtp.Auth,to []string, subject, body,senderEmail string) error {

    mime := "MIME-version: 1.0\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n"
    // formatedSubject := fmt.Sprintf("Subject: %s\r\n\r\n",subject)

    headers := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n%s\r\n", 
        senderEmail, strings.Join(to, ","), subject, mime)
    

    msg := []byte(headers + body)
    err := smtp.SendMail(smtpServerAddr+":"+smtpServerPort, auth, senderEmail, to, msg)
   
    if err != nil {
        fmt.Printf("Error to sending email: %s", err)
        return err
    }
    
    fmt.Println("email sent success")
    return nil
}
