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

	time.Sleep(3 * time.Second)
	mid, asks, bids := realtime.GetBoard()
	if asks[0].Price < mid {
		t.Fatalf("asks[0].Price > bids[0].Price")
	}
	if mid < bids[0].Price {
		t.Fatalf("asks[0].Price > bids[0].Price")
	}

	exs := realtime.GetExecutions(3 * time.Second)
	if len(exs) == 0 {
		t.Fatal("len(exs) must not be zero")
	}

	cancel()
	if err := eg.Wait(); err != nil {
		t.Fatal(err)
	}
}
