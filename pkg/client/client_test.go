package client

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestClientCall(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("5.12333"))
	}))
	defer server.Close()
	os.Remove("cotacao.txt")
	defer os.Remove("cotacao.txt")

	FetchExchangeRate(server.URL)

	_, err := os.Stat("cotacao.txt")
	if err != nil {
		t.Errorf("Expected file cotacao.txt to exist, got error: %v", err)
	}

	file, err := os.ReadFile("cotacao.txt")
	if err != nil {
		t.Errorf("Failed to read cotacao.txt: %v", err)
	}

	expected := "Dólar: 5.12333"
	if string(file) != expected {
		t.Errorf("Expected file content %s, got %s", expected, string(file))
	}
}
