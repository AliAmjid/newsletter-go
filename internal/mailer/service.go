package mailer

import (
    "bytes"
    "fmt"
    "html/template"
    "io"
    "net/http"
    "net/url"
    "strings"
)

type Service struct {
    domain    string
    apiKey    string
    fromEmail string
    templates *template.Template
}

func NewService(domain, apiKey, from string) (*Service, error) {
    tmpl, err := template.ParseGlob("templates/*.html")
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
    endpoint := fmt.Sprintf("https://api.mailgun.net/v3/%s/messages", s.domain)
    values := url.Values{}
    values.Set("from", s.fromEmail)
    values.Set("to", to)
    values.Set("subject", subject)
    values.Set("html", body)

    req, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(values.Encode()))
    if err != nil {
        return err
    }
    req.SetBasicAuth("api", s.apiKey)
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    if resp.StatusCode >= 300 {
        b, _ := io.ReadAll(resp.Body)
        return fmt.Errorf("mailgun error: %s", string(b))
    }
    return nil
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
