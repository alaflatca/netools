package test

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConnection(t *testing.T) {
	conn, err := net.Dial("tcp", "127.0.0.1:9001")
	assert.NoError(t, err)
	time.Sleep(10 * time.Second)
	conn.Close()
}
