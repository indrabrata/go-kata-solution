package aggregator

import (
	"concurrency-aggregator/model"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"
)

type UserAggregator struct {
	timeout int
	logger  *slog.Logger
}

type Option func(*UserAggregator)

func NewUserAggregator(opts ...Option) *UserAggregator {
	aggr := &UserAggregator{
		timeout: 10,
		logger:  slog.Default(),
	}

	for _, opt := range opts {
		opt(aggr)
	}

	return aggr
}

func WithTimeout(timeout int) Option {
	return func(aggr *UserAggregator) {
		aggr.timeout = timeout
	}
}

func WithLogger(logger *slog.Logger) Option {
	return func(aggr *UserAggregator) {
		aggr.logger = logger
	}
}

func (s *UserAggregator) Aggregate(ctx context.Context, id int) (string, error) {
	client := http.Client{
		Timeout: time.Duration(s.timeout) * time.Second,
	}

	var profile model.Profile
	var orders model.Order

	// Note : When one goroutine returns an error, the context gets cancelled, which should cancel the other goroutine's HTTP request
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		url := fmt.Sprintf("https://yourcoolapi.com/profiles/%d", id)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return err
		}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(body, &profile); err != nil {
			return err
		}

		s.logger.Info("profile fetched", "id", id, "body", string(body))
		return nil
	})

	g.Go(func() error {
		url := fmt.Sprintf("https://yourcoolapi.com/orders/%d", id)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return err
		}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(body, &orders); err != nil {
			return err
		}

		s.logger.Info("orders fetched", "id", id, "body", string(body))
		return nil
	})

	if err := g.Wait(); err != nil {
		return "", err
	}

	s.logger.Info("aggregation completed", "id", id)

	return fmt.Sprintf("Name : %s | Orders : %d", profile.Name, orders.Orders), nil
}
