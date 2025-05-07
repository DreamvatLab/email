package svc

import (
	"bytes"
	"context"
	"crypto/tls"
	"strings"

	"github.com/DeamvatLab/email"
	"github.com/DeamvatLab/email/main/dal"
	"github.com/DreamvatLab/go/xconv"
	"github.com/DreamvatLab/go/xerr"
	"gopkg.in/mail.v2"
)

var (
	_dal = dal.NewDAL()
)

type EmailService struct{}

// Send implements email.EmailServiceServer.
func (*EmailService) Send(ctx context.Context, dto *email.EmailDTO) (*email.SendResponse, error) {
	if dto.AccountID == "" {
		err := xerr.New("AccountID cannot be empty")
		xerr.LogError(err)
		return nil, err
	}

	account, err := _dal.GetEmailAccount(dto.AccountID)
	if xerr.LogError(err) {
		return nil, err
	} else if account == nil {
		err = xerr.Errorf("Email account '%s' does not exist", dto.AccountID)
		xerr.LogError(err)
		return nil, err
	}

	// err = sendByNativePackage(account, dto)
	err = sendByMailPackage(account, dto)
	r := new(email.SendResponse)
	if err != nil {
		r.Success = false
		r.Message = err.Error()
		xerr.LogError(err)
		return r, err
	}

	r.Success = true
	return r, nil
}

// func sendByNativePackage(account *email.EmailAccount, dto *email.EmailDTO) error {

// 	// 使用smtp.PlainAuth创建SMTP认证
// 	auth := smtp.PlainAuth("", account.Username, account.Password, account.SMTPAddress)

// 	// 创建邮件头
// 	header := make(map[string]string)
// 	header["From"] = dto.From
// 	if len(dto.To) > 0 {
// 		header["To"] = strings.Join(dto.To, ";")
// 	}
// 	if len(dto.CC) > 0 {
// 		header["Cc"] = strings.Join(dto.CC, ";")
// 	}
// 	header["Subject"] = dto.Subject
// 	header["MIME-Version"] = "1.0"

// 	if len(dto.Attachments) > 0 { // 有附件
// 		header["Content-Type"] = `multipart/mixed; boundary="mix"`
// 	} else { // 没附件
// 		header["Content-Type"] = `text/html; charset="UTF-8"`
// 	}

// 	// 创建邮件
// 	var msg bytes.Buffer
// 	for k, v := range header {
// 		msg.WriteString(fmt.Sprintf("%s: %s\n", k, v))
// 	}
// 	msg.WriteString("\n")

// 	if len(dto.Attachments) > 0 { // 有附件
// 		// 添加邮件正文
// 		msg.WriteString("--mix\n")
// 		msg.WriteString("Content-Type: text/html; charset=UTF-8\n")
// 		msg.WriteString("Content-Transfer-Encoding: quoted-printable\n\n")
// 		msg.WriteString(dto.Body)
// 		msg.WriteString("\n")

// 		// 添加每个附件
// 		for _, attachment := range dto.Attachments {
// 			// Base64 编码附件
// 			encodedAttachment := base64.StdEncoding.EncodeToString(attachment.Data)

// 			// 添加附件
// 			msg.WriteString("--mix\n")
// 			msg.WriteString("Content-Type: application/octet-stream\n")
// 			msg.WriteString("Content-Transfer-Encoding: base64\n")
// 			msg.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"\n\n", attachment.Name))
// 			msg.WriteString(encodedAttachment)
// 			msg.WriteString("\n")
// 		}
// 		msg.WriteString("--mix--")
// 	} else { // 没附件
// 		// 添加邮件正文
// 		msg.WriteString(dto.Body)
// 	}

// 	// 合并所有收件人，包括 To, Cc 和 Bcc
// 	var allRecipients []string
// 	allRecipients = append(allRecipients, dto.To...)
// 	allRecipients = append(allRecipients, dto.CC...)
// 	allRecipients = append(allRecipients, dto.BCC...)

// 	// 发送邮件
// 	err := smtp.SendMail(account.SMTPAddress, auth, dto.From, allRecipients, msg.Bytes())
// 	if err != nil {
// 		return xerr.WithStack(err)
// 	}

// 	return nil
// }

func sendByMailPackage(account *email.EmailAccount, dto *email.EmailDTO) error {
	// 创建邮件
	m := mail.NewMessage()
	// 设置邮件的发件人
	m.SetHeader("From", dto.From)
	// 设置邮件的收件人
	m.SetHeader("To", dto.To...)

	// 设置邮件的抄送人
	m.SetHeader("Cc", dto.CC...)
	// for i, x := range dto.CC {
	// 	m.SetAddressHeader("Cc", x, "CC"+xconv.ToString(i+1))
	// }

	// 设置邮件的密送人
	m.SetHeader("Bcc", dto.BCC...)
	// for i, x := range dto.BCC {
	// 	m.SetAddressHeader("Bcc", x, "BCC"+xconv.ToString(i+1))
	// }

	// 设置邮件的主题
	m.SetHeader("Subject", dto.Subject)

	// 设置邮件的正文
	m.SetBody(`text/html; charset="UTF-8"`, dto.Body)

	for _, x := range dto.Attachments {
		m.AttachReader(x.Name, bytes.NewReader(x.Data), mail.SetHeader(map[string][]string{
			"Content-Type": {"application/octet-stream"},
		}))
	}

	array := strings.Split(account.SMTPAddress, ":")
	if len(array) < 2 {
		return xerr.Errorf("Invalid smtp address: %s", account.SMTPAddress)
	}

	host := array[0]
	port := xconv.ToInt(array[1])

	// 设置SMTP服务器
	d := mail.NewDialer(host, port, account.Username, account.Password)

	// 启用 TLS
	d.StartTLSPolicy = mail.MandatoryStartTLS
	d.TLSConfig = &tls.Config{
		InsecureSkipVerify: true, // 跳过证书验证
	}

	// 发送
	err := d.DialAndSend(m)
	if err != nil {
		return xerr.WithStack(err)
	}

	return nil
}
