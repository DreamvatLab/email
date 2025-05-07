package svc

import (
	"context"
	"os"
	"testing"

	"github.com/DreamvatLab/email"
	"github.com/stretchr/testify/require"
)

func TestEmailService_Send(t *testing.T) {
	s := new(EmailService)
	mail := &email.EmailDTO{
		AccountID: "dev@DreamvatLab.com",
		From:      "dev@DreamvatLab.com",
		To:        []string{"lukiya@outlook.com"},
		CC:        []string{"yunuo215@outlook.com", "jcc1275@outlook.com"},
		BCC:       []string{"eyong551@outlook.com", "eyun551@outlook.com"},
		Subject:   "LF Ledger Summary",
		Body:      `<html><body><h1>Headline</h1><span style="color:red">Test</span> message</body></html>`,
	}

	attachmentData, _ := os.ReadFile("M:/FEES-2023-09-27.xlsx")

	mail.Attachments = append(mail.Attachments, &email.AttachmentDTO{
		Name: "FEES-2023-09-27.xlsx",
		Data: attachmentData,
	})
	rs, err := s.Send(context.Background(), mail)
	require.NoError(t, err)
	t.Log(rs.Message)
	require.True(t, rs.Success)
}
