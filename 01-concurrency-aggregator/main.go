package main

import (
	"concurrency-aggregator/aggregator"
	"fmt"
	"log/slog"
	"os"
)

func main() {
	userAggregator := aggregator.NewUserAggregator(
		aggregator.WithTimeout(20),
		aggregator.WithLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil))),
	)

	fmt.Println(userAggregator.Aggregate(1))
}
