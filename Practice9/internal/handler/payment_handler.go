package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type PaymentHandler struct{}

func (h *PaymentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	time.Sleep(2 * time.Second)

	resp := map[string]interface{}{
		"status":         "paid",
		"amount":         1000,
		"transaction_id": uuid.NewString(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
