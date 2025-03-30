package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type ExchangeRate struct {
	Bid string `json:"bid"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	exchangeRate, err := getExchangeRate(ctx)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Print("Timeout when getting exchange rate from server")
		} else {
			log.Printf("Error getting exchange rate: %v", err)
		}
		return
	}

	if exchangeRate == nil {
		log.Print("Received empty exchange rate data")
		return
	}

	if err := saveExchangeRateToFile(exchangeRate); err != nil {
		log.Printf("Error saving exchange rate to file: %v", err)
	}
}

func getExchangeRate(ctx context.Context) (*ExchangeRate, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080/cotacao", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusRequestTimeout {
		return nil, fmt.Errorf("timeout when getting the exchange rate")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d - %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	var exchangeRate ExchangeRate
	if err := json.NewDecoder(resp.Body).Decode(&exchangeRate); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &exchangeRate, nil
}

func saveExchangeRateToFile(exchangeRate *ExchangeRate) error {
	if exchangeRate == nil {
		return errors.New("exchange rate data is nil")
	}

	content := fmt.Sprintf("Dólar: %s\n", exchangeRate.Bid)
	log.Print("Saving exchange rate to file: ", content)

	if err := os.WriteFile("cotacao.txt", []byte(content), 0o644); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}
