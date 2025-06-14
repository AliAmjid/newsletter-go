package mailer

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log"
	"path/filepath"
	"runtime"
	"time"

	"github.com/mailgun/mailgun-go/v4"

	"newsletter-go/domain"
)

type Service struct {
	domain    string
	apiKey    string
	fromEmail string
	templates *template.Template
}

func NewService(domain, apiKey, from string) (*Service, error) {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	tmpl, err := template.ParseGlob(filepath.Join(dir, "templates", "*.html"))
	if err != nil {
		return nil, err
	}
	return &Service{domain: domain, apiKey: apiKey, fromEmail: from, templates: tmpl}, nil
}

func (s *Service) render(name string, data interface{}) (string, error) {
	var buf bytes.Buffer
	if err := s.templates.ExecuteTemplate(&buf, name, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (s *Service) send(to, subject, body string) error {
	mg := mailgun.NewMailgun(s.domain, s.apiKey)

	from := s.fromEmail
	// Pokud fromEmail neobsahuje "<" a ">", přidej výchozí jméno
	if !containsAngleBrackets(from) {
		from = fmt.Sprintf("Mailgun Sandbox <%s>", s.fromEmail)
	}
	m := mg.NewMessage(
		from,
		subject,
		"",
		to,
	)
	m.SetHtml(body)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, _, err := mg.Send(ctx, m)
	if err != nil {
		log.Printf("mailgun send error: %v", err)
		return fmt.Errorf("unable to send email")
	}
	return nil
}

func containsAngleBrackets(s string) bool {
	return len(s) > 0 && (contains(s, "<") && contains(s, ">"))
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (len(substr) == 0 || (len(s) > 0 && (s[0:len(substr)] == substr || contains(s[1:], substr))))
}

type ForgotPasswordData struct {
	Token string
}

func (s *Service) SendForgotPasswordEmail(to, token string) error {
	body, err := s.render("password_reset.html", ForgotPasswordData{Token: token})
	if err != nil {
		return err
	}
	return s.send(to, "Password Reset", body)
}

type SubscriptionData struct {
	Token string
}

func (s *Service) SendSubscriptionConfirmEmail(to, token string) error {
	body, err := s.render("subscription_confirm.html", SubscriptionData{Token: token})
	if err != nil {
		return err
	}
	return s.send(to, "Confirm Subscription", body)
}

type PostEmailData struct {
	Title      string
	Content    template.HTML
	PixelURL   string
	UnsubToken string
}

func (s *Service) SendPostEmail(to, token string, p *domain.Post, deliveryID string) error {
	data := PostEmailData{
		Title:      p.Title,
		Content:    template.HTML(p.Content),
		PixelURL:   fmt.Sprintf("http://localhost:3000/post-deliveries/%s/pixel", deliveryID),
		UnsubToken: token,
	}
	body, err := s.render("post.html", data)
	if err != nil {
		return err
	}
	return s.send(to, p.Title, body)
}
