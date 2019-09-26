package bitflyer

import (
	"context"
	"testing"
	"time"

	"golang.org/x/sync/errgroup"
)

func TestRealtimeAPIClient(t *testing.T) {
	bf := NewBitflyer()
	realtime, err := bf.GetRealtimeAPIClient()
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return realtime.Subscribe(ctx)
	})

	time.Sleep(1 * time.Second)
	cancel()
	if err := eg.Wait(); err != nil {
		t.Fatal(err)
	}
}
