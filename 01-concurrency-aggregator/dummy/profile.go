package dummy

import (
	"context"
	"fmt"
)

var profile = map[int]string{
	1: "John Doe",
	2: "Jane Smith",
	3: "Bob Johnson",
}

func GetProfile(ctx context.Context, id int) (string, error) {
	if profile[id] == "" {
		return "", fmt.Errorf("profile not found")
	}
	return profile[id], nil
}
