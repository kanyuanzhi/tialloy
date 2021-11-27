package test

import (
	"testing"
	"time"
)

func TestClient1(t *testing.T) {
	client := NewClient(1, 1)
	client.Start()
	time.Sleep(5 * time.Second)
	client.Conn.Close()
	time.Sleep(100 * time.Second)

}
