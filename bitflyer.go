package bitflyer

import (
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
