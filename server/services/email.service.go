package services

import (
	"fmt"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/snowlynxsoftware/oto-api/server/util"
)

type IEmailService interface {
	SendEmail(options *EmailSendOptions) bool
	GetTemplates() IEmailTemplates
}

type EmailService struct {
	client    *sendgrid.Client
	templates IEmailTemplates
}

func NewEmailService(apiKey string, templates IEmailTemplates) IEmailService {
	emailClient := sendgrid.NewSendClient(apiKey)
	return &EmailService{
		client:    emailClient,
		templates: templates,
	}
}

type EmailSendOptions struct {
	FromEmail   string
	ToEmail     string
	Subject     string
	HTMLContent string
}

func (s *EmailService) SendEmail(options *EmailSendOptions) bool {
	from := mail.NewEmail("", options.FromEmail)
	subject := options.Subject
	to := mail.NewEmail("", options.ToEmail)
	plainTextContent := "Open Trivia Online does not support plain text emails. Please enable HTML."
	htmlContent := options.HTMLContent
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	response, err := s.client.Send(message)
	if err != nil {
		util.LogErrorWithStackTrace(err)
		return false
	} else {
		if response.StatusCode != 202 {
			util.LogWarning(string(rune(response.StatusCode)))
			util.LogWarning(response.Body)
			util.LogWarning(fmt.Sprintf("%v", response.Headers))
			return false
		}
		return true
	}
}

func (s *EmailService) GetTemplates() IEmailTemplates {
	return s.templates
}
