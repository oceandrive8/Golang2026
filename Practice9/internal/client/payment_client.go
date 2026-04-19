package client

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"Practice9/internal/retry"
)

type PaymentClient struct {
	Client     *http.Client
	MaxRetries int
	URL        string
}

func (p *PaymentClient) ExecutePayment(ctx context.Context) error {
	for attempt := 1; attempt <= p.MaxRetries; attempt++ {

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.URL, nil)
		if err != nil {
			return err
		}

		resp, err := p.Client.Do(req)

		if resp != nil {
			defer resp.Body.Close()
		}

		if err == nil && resp.StatusCode == http.StatusOK {
			fmt.Printf("Attempt %d: Success!\n", attempt)
			return nil
		}

		if !retry.IsRetryable(resp, err) {
			return fmt.Errorf("non-retryable error")
		}

		if attempt == p.MaxRetries {
			return fmt.Errorf("max retries exceeded")
		}

		delay := retry.CalculateBackoff(attempt)

		fmt.Printf("Attempt %d failed: waiting %v...\n", attempt, delay)

		select {
		case <-time.After(delay):
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return fmt.Errorf("unexpected exit")
}
