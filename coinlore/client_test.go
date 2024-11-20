package coinlore

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupMockServer() *httptest.Server {
	handler := http.NewServeMux()

	handler.HandleFunc("/api/tickers/", func(w http.ResponseWriter, r *http.Request) {
		response := CoinListResponse{
			Data: []Coin{
				{ID: "1", Symbol: "BTC", Name: "Bitcoin", PriceUSD: "50000"},
				{ID: "2", Symbol: "ETH", Name: "Ethereum", PriceUSD: "4000"},
			},
		}
		json.NewEncoder(w).Encode(response)
	})

	handler.HandleFunc("/api/ticker/", func(w http.ResponseWriter, r *http.Request) {
		response := []*Coin{
			{ID: "1", Symbol: "BTC", Name: "Bitcoin", PriceUSD: "50000"},
		}
		json.NewEncoder(w).Encode(response)
	})

	handler.HandleFunc("/api/global/", func(w http.ResponseWriter, r *http.Request) {
		response := []Global{
			{
				CoinsCount:   5000,
				McapChange:   "10%",
				VolumeChange: "5%",
				BtcDominance: "60%",
			},
		}
		json.NewEncoder(w).Encode(response)
	})

	return httptest.NewServer(handler)
}

func TestGetCoins(t *testing.T) {
	server := setupMockServer()
	defer server.Close()

	client := NewClient(server.URL)
	coins, err := client.GetCoins(0, 2)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(*coins) != 2 {
		t.Fatalf("expected 2 coins, got %d", len(*coins))
	}

	if (*coins)[0].Symbol != "BTC" {
		t.Errorf("expected BTC, got %s", (*coins)[0].Symbol)
	}
}

func TestGetCoin(t *testing.T) {
	server := setupMockServer()
	defer server.Close()

	client := NewClient(server.URL)
	coin, err := client.GetCoin(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if coin.Symbol != "BTC" {
		t.Errorf("expected BTC, got %s", coin.Symbol)
	}
}

func TestGetGlobal(t *testing.T) {
	server := setupMockServer()
	defer server.Close()

	client := NewClient(server.URL)
	global, err := client.GetGlobal()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if global.CoinsCount != 5000 {
		t.Errorf("expected 5000, got %d", global.CoinsCount)
	}
}

func TestGetCoins_NetworkError(t *testing.T) {
	client := NewClient("http://invalid-url")
	_, err := client.GetCoins(0, 2)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestGetCoins_Non200StatusCode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not found", http.StatusNotFound)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	_, err := client.GetCoins(0, 2)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestGetCoins_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	_, err := client.GetCoins(0, 2)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestGetCoin_NetworkError(t *testing.T) {
	client := NewClient("http://invalid-url")
	_, err := client.GetCoin(1)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestGetCoin_Non200StatusCode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not found", http.StatusNotFound)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	_, err := client.GetCoin(1)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestGetCoin_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	_, err := client.GetCoin(1)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestGetGlobal_NetworkError(t *testing.T) {
	client := NewClient("http://invalid-url")
	_, err := client.GetGlobal()
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestGetGlobal_Non200StatusCode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not found", http.StatusNotFound)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	_, err := client.GetGlobal()
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestGetGlobal_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	_, err := client.GetGlobal()
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestGetGlobal_NoData(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[]"))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	_, err := client.GetGlobal()
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if clientErr, ok := err.(*ClientError); ok {
		if clientErr.Message != "no global data found" {
			t.Errorf("expected message 'no global data found', got %s", clientErr.Message)
		}
	} else {
		t.Fatalf("expected ClientError, got %v", err)
	}
}

func TestClientError_Error(t *testing.T) {
	err := &ClientError{
		Status:  http.StatusInternalServerError,
		Message: "internal server error",
	}

	expected := "error: internal server error"
	if err.Error() != expected {
		t.Errorf("expected %s, got %s", expected, err.Error())
	}
}
