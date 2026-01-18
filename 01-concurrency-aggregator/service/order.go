package service

import "context"

type OrderService struct {
	// Define fields here
}

func NewOrderService() *OrderService {
	return &OrderService{}
}

func (s *OrderService) GetOrder(ctx context.Context, id int) (int, error) {
	return 5, nil
}
