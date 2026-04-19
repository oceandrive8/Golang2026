package main

import (
	"fmt"
	"net/http"
	"sync"

	"Practice9/internal/handler"
	"Practice9/internal/middleware"
	"Practice9/internal/storage"
)

func main() {

	store := storage.NewStore()
	h := &handler.PaymentHandler{}

	mux := http.NewServeMux()

	mux.Handle("/pay",
		middleware.Idempotency(store, h),
	)

	go func() {
		http.ListenAndServe(":8080", mux)
	}()

	var wg sync.WaitGroup
	key := "loan-123"

	for i := 0; i < 10; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			req, _ := http.NewRequest("POST", "http://localhost:8080/pay", nil)
			req.Header.Set("Idempotency-Key", key)

			client := &http.Client{}
			resp, err := client.Do(req)

			if err != nil {
				fmt.Println(i, "error")
				return
			}

			fmt.Println(i, resp.StatusCode)
		}(i)
	}

	wg.Wait()
	fmt.Println("---- After completion ----")

	req, _ := http.NewRequest("POST", "http://localhost:8080/pay", nil)
	req.Header.Set("Idempotency-Key", key)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("final request error")
		return
	}

	fmt.Println("Final request status:", resp.StatusCode)
}
