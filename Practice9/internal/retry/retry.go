package retry

import (
	"math/rand"
	"net"
	"net/http"
	"time"
)

var BaseDelay = 500 * time.Millisecond

func init() {
	rand.Seed(time.Now().UnixNano())
}

func IsRetryable(resp *http.Response, err error) bool {
	if err != nil {
		if netErr, ok := err.(net.Error); ok {
			return netErr.Timeout() || netErr.Temporary()
		}
		return false
	}

	switch resp.StatusCode {
	case 429, 500, 502, 503, 504:
		return true
	case 401, 404:
		return false
	default:
		return false
	}
}

func CalculateBackoff(attempt int) time.Duration {
	maxDelay := BaseDelay * time.Duration(1<<(attempt-1))
	return time.Duration(rand.Int63n(int64(maxDelay)))
}
