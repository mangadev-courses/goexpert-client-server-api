package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"client-server-api/pkg/db"
)

type CurrencyData struct {
	Code       string `json:"code"`
	CodeIn     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

type APIResponse struct {
	USDBRL CurrencyData `json:"USDBRL"`
}

func (r *APIResponse) validate() error {
	if err := r.USDBRL.validate(); err != nil {
		return err
	}

	return nil
}

func (c *CurrencyData) validate() error {
	if c.Bid == "" {
		baseErr := errors.New("bid is missing")
		return fmt.Errorf("validation failed: %w; context: %+v", baseErr, c)
	}

	return nil
}

func ExchangeRate(dbClient *sql.DB, apiUrl string, w http.ResponseWriter, r *http.Request) {
	ctxAPI, cancelAPI := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancelAPI()

	request, err := http.NewRequestWithContext(ctxAPI, "GET", apiUrl, nil)
	if err != nil {
		http.Error(w, "Failed to create request: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		http.Error(w, "Failed to fetch data: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		http.Error(w, "Failed to read response body: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var result APIResponse
	if err := json.Unmarshal(body, &result); err != nil {
		http.Error(w, "Failed to parse response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := result.validate(); err != nil {
		http.Error(w, "Invalid response data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	ctxDB, cancelDB := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancelDB()

	err = db.InsertBid(ctxDB, dbClient, result.USDBRL.Bid)
	if err != nil {
		http.Error(w, "Failed to save data: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(result.USDBRL.Bid))
}

func Healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func Start(dbClient *sql.DB, ready chan struct{}) {
	http.HandleFunc("/healthz", Healthz)
	http.HandleFunc("/cotacao", func(w http.ResponseWriter, r *http.Request) {
		ExchangeRate(dbClient, "https://economia.awesomeapi.com.br/json/last/usd-brl", w, r)
	})

	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Failed to start server: %v\n", err)
		}
	}()

	for {
		resp, err := http.Get("http://localhost:8080/healthz")
		if err == nil && resp.StatusCode == http.StatusOK {
			break
		}

		fmt.Println("Waiting for server to start...")
		time.Sleep(200 * time.Millisecond)
	}

	fmt.Println("Server started on port 8080")
	close(ready)
}
