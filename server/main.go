package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type ExchangeDB struct {
	ID    int `gorm:"primaryKey"`
	Value float64
}

const url = "https://economia.awesomeapi.com.br/json/last/USD-BRL"

type InfoExchangeUSDBRL struct {
	Usdbrl struct {
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"varBid"`
		PctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	} `json:"USDBRL"`
}

func main() {
	http.HandleFunc("/cotacao", HandleExchangeInfo)
	http.ListenAndServe(":8080", nil)

}

func HandleExchangeInfo(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/cotacao" {
		w.WriteHeader(http.StatusNotFound)
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 200*time.Millisecond)
	defer cancel()

	info, err := ExchangeUSDBRL(ctx)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	exchangeFloat64, err := strconv.ParseFloat(info.Usdbrl.Bid, 64)
	if err != nil {
		panic(err)
	}

	ctx, cancel = context.WithTimeout(ctx, 10*time.Millisecond)
	defer cancel()

	exchange := ExchangeDB{Value: exchangeFloat64}

	insertExchange(ctx, &exchange)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	result, err := json.Marshal(exchange)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(result)
}

func ExchangeUSDBRL(ctx context.Context) (*InfoExchangeUSDBRL, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var info InfoExchangeUSDBRL
	err = json.Unmarshal(body, &info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

func insertExchange(ctx context.Context, exchange *ExchangeDB) {
	db, err := gorm.Open(sqlite.Open("cotacaoDB.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&ExchangeDB{})

	db.Create(&exchange)

}
