package service

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"time"
)

type EmailService struct {
	Email, Password, Host, Port string
}

type VerificationEmailData struct {
	Username, Email, VerificationCode string
	RegistrationDate                  time.Time
}

func NewEmailService(email, password, host, port string) *EmailService {
	return &EmailService{
		Email:    email,
		Password: password,
		Host:     host,
		Port:     port,
	}
}

func (e *EmailService) SendVerificationEmail(toEmail string, username string, verificationCode string, registrationDate time.Time) error {
	tmpl, fileErr := template.ParseFiles("/app/templates/verification_email.html")
	if fileErr != nil {
		return fileErr
	}
	EmailData := VerificationEmailData{
		Username:         username,
		Email:            toEmail,
		VerificationCode: verificationCode,
		RegistrationDate: registrationDate,
	}
	buffer := bytes.Buffer{}
	tmplErr := tmpl.Execute(&buffer, EmailData)
	if tmplErr != nil {
		return tmplErr
	}

	fallback := fmt.Sprintf("Hi %s!\nYour verification code is: %s\nRegistered on: %s\n", username, verificationCode, registrationDate.Format(time.RFC1123))

	boundary := "BOUNDARY-7e1f2c3"

	headers := ""
	headers += fmt.Sprintf("From: %s\r\n", e.Email)
	headers += fmt.Sprintf("To: %s\r\n", toEmail)
	headers += "Subject: Your Verification Code\r\n"
	headers += "MIME-Version: 1.0\r\n"
	headers += fmt.Sprintf("Content-Type: multipart/alternative; boundary=%q\r\n", boundary)

	body := ""
	body += fmt.Sprintf("--%s\r\n", boundary)
	body += "Content-Type: text/plain; charset=\"UTF-8\"\r\n"
	body += "Content-Transfer-Encoding: 7bit\r\n\r\n"
	body += fallback + "\r\n"

	body += fmt.Sprintf("--%s\r\n", boundary)
	body += "Content-Type: text/html; charset=\"UTF-8\"\r\n"
	body += "Content-Transfer-Encoding: 7bit\r\n\r\n"
	body += buffer.String() + "\r\n"

	body += fmt.Sprintf("--%s--\r\n", boundary)

	msg := []byte(headers + "\r\n" + body)

	auth := smtp.PlainAuth("", e.Email, e.Password, e.Host)
	addr := fmt.Sprintf("%s:%s", e.Host, e.Port)

	return smtp.SendMail(addr, auth, e.Email, []string{toEmail}, msg)
}
