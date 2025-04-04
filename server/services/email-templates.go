package services

import "fmt"

type IEmailTemplates interface {
	GetNewUserEmailTemplate(baseURL string, verificationToken string) string
	GetLoginEmailTemplate(baseURL string, verificationToken string) string
}

type EmailTemplates struct {
}

func NewEmailTemplates() IEmailTemplates {
	return &EmailTemplates{}
}

func (e *EmailTemplates) GetNewUserEmailTemplate(baseURL string, verificationToken string) string {
	return fmt.Sprintf(`
		<p>
			Hello and Welcome to OpenZooSim! Please verify your account by 
			<a href="%v/auth/verify?token=%v">Clicking Here!</a>
		</p>
	`, baseURL, verificationToken)
}

func (e *EmailTemplates) GetLoginEmailTemplate(baseURL string, verificationToken string) string {
	return fmt.Sprintf(`
		<p>
			Hello! You can login to your account by 
			<a href="%v/auth/login-with-email?token=%v">Clicking Here!</a>
		</p>

		<p>
			If you did not request this email, please ignore it.
		</p>
	`, baseURL, verificationToken)
}
