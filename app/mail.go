package app

import (
	"bytes"
	"context"
	"html/template"
	"time"

	"github.com/mailgun/mailgun-go/v4"
)

// MailService handles sending emails
type MailService struct {
	FromAddress string
	Hostname    string
	MailClient  mailgun.Mailgun
	Templates   map[string]*template.Template
}

// SendWelcomeEmail sends a welcome email to a user. A verification link is included
func (s *MailService) SendWelcomeEmail(email, name, verificationID string) error {
	templateData := struct {
		Name string
		URL  string
	}{
		Name: name,
		URL:  s.Hostname + "/v1/auth/verify-email/" + verificationID,
	}

	tpl := s.Templates["email/welcome"]
	buf := new(bytes.Buffer)
	if err := tpl.Execute(buf, templateData); err != nil {
		return err
	}

	return s.send(email, "Welcome to Playtest Co-op!", buf.String())
}

// SendVerifyEmail sends an email to a user to verify that their email address is legit
func (s *MailService) SendVerifyEmail(email, name, verificationID string) error {
	templateData := struct {
		Name string
		URL  string
	}{
		Name: name,
		URL:  s.Hostname + "/v1/auth/verify-email/" + verificationID,
	}

	tpl := s.Templates["email/verify-email"]
	buf := new(bytes.Buffer)
	if err := tpl.Execute(buf, templateData); err != nil {
		return err
	}

	return s.send(email, "Verify your email", buf.String())
}

// SendPasswordResetEmail sends an email with an included one-time-password for users to use to set their password again
func (s *MailService) SendPasswordResetEmail(email, name, otp string) error {
	templateData := struct {
		Name string
		URL  string
	}{
		Name: name,
		URL:  s.Hostname + "/v1/auth/reset-password/" + otp,
	}

	tpl := s.Templates["email/reset-password"]
	buf := new(bytes.Buffer)
	if err := tpl.Execute(buf, templateData); err != nil {
		return err
	}

	return s.send(email, "Password reset requested", buf.String())
}

func (s *MailService) send(toAddress, subject, body string) error {
	message := s.MailClient.NewMessage(
		s.FromAddress,
		subject,
		"",
		toAddress,
	)
	message.SetHtml(body)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	_, _, err := s.MailClient.Send(ctx, message)
	if err != nil {
		return err
	}

	return nil
}
