package tialloy

import "testing"

func TestClient2(t *testing.T) {
	client := NewClient(2,2)
	client.Start()
}
