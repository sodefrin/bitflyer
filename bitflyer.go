package bitflyer

import "github.com/sodefrin/wsjsonrpc"

const (
	endpoint         = "https://api.bitflyer.com"
	realtimeEndpoint = "wss://ws.lightstream.bitflyer.com/json-rpc"
	origin           = "https://ws.lightstream.bitflyer.com/json-rpc"
)

type Bitflyer struct{}

func NewBitflyer() *Bitflyer {
	return &Bitflyer{}
}

type RealtimeAPIClient struct {
	rpc *wsjsonrpc.JsonRPC
}

func (b *Bitflyer) GetRealtimeAPIClient() (*RealtimeAPIClient, error) {
	rpc, err := wsjsonrpc.NewJsonRPC(realtimeEndpoint, "", origin)
	if err != nil {
		return nil, err
	}
	return &RealtimeAPIClient{rpc: rpc}, nil
}

type PublicAPIClient struct{}

func (b *Bitflyer) GetPublicAPIClient() (*PublicAPIClient, error) {
	return &PublicAPIClient{}, nil
}

type PrivateAPIClient struct {
	key    string
	secret string
}

func (b *Bitflyer) PrivateAPIClient(key, secret string) (*PrivateAPIClient, error) {
	return &PrivateAPIClient{key: key, secret: secret}, nil
}
