package test

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConnection(t *testing.T) {
	conn, err := net.Dial("tcp", "127.0.0.1:9004")
	assert.NoError(t, err)
	conn.Write([]byte("hello, remote"))
	time.Sleep(10 * time.Second)
	conn.Close()
}
