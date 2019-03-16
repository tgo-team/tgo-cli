package client

import (
	"github.com/tgo-team/tgo-talk/test"
	"testing"
)

func TestClient_Register(t *testing.T) {
	c := New(NewOptions())
	status,err := c.Register(1,"123456")
	test.Nil(t,err)
	test.Equal(t,uint16(200),status)
}

func TestClient_Login(t *testing.T) {
	c := New(NewOptions())
	err := c.Login(1,"123456")
	test.Nil(t,err)
}