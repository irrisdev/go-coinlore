package coinlore

import (
	"net/http"
)

// Client is the API client for CoinLore
type Client struct {
	httpClient *http.Client
	baseURL    string
}

// Coin represents a cryptocurrency
type Coin struct {
	// Basic Information
	ID     string `json:"id"`
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
	Rank   int    `json:"rank"`

	// Price and Market Data
	PriceUSD     string  `json:"price_usd"`
	MarketCapUSD string  `json:"market_cap_usd"`
	Volume24     float64 `json:"volume24"`

	// Supply Information
	CirculatingSupply string `json:"csupply"`
	TotalSupply       string `json:"tsupply"`
	MaxSupply         string `json:"msupply"`

	// Percent Changes
	PercentChange1h  string `json:"percent_change_1h"`
	PercentChange24h string `json:"percent_change_24h"`
	PercentChange7d  string `json:"percent_change_7d"`
}

// CoinListResponse represents the response from the CoinLore API for a list of coins
type CoinListResponse struct {
	Data []Coin `json:"data"`
}

// Global represents the global cryptocurrency market data
type Global struct {
	CoinsCount   int    `json:"coins_count"`
	BtcDominance string `json:"btc_d"`
	McapChange   string `json:"mcap_change"`
	VolumeChange string `json:"volume_change"`
}

// ClientError represents an error when calling the CoinLore API
type ClientError struct {
	Status  int
	Message string
}

// Error returns the error message
func (e *ClientError) Error() string {
	return e.Message
}
