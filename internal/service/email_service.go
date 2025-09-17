package service

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"html/template"
	"net"
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

func (e *EmailService) buildVerificationMessage(toEmail, username, verificationCode string, registrationDate time.Time) ([]byte, error) {
	tmpl, fileErr := template.ParseFiles("/app/templates/verification_email.html")
	if fileErr != nil {
		return nil, fileErr
	}
	EmailData := VerificationEmailData{
		Username:         username,
		Email:            toEmail,
		VerificationCode: verificationCode,
		RegistrationDate: registrationDate,
	}
	var buffer bytes.Buffer
	if err := tmpl.Execute(&buffer, EmailData); err != nil {
		return nil, err
	}

	fallback := fmt.Sprintf("Hi %s!\nYour verification code is: %s\nRegistered on: %s\n",
		username, verificationCode, registrationDate.Format(time.RFC1123))

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
	return msg, nil
}

func (e *EmailService) SendVerificationEmailCtx(ctx context.Context, toEmail, username, verificationCode string, registrationDate time.Time) error {
	msg, err := e.buildVerificationMessage(toEmail, username, verificationCode, registrationDate)
	if err != nil {
		return err
	}
	host, port := e.Host, e.Port
	addr := net.JoinHostPort(host, port)

	// Dial with context -> hard timeout
	dialer := &net.Dialer{}
	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}
	defer c.Quit()

	// STARTTLS upgrade (587)
	if ok, _ := c.Extension("STARTTLS"); ok {
		tlsCfg := &tls.Config{ServerName: host}
		if err := c.StartTLS(tlsCfg); err != nil {
			return err
		}
	}

	// AUTH
	if ok, _ := c.Extension("AUTH"); ok {
		if err := c.Auth(smtp.PlainAuth("", e.Email, e.Password, host)); err != nil {
			return err
		}
	}

	if err := c.Mail(e.Email); err != nil {
		return err
	}
	if err := c.Rcpt(toEmail); err != nil {
		return err
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	if _, err := w.Write(msg); err != nil {
		return err
	}
	return w.Close()
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
