// Package mail implements a parser for electronic mail messages as specified
// in RFC5322.
//
// We allow both CRLF and LF to be used in the input, possibly mixed.
package mail

import (
	"crypto/sha1"
	"encoding/base64"
	"time"
)

var benc = base64.URLEncoding

func mkId(s []byte) string {
	h := sha1.New()
	h.Write(s)
	hash := h.Sum(nil)
	ed := benc.EncodeToString(hash)
	return ed[0:20]
}

// structure containing the header vars as defined in RFC5322
type HeaderInfo struct {
	FullHeaders []Header // all headers
	OptHeaders  []Header // unprocessed headers

	MessageId   string
	Id          string
	Date        time.Time
	From        []Address
	Sender      Address
	ReplyTo     []Address
	To          []Address
	Cc          []Address
	Bcc         []Address
	Subject     string
	Comments    []string
	Keywords    []string
	ContentType string

	InReply     []string
	References  []string
}

// structure containing a multipart message
type Message struct {
	HeaderInfo

	Text  string
	Body  []byte
	Parts []Part
}

// structure containing a header value (respecting RFC5322 or not)
type Header struct {
	Key, Value string
}


// structure containing a raw header
type RawHeader struct {
	Key, Value []byte
}

// structure containing a raw message (arrays of bytes)
type RawMessage struct {
	RawHeaders []RawHeader
	Body       []byte
}

