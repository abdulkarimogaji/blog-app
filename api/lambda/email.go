package lambda

import (
	"fmt"
	"net/http"
	"net/smtp"

	"github.com/abdulkarimogaji/blognado/config"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/jordan-wright/email"
)

const (
	smtpAuthAddress   = "smtp.gmail.com"
	smtpServerAddress = "smtp.gmail.com:587"
)

type EmailSender interface {
	SendEmail(
		subject string,
		content string,
		to []string,
		cc []string,
		bcc []string,
		attachFiles []string,
	) error
}

type GmailSender struct {
	name              string
	fromEmailAddress  string
	fromEmailPassword string
}

func NewGmailSender(name string, fromEmailAddress string, fromEmailPassword string) EmailSender {
	return &GmailSender{
		name:              name,
		fromEmailAddress:  fromEmailAddress,
		fromEmailPassword: fromEmailPassword,
	}
}

func (sender *GmailSender) SendEmail(
	subject string,
	content string,
	to []string,
	cc []string,
	bcc []string,
	attachFiles []string,
) error {
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", sender.name, sender.fromEmailAddress)
	e.Subject = subject
	e.HTML = []byte(content)
	e.To = to
	e.Cc = cc
	e.Bcc = bcc

	for _, f := range attachFiles {
		_, err := e.AttachFile(f)
		if err != nil {
			return fmt.Errorf("failed to attach file %s: %w", f, err)
		}
	}

	smtpAuth := smtp.PlainAuth("", sender.fromEmailAddress, sender.fromEmailPassword, smtpAuthAddress)
	return e.Send(smtpServerAddress, smtpAuth)
}

type sendEmailRequest struct {
	Subject string   `json:"subject" binding:"required"`
	Content string   `json:"content" binding:"required"`
	To      []string `json:"to" binding:"min=1"`
	Cc      []string `json:"cc,omitempty"`
	Bcc     []string `json:"bcc,omitempty"`
}

func sendMailAPI(c *gin.Context) {
	var body sendEmailRequest
	err := c.ShouldBindBodyWith(&body, binding.JSON)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "bad request",
			"message": err.Error(),
			"error":   "true",
			"data":    nil,
		})
		return
	}

	gmailSender := NewGmailSender(config.AppConfig.GMAIL_NAME, config.AppConfig.GMAIL_ADDRESS, config.AppConfig.GMAIL_PASSWORD)
	err = gmailSender.SendEmail(body.Subject, body.Content, body.To, body.Cc, body.Bcc, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failure",
			"message": err.Error(),
			"error":   true,
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Email sent successfully",
		"error":   false,
		"data":    nil,
	})
}
