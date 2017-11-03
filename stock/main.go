package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/fatih/color"
)

type Stock struct {
	Quote struct {
		Symbol           string  `json:"symbol"`
		CompanyName      string  `json:"companyName"`
		PrimaryExchange  string  `json:"primaryExchange"`
		Sector           string  `json:"sector"`
		CalculationPrice string  `json:"calculationPrice"`
		Open             float64 `json:"open"`
		OpenTime         int64   `json:"openTime"`
		Close            float64 `json:"close"`
		CloseTime        int64   `json:"closeTime"`
		LatestPrice      float64 `json:"latestPrice"`
		LatestSource     string  `json:"latestSource"`
		LatestTime       string  `json:"latestTime"`
		LatestUpdate     int64   `json:"latestUpdate"`
		LatestVolume     int     `json:"latestVolume"`
		IexRealtimePrice float64 `json:"iexRealtimePrice"`
		IexRealtimeSize  int     `json:"iexRealtimeSize"`
		IexLastUpdated   int64   `json:"iexLastUpdated"`
		DelayedPrice     float64 `json:"delayedPrice"`
		DelayedPriceTime int64   `json:"delayedPriceTime"`
		PreviousClose    float64 `json:"previousClose"`
		Change           float64 `json:"change"`
		ChangePercent    float64 `json:"changePercent"`
		IexMarketPercent float64 `json:"iexMarketPercent"`
		IexVolume        int     `json:"iexVolume"`
		AvgTotalVolume   int     `json:"avgTotalVolume"`
		IexBidPrice      float64 `json:"iexBidPrice"`
		IexBidSize       int     `json:"iexBidSize"`
		IexAskPrice      float64 `json:"iexAskPrice"`
		IexAskSize       int     `json:"iexAskSize"`
		MarketCap        int64   `json:"marketCap"`
		PeRatio          float64 `json:"peRatio"`
		Week52High       float64 `json:"week52High"`
		Week52Low        float64 `json:"week52Low"`
		YtdChange        float64 `json:"ytdChange"`
	} `json:"quote"`
	News []struct {
		Datetime string `json:"datetime"`
		Headline string `json:"headline"`
		Source   string `json:"source"`
		URL      string `json:"url"`
		Summary  string `json:"summary"`
		Related  string `json:"related"`
	} `json:"news"`
}

func fetchSymbols(symbols []string) ([]Stock, error) {
	if len(symbols) == 0 {
		return []Stock{}, nil
	}

	u, err := url.Parse("https://api.iextrading.com/1.0/stock/market/batch")
	if err != nil {
		panic(err)
	}
	v := url.Values{}
	v.Add("symbols", strings.Join(symbols, ","))
	v.Add("last", "1")
	v.Add("types", "quote")
	u.RawQuery = v.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	results := make(map[string]Stock)
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, err
	}

	stocks := make([]Stock, 0, len(symbols))
	for _, symbol := range symbols {
		stocks = append(stocks, results[symbol])
	}

	return stocks, err
}

var (
	hiWhite = color.New(color.FgHiWhite)
	hiRed   = color.New(color.FgHiRed)
	hiGreen = color.New(color.FgHiGreen)
)

func main() {
	stocks, err := fetchSymbols(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed get stock price: %v\n", err)
	}

	for _, s := range stocks {
		hiWhite.Printf("%8s: ", s.Quote.Symbol)
		fmt.Printf("%7.2f ", s.Quote.LatestPrice)

		changeColor := hiGreen
		if s.Quote.Change < 0 {
			changeColor = hiRed
		}
		changeColor.Printf("%+6.2f (%+0.2f%%)\n", s.Quote.Change, s.Quote.ChangePercent*100)
	}
}
