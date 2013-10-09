package mail

func newMessage(subject, from string, to []string) (*Message, error) {
	var err error
	message := new(Message)
	// Set charset and encoding for email subject.
	// See RFC 1342.
	message.Subject = "=?UTF-8?Q?" + subject + "?="
	message.From, err = parseAddressList([]byte(from))
	if err != nil {
		return nil, err
	}
	message.ReplyTo = append(message.ReplyTo, message.From...)
	for _, recipient := range to {
		address, err := ParseAddress([]byte(recipient))
		if err != nil {
			return nil, err
		}
		message.To = append(message.To, address)
	}
	return message, nil
}

// returns a message with minimum headers and textual content
func NewTextMessage(subject, content, from string, to ...string) (*Message, error) {
	message, err := newMessage(subject, from, to)
	if err != nil {
		return nil, err
	}
	message.Text = content
	message.ContentType = "text/plain"
	return message, nil
}

// returns a message with minimum headers and html content
func NewHTMLMessage(subject, content, from string, to ...string) (*Message, error) {
	message, err := newMessage(subject, from, to)
	if err != nil {
		return nil, err
	}
	message.Text = content
	message.ContentType = "text/html"
	return message, nil
}
