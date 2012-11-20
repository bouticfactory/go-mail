package mail

import (
	"net/smtp"
	"errors"
	"fmt"
	"strings"
	"io"
)

// implements io.Writer interface to use template.ExecuteTemplate()
//
// Example:
// 	 body := new(mail.SMTPWriter)
//
//	 templates.ExecuteTemplate(body, "account_validation", boutique)
//
//   message, err := mail.NewHTMLMessage(subject, body.String(), from, to)
//   mailer.Send(message)
type SMTPWriter struct {
	io.Writer
	content []string
}

// simply store to []string
func (self *SMTPWriter)Write(p []byte) (n int, err error) {
	self.content = append(self.content, string(p))
	return len(p), nil
}

func (self *SMTPWriter)String() string {
	return strings.Join(self.content, "")
}

// structure to connect to a SMTP server
type SMTPServer struct {
	Host     string
	Port     uint16
	Username string
	Password string
	From     Address
	Auth     *smtp.Auth
}

func (self *SMTPServer) Name() string {
	return "email"
}

func (self *SMTPServer) Close() error {
	return nil
}

// returns the simplest possible server (localhost:25)
func InitLocalHost(fromAddress string) (*SMTPServer, error) {
	from, err := ParseAddress([]byte(fromAddress))
	if err != nil {
			return nil, err
	}

	config := SMTPServer{
		Host: "localhost",
		Port: 25,
		From: from,
	}
	
	return &config, nil
}

// returns a server to a gmail account
// uses a different fromAddress for sending emails
// than the loginAddress used together with the password for authentication
// The fromAddress has to be verified as a valid sender address in Gmail.
func InitGmail(fromAddress, loginAddress, password string) (*SMTPServer, error) {
	if _, err := ParseAddress([]byte(loginAddress)); err != nil {
		return nil, err
	}
	if len(fromAddress) == 0 {
		fromAddress = loginAddress
	}
	from, err := ParseAddress([]byte(fromAddress))
	if err != nil {
			return nil, err
	}

	config := SMTPServer{
		Host:     "smtp.gmail.com",
		Port:     587,
		Username: loginAddress,
		Password: password,
		From:     from,
	}
	return &config, nil
}

// checks minimum requirements before trying to send a message
func (m *Message) Validate() error {
	if len(m.To) == 0 {
		return errors.New("Missing email recipient (email.Message.To)")
	}
	return nil
}

func sendWithAuth(server *SMTPServer, m *Message) error {
	addr := fmt.Sprintf("%s:%d", server.Host, server.Port)

	if err := m.Validate(); err != nil {
		return err
	}
	auth := smtp.PlainAuth("", server.Username, server.Password, server.Host)
	content, _ := m.Buffer()
	return smtp.SendMail(addr, auth, strings.Join(addrToString(m.From), ";"), addrToString(m.To), []byte(content.String()))
}


func sendWithoutAuth(server *SMTPServer, m *Message) error {
	addr := fmt.Sprintf("%s:%d", server.Host, server.Port)

	c, err := smtp.Dial(addr)
	if err != nil {
		return err
	}

	c.Mail(strings.Join(addrToString(m.From), ";"))
	c.Rcpt(strings.Join(addrToString(m.To), ";"))

	wc, err := c.Data()
	if err != nil {
		return err
	}
	defer wc.Close()

	buf, _ := m.Buffer()
	_, err = buf.WriteTo(wc)
	return err
}

// sends a message via the SMTP server
func (server *SMTPServer) Send(m *Message) error {
	if m == nil {
		return errors.New("Message is nil")
	}
	if err := m.Validate(); err != nil {
		return err
	}

	if server.Auth != nil {
		return sendWithAuth(server, m)
	}
	return sendWithoutAuth(server, m)
}
