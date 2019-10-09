package bitflyer

import (
	"errors"
	"sync"

	"github.com/sodefrin/wsjsonrpc"
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
var ErrTimeout = errors.New("timeout")

func (b *Bitflyer) GetRealtimeAPIClient() (*RealtimeAPIClient, error) {
	rpc, err := wsjsonrpc.NewJsonRPC("2.0", realtimeEndpoint, origin)
	if err != nil {
		return nil, err
	}

	return &RealtimeAPIClient{
		rpc:         rpc,
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
