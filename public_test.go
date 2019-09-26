package bitflyer

import (
	"testing"
)

func TestPublicAPI(t *testing.T) {
	bf := NewBitflyer()
	public, err := bf.GetPublicAPIClient()
	if err != nil {
		t.Fatal(err)
	}

	ticker, err := public.GetTicker("FX_BTC_JPY")
	if err != nil {
		t.Fatal(err)
	}

	if ticker.ProductCode != "FX_BTC_JPY" {
		t.Fatalf("invalid ProductCode want FX_BTC_JPY have %s", ticker.ProductCode)
	}
}
