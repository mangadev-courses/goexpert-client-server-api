package client

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func FetchExchangeRate(url string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2000*time.Millisecond)
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return errors.New("Failed to create request: " + err.Error())
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return errors.New("Filed to fetch data: " + err.Error())
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return errors.New("Failed to read response body: " + err.Error())
	}

	if response.StatusCode != http.StatusOK {
		err := string(body)
		return fmt.Errorf("unexpected status code: %d; body: %s", response.StatusCode, err)
	}

	bid := string(body)
	formattedBid := fmt.Sprintf("Dólar: %s", bid)

	os.WriteFile("cotacao.txt", []byte(formattedBid), 0644)
	return nil
}
