package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

type ExchangeRate struct {
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

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("mysql", "root:root@tcp(localhost:3306)/goexpert")
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal("Database is unreachable:", err)
	}

	http.HandleFunc("/cotacao", handleExchangeRate)
	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleExchangeRate(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	exchangeRate, err := getExchangeRate(ctx)
	if err != nil {
		handleError(w, ctx, err, "Timeout when obtaining the exchange rate", "request timeout", http.StatusRequestTimeout)
		return
	}

	ctxDB, cancelDB := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancelDB()

	if err := saveExchangeRate(ctxDB, *exchangeRate); err != nil {
		handleError(w, ctxDB, err, "Timeout when saving exchange rate to database", "database timeout", http.StatusRequestTimeout)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(exchangeRate)
	log.Println("Request Successful")
}

func getExchangeRate(ctx context.Context) (*ExchangeRate, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data map[string]ExchangeRate
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	exchangeRate, ok := data["USDBRL"]
	if !ok {
		return nil, errors.New("exchange rate data not found")
	}

	return &exchangeRate, nil
}

func saveExchangeRate(ctx context.Context, exchangeRate ExchangeRate) error {
	stmt, err := db.PrepareContext(ctx, `INSERT INTO exchange_rate (id, code, code_in, name, high, low, var_bid, pct_change, bid, ask, timestamp, create_date) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	id := uuid.New().String()
	log.Println("Generated UUID:", id)
	_, err = stmt.ExecContext(ctx, id, exchangeRate.Code, exchangeRate.CodeIn, exchangeRate.Name,
		exchangeRate.High, exchangeRate.Low, exchangeRate.VarBid, exchangeRate.PctChange,
		exchangeRate.Bid, exchangeRate.Ask, exchangeRate.Timestamp, exchangeRate.CreateDate)

	return err
}

func handleError(w http.ResponseWriter, ctx context.Context, err error, timeoutMsg, userMsg string, statusCode int) {
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		log.Print(timeoutMsg)
		http.Error(w, userMsg, statusCode)
	} else {
		log.Printf("Error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
