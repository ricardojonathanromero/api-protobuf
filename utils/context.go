package utils

import (
	"context"
	"time"
)

func ContextWithTimeout(seconds int) (context.Context, context.CancelFunc) {
	if seconds <= 0 {
		seconds = 10
	}

	return context.WithTimeout(context.Background(), time.Duration(seconds)*time.Second)
}
