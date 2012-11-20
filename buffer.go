package mail

import (
	"bytes"
	"mime"
	"errors"
	"fmt"
	"strings"
)

const crlf = "\r\n"

func addrToString(list []Address) []string {
	var output []string
	for _, address := range list {
		output = append(output, address.String())
	}
	return output
}

// converts a message to a buffer for SMTP sending
func (m *Message) Buffer() (*bytes.Buffer, error) {
	var buf bytes.Buffer 

	if m.Sender == nil && len(m.From) > 0 {
		m.Sender = m.From[0]
	}

	fmt.Fprintf(&buf, "Content-Type: %s%s", m.ContentType, crlf)
	fmt.Fprintf(&buf, "Message-ID: %s%s", m.MessageId, crlf)
	fmt.Fprintf(&buf, "In-Reply-To: %s%s", strings.Join(m.InReply, " "), crlf)
	fmt.Fprintf(&buf, "References: %s%s", strings.Join(m.References, ","), crlf)
	fmt.Fprintf(&buf, "Date: %s%s", m.Date.Format(dateFormats[0]), crlf)
	fmt.Fprintf(&buf, "From: %s%s", strings.Join(addrToString(m.From), ","), crlf)
	fmt.Fprintf(&buf, "Sender: %s%s", m.Sender, crlf)
	fmt.Fprintf(&buf, "Reply-To: %s%s", strings.Join(addrToString(m.ReplyTo), ","), crlf)
	fmt.Fprintf(&buf, "To: %s%s", strings.Join(addrToString(m.To), ","), crlf)
	fmt.Fprintf(&buf, "Cc: %s%s", strings.Join(addrToString(m.Cc), ","), crlf)
	fmt.Fprintf(&buf, "Bcc: %s%s", strings.Join(addrToString(m.Bcc), ","), crlf)
	fmt.Fprintf(&buf, "Subject: %s%s", m.Subject, crlf)
	fmt.Fprintf(&buf, "Comments: %s%s", strings.Join(m.Comments, " "), crlf)
	fmt.Fprintf(&buf, "Keywords: %s%s", strings.Join(m.Keywords, ","), crlf)

	for _, h := range m.HeaderInfo.OptHeaders {
		fmt.Fprintf(&buf, "%s: %s%s", h.Key, h.Value, crlf)
	}

	if m.Parts == nil {
		fmt.Fprintf(&buf, "%s", m.Text)
	} else {
		_, ps, err := mime.ParseMediaType(m.ContentType)
		if err != nil { return nil, err }
		boundary, ok := ps["boundary"]
		if !ok {
			return nil, errors.New("multipart specified without boundary")
		}
		for _, part := range m.Parts {
			fmt.Fprintf(&buf, "%s%s", boundary, crlf)
			for key, value := range part.Headers {
				fmt.Fprintf(&buf, "%s: %s%s", key, value, crlf)
			}
			fmt.Fprintf(&buf, "%s%s%s%s", crlf, crlf, part.Data, crlf)
		}
		fmt.Fprintf(&buf, "%s%s", boundary, crlf)
	}

	return &buf, nil
}
