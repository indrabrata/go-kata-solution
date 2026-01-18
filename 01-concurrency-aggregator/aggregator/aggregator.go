package aggregator

import (
	"log/slog"
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

func WithTimeout(timemout int) Option {
	return func(aggr *UserAggregator) {
		aggr.timeout = timemout
	}
}

func WithLogger(logger *slog.Logger) Option {
	return func(aggr *UserAggregator) {
		aggr.logger = logger
	}
}

func (s *UserAggregator) Aggregate(id int) (string, error) {
	return "", nil
}
