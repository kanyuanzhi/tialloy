package test

import (
	"testing"
)

func TestClient1(t *testing.T) {
	client := NewClient(1, 1)
	client.Start()
}
