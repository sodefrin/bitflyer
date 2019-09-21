package bitflyer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Ticker struct {
	ProductCode     string    `json:"product_code"`
	Timestamp       time.Time `json:"timestamp"`
	TickID          int       `json:"tick_id"`
	BestBid         float64   `json:"best_bid"`
	BestAsk         float64   `json:"best_ask"`
	BestBidSize     float64   `json:"best_bid_size"`
	BestAskSize     float64   `json:"best_ask_size"`
	TotalBidDepth   float64   `json:"total_bid_depth"`
	TotalAskDepth   float64   `json:"total_ask_depth"`
	Ltp             float64   `json:"ltp"`
	Volume          float64   `json:"volume"`
	VolumeByProduct float64   `json:"volume_by_product"`
}

func (p *PublicAPIClient) GetTicker(productCode string) (*Ticker, error) {
	res, err := http.Get(endpoint + "/v1/ticker?product_code=" + productCode)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		bytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("%w; invalid status code (want: %d, have: %d)", err, 200, res.StatusCode)
		}
		return nil, fmt.Errorf("invalid status code (want: %d, have: %d, msg: %s)", 200, res.StatusCode, string(bytes))
	}

	t := struct {
		ProductCode     string  `json:"product_code"`
		Timestamp       string  `json:"timestamp"`
		TickID          int     `json:"tick_id"`
		BestBid         float64 `json:"best_bid"`
		BestAsk         float64 `json:"best_ask"`
		BestBidSize     float64 `json:"best_bid_size"`
		BestAskSize     float64 `json:"best_ask_size"`
		TotalBidDepth   float64 `json:"total_bid_depth"`
		TotalAskDepth   float64 `json:"total_ask_depth"`
		Ltp             float64 `json:"ltp"`
		Volume          float64 `json:"volume"`
		VolumeByProduct float64 `json:"volume_by_product"`
	}{}

	if err := json.NewDecoder(res.Body).Decode(&t); err != nil {
		return nil, err
	}

	timestamp, err := parseTimeString(t.Timestamp)
	if err != nil {
		return nil, err
	}

	return &Ticker{
		ProductCode:     t.ProductCode,
		Timestamp:       timestamp,
		TickID:          t.TickID,
		BestBid:         t.BestBid,
		BestAsk:         t.BestAsk,
		BestBidSize:     t.BestBidSize,
		BestAskSize:     t.BestAskSize,
		TotalBidDepth:   t.TotalBidDepth,
		TotalAskDepth:   t.TotalAskDepth,
		Ltp:             t.Ltp,
		Volume:          t.Volume,
		VolumeByProduct: t.VolumeByProduct,
	}, nil
}
