package bitflyer

import (
	"errors"
	"sync"
)

const (
	endpoint         = "https://api.bitflyer.com"
	realtimeEndpoint = "wss://ws.lightstream.bitflyer.com/json-rpc"
	origin           = "https://ws.lightstream.bitflyer.com/json-rpc"
)

type Bitflyer struct{}

func NewBitflyer() *Bitflyer {
	return &Bitflyer{}
}

var ErrInvalidStatusCode = errors.New("invalid status code")
var ErrInvalidResponse = errors.New("invalid response")

func (b *Bitflyer) GetRealtimeAPIClient() (*RealtimeAPIClient, error) {
	return &RealtimeAPIClient{
		boardMu:     &sync.Mutex{},
		board:       &Board{bids: map[float64]float64{}, asks: map[float64]float64{}},
		executionMu: &sync.Mutex{},
		executions:  []*Execution{},
	}, nil
}

func (b *Bitflyer) GetPublicAPIClient() (*PublicAPIClient, error) {
	return &PublicAPIClient{}, nil
}

func (b *Bitflyer) PrivateAPIClient(key, secret string) (*PrivateAPIClient, error) {
	return &PrivateAPIClient{key: key, secret: secret}, nil
}
