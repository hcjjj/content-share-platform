package time

import (
	"context"
	"testing"
	"time"
)

func TestTimer(t *testing.T) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Second*10)
	defer cancel()
	timer := time.NewTimer(time.Second)
	defer timer.Stop()
	for {
		select {
		case now := <-timer.C:
			t.Log(now.String())
		case <-ctx.Done():
			// 退出
			return
		}
	}
}
