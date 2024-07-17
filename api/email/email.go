package email

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
)

func SendCode(email string, code string) {
	// sender data
	from := "articanconnection@gmail.com" // shaxsiy email kiriting
	password := "colo twdh fabv kcvj"  

	to := []string{
		email,
	}

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	auth := smtp.PlainAuth("", from, password, smtpHost)

	t, err := template.ParseFiles("template.html")
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
		return
	}

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: Your verification code \n%s\n\n", mimeHeaders)))

	err = t.Execute(&body, struct {
		Passwd string
	}{
		Passwd: code,
	})
	if err != nil {
		log.Fatalf("Error executing template: %v", err)
		return
	}

	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
	if err != nil {
		log.Fatalf("Error sending email: %v", err)
		return
	}
}
