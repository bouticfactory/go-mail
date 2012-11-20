package mail

import (
	"testing"
)

func TestNewTextMessage(t *testing.T) {

	m, err := NewTextMessage("test", "bonjour", "donotreply@mydomain.com", "donotreceive@mydomain.com")

	server, _ := InitLocalHost("donotreply@mydomain.com")
	if (err == nil) {
		server.Send(m)
	}
}