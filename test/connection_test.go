package test

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConnection(t *testing.T) {
	conn, err := net.Dial("tcp", "127.0.0.1:9006")
	assert.NoError(t, err)
	conn.Write([]byte("hello, remote"))
	time.Sleep(15 * time.Second)
	conn.Close()
}
