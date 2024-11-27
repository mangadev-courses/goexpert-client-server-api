package server

import (
	"client-server-api/pkg/db"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServerCotacaoEndpoint(t *testing.T) {
	mockDB, _ := sql.Open("sqlite3", ":memory:")
	db.CreateTable(mockDB)
	defer mockDB.Close()

	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"USDBRL":{"bid":"5.12333"}}`))
	}))
	defer mockAPI.Close()

	req := httptest.NewRequest("GET", "/cotacao", nil)
	w := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ExchangeRate(mockDB, mockAPI.URL, w, r)
	})

	handler.ServeHTTP(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}

	body := w.Body.String()
	if body != "5.12333" {
		t.Errorf("Expected body 5.123, got %s", body)
	}
}
