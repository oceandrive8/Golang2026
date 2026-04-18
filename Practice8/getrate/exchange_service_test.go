package getrate

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// 1. SUCCESS
func TestGetRate_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, _ = fmt.Fprint(w, `{"base":"USD","target":"EUR","rate":0.92}`)
	}))
	defer server.Close()

	service := NewExchangeService(server.URL)

	rate, err := service.GetRate("USD", "EUR")

	assert.NoError(t, err)
	assert.InEpsilon(t, 0.92, rate, 0.0001)
}

// 2. BUSINESS ERROR (400 / 404)
func TestGetRate_BusinessError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)

		_, _ = fmt.Fprint(w, `{"error":"invalid currency pair"}`)
	}))
	defer server.Close()

	service := NewExchangeService(server.URL)

	_, err := service.GetRate("USD", "XXX")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid currency pair")
}

// 3. MALFORMED JSON
func TestGetRate_MalformedJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		_, _ = fmt.Fprint(w, "NOT_JSON")
	}))
	defer server.Close()

	service := NewExchangeService(server.URL)

	_, err := service.GetRate("USD", "EUR")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "decode error")
}

// 4. TIMEOUT
func TestGetRate_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Second)
	}))
	defer server.Close()

	service := NewExchangeService(server.URL)

	_, err := service.GetRate("USD", "EUR")

	assert.Error(t, err)
}

// 5. SERVER ERROR (500)
func TestGetRate_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)

		_, _ = fmt.Fprint(w, `{"error":"server crashed"}`)
	}))
	defer server.Close()

	service := NewExchangeService(server.URL)

	_, err := service.GetRate("USD", "EUR")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "api error")
}

// 6. EMPTY BODY
func TestGetRate_EmptyBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	service := NewExchangeService(server.URL)

	_, err := service.GetRate("USD", "EUR")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "decode error")
}

// 7. SERVER PANIC
func TestGetRate_ServerPanic(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("mock server crash")
	}))
	defer server.Close()

	service := NewExchangeService(server.URL)

	_, err := service.GetRate("USD", "EUR")

	assert.Error(t, err)
}
