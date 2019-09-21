package bitflyer

import (
	"testing"
)

func TestPublicAPI(t *testing.T) {
	p := &PublicAPIClient{}

	ticker, err := p.GetTicker("FX_BTC_JPY")
	if err != nil {
		t.Fatal(err)
	}

	if ticker.ProductCode != "FX_BTC_JPY" {
		t.Fatalf("invalid ProductCode want FX_BTC_JPY have %s", ticker.ProductCode)
	}
}
