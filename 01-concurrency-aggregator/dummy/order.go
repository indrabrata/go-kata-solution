package dummy

import (
	"context"
	"fmt"
)

var orders = map[int]int{
	1: 10,
	2: 20,
	3: 30,
}

func GetOrder(ctx context.Context, id int) (int, error) {
	if orders[id] == 0 {
		return 0, fmt.Errorf("order not found")
	}
	return orders[id], nil
}
