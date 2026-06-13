package sendgrid

import (
	"backend/internal/domain/auth"
	"context"
	"fmt"
	"net/http"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridAdapter struct {
	client    *sendgrid.Client
	fromEmail string
	fromName  string
	baseURL   string
}

func NewService(apiKey, fromEmail, fromName, baseURL string) auth.EmailPort {
	client := sendgrid.NewSendClient(apiKey)
	return &SendGridAdapter{
		client:    client,
		fromEmail: fromEmail,
		fromName:  fromName,
		baseURL:   baseURL,
	}
}

func (s *SendGridAdapter) SendEmailVerification(ctx context.Context, toEmail, token string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	from := mail.NewEmail(s.fromName, s.fromEmail)
	to := mail.NewEmail("RoadBuddy User", toEmail)

	subject := "Verify Your RoadBuddy Account"

	verificationURL := fmt.Sprintf("%s/api/v1/auth/verify?token=%s", s.baseURL, token)

	htmlContent := fmt.Sprintf(`
		<h3>Welcome to RoadBuddy!</h3>
		<p>Please verify your email address by clicking the link below:</p>
		<p><a href="%s" style="background-color: #4CAF50; color: white; padding: 10px 20px; text-decoration: none; display: inline-block; border-radius: 4px;">Verify Email</a></p>
	`, verificationURL)

	textContent := fmt.Sprintf("Welcome to RoadBuddy! Verify your email by visiting: %s", verificationURL)

	message := mail.NewSingleEmail(from, subject, to, textContent, htmlContent)

	response, err := s.client.Send(message)
	if err != nil {
		return fmt.Errorf("Failed to send verification email: %w", err)
	}

	if response.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("SendGrid API error: status %d, body %s", response.StatusCode, response.Body)
	}

	return nil
}
