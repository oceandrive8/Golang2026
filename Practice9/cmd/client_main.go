package main

import (
	"context"
	"net/http"
	"time"

	"Practice9/internal/client"
	"Practice9/internal/server"
)

func main() {
	mockServer := server.NewMockServer()
	defer mockServer.Close()

	paymentClient := client.PaymentClient{
		Client: &http.Client{
			Timeout: 3 * time.Second,
		},
		MaxRetries: 5,
		URL:        mockServer.URL,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := paymentClient.ExecutePayment(ctx)
	if err != nil {
		println("Error:", err.Error())
	}
}
