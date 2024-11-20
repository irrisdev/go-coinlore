package coinlore

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	ticker  = "api/ticker/?id=%d"
	tickers = "api/tickers/?start=%d&limit=%d"
	global  = "api/global/"
)

// NewClient creates a new CoinLore API client
func NewClient(url string) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		baseURL:    url,
	}
}

// GetCoins fetches the list of coins from CoinLore
func (c *Client) GetCoins(start int, limit int) (*[]Coin, error) {
	endpoint := fmt.Sprintf(tickers, start, limit)
	resp, err := c.httpClient.Get(fmt.Sprintf("%s/%s", c.baseURL, endpoint))
	if err != nil {
		return nil, &ClientError{
			Status:  http.StatusInternalServerError,
			Message: fmt.Sprintf("failed to make request: %v", err),
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result CoinListResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, &ClientError{
			Status:  http.StatusInternalServerError,
			Message: fmt.Sprintf("failed to decode response: %v", err),
		}
	}

	return &result.Data, nil
}

// GetCoin fetches a single coin by its ID
func (c *Client) GetCoin(id int) (*Coin, error) {

	endpoint := fmt.Sprintf(ticker, id)
	resp, err := c.httpClient.Get(fmt.Sprintf("%s/%s", c.baseURL, endpoint))
	if err != nil {
		return nil, &ClientError{
			Status:  http.StatusInternalServerError,
			Message: fmt.Sprintf("failed to make request: %v", err),
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, &ClientError{
			Status:  resp.StatusCode,
			Message: fmt.Sprintf("unexpected status code: %d", resp.StatusCode),
		}
	}

	var result []Coin

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, &ClientError{
			Status:  http.StatusInternalServerError,
			Message: fmt.Sprintf("failed to decode response: %v", err),
		}
	}

	if len(result) == 0 {
		return nil, &ClientError{
			Status:  http.StatusNotFound,
			Message: fmt.Sprintf("coin with ID %d not found", id),
		}
	}

	return &result[0], nil

}

// GetGlobal fetches the global data from CoinLore
func (c *Client) GetGlobal() (*Global, error) {
	resp, err := c.httpClient.Get(fmt.Sprintf("%s/%s", c.baseURL, global))
	if err != nil {
		return nil, &ClientError{
			Status:  http.StatusInternalServerError,
			Message: fmt.Sprintf("failed to make request: %v", err),
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, &ClientError{
			Status:  resp.StatusCode,
			Message: fmt.Sprintf("unexpected status code: %d", resp.StatusCode),
		}
	}

	var result []Global

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, &ClientError{
			Status:  http.StatusInternalServerError,
			Message: fmt.Sprintf("failed to decode response: %v", err),
		}
	}

	if len(result) == 0 {
		return nil, &ClientError{
			Status:  http.StatusInternalServerError,
			Message: "no global data found",
		}
	}

	return &result[0], nil
}
