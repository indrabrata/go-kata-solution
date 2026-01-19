package test

import (
	"concurrency-aggregator/aggregator"
	"context"
	"log/slog"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
)

var profileMap = map[string]interface{}{
	"id":   1,
	"name": "Alice",
}

var orderMap = map[string]interface{}{
	"id":     1,
	"orders": 5,
}

func TestAggregatorSuccess(t *testing.T) {
	httpmock.Activate(t)

	httpmock.RegisterResponder("GET", "https://yourcoolapi.com/profiles/1",
		func(r *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, profileMap)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		})

	httpmock.RegisterResponder("GET", "https://yourcoolapi.com/orders/1",
		func(r *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, orderMap)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		})

	userAggregator := aggregator.NewUserAggregator(
		aggregator.WithTimeout(5),
		aggregator.WithLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil))),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	resp, err := userAggregator.Aggregate(ctx, 1)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if resp != "Name : Alice | Orders : 5" {
		t.Errorf("unexpected response: %v", resp)
	}
}

func TestAggregatorTimeout(t *testing.T) {
	httpmock.Activate(t)

	httpmock.RegisterResponder("GET", "https://yourcoolapi.com/profiles/1",
		func(r *http.Request) (*http.Response, error) {
			time.Sleep(2 * time.Second)
			resp, err := httpmock.NewJsonResponse(200, profileMap)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		})

	httpmock.RegisterResponder("GET", "https://yourcoolapi.com/orders/1",
		func(r *http.Request) (*http.Response, error) {
			time.Sleep(2 * time.Second)
			resp, err := httpmock.NewJsonResponse(200, orderMap)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		})

	userAggregator := aggregator.NewUserAggregator(
		aggregator.WithTimeout(5),
		aggregator.WithLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil))),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	resp, err := userAggregator.Aggregate(ctx, 1)

	if err == nil {
		t.Errorf("expected error")
	}
	if resp != "" {
		t.Errorf("unexpected response: %v", resp)
	}
}

func TestAggregatorOneServiceError(t *testing.T) {
	httpmock.Activate(t)

	orderServiceCalled := false

	httpmock.RegisterResponder("GET", "https://yourcoolapi.com/profiles/1",
		httpmock.NewStringResponder(500, `{"error": "internal server error"}`))

	httpmock.RegisterResponder("GET", "https://yourcoolapi.com/orders/1",
		func(r *http.Request) (*http.Response, error) {
			select {
			case <-r.Context().Done():
				return nil, r.Context().Err()
			case <-time.After(2 * time.Second):
				orderServiceCalled = true
				resp, err := httpmock.NewJsonResponse(200, orderMap)
				if err != nil {
					return httpmock.NewStringResponse(500, ""), nil
				}
				return resp, nil
			}
		})

	userAggregator := aggregator.NewUserAggregator(
		aggregator.WithTimeout(5),
		aggregator.WithLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil))),
	)

	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	resp, err := userAggregator.Aggregate(ctx, 1)
	elapsed := time.Since(start)

	if err == nil {
		t.Errorf("expected error")
	}
	if resp != "" {
		t.Errorf("unexpected response: %v", resp)
	}
	// Should complete quickly (not wait 2s for orders service)
	if elapsed > 500*time.Millisecond {
		t.Errorf("expected fast failure, but took %v", elapsed)
	}
	if orderServiceCalled {
		t.Errorf("orders service should have been cancelled")
	}
}
