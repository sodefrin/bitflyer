package bitflyer

import (
	"context"
	"encoding/json"
	"fmt"
)

type Board struct{}

func (r *RealtimeAPIClient) GetBoard() (*Board, error) {
	return &Board{}, nil
}

type channelMessage struct {
	Channel string          `json:"channel"`
	Message json.RawMessage `json:"message,omitempty"`
}

type boardMessage struct {
	MidPrice float64 `json:"mid_price"`
	Bids     []price `json:"bids"`
	Asks     []price `json:"asks"`
}

type price struct {
	Price float64 `json:"price"`
	Size  float64 `json:"size"`
}

type executionMessage struct {
	ID                         int64   `json:"id"`
	Side                       string  `json:"side"`
	Price                      float64 `json:"price"`
	Size                       float64 `json:"size"`
	ExecDate                   string  `json:"exec_date"`
	BuyChildOrderAcceptanceID  string  `json:"buy_child_order_acceptance_id"`
	SellChildOrderAcceptanceID string  `json:"sell_child_order_acceptance_id"`
}

func (r *RealtimeAPIClient) Subscribe(ctx context.Context) error {
	if err := r.rpc.Send("subscribe", &channelMessage{
		Channel: "lightning_board_FX_BTC_JPY",
	}, nil); err != nil {
		return err
	}

	if err := r.rpc.Send("subscribe", &channelMessage{
		Channel: "lightning_board_snapshot_FX_BTC_JPY",
	}, nil); err != nil {
		return err
	}

	if err := r.rpc.Send("subscribe", &channelMessage{
		Channel: "lightning_executions_FX_BTC_JPY",
	}, nil); err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			method, msg, _, err := r.rpc.Recv()
			if err != nil {
				return err
			}
			if method == "channelMessage" {
				channelMsg := channelMessage{}
				if err := json.Unmarshal(msg, &channelMsg); err != nil {
					return err
				}
				switch channelMsg.Channel {
				case "lightning_board_FX_BTC_JPY":
					boardMsg := &boardMessage{}
					if err := json.Unmarshal(channelMsg.Message, boardMsg); err != nil {
						break
					}
					r.initBoard(boardMsg)
				case "lightning_board_snapshot_FX_BTC_JPY":
					boardMsg := &boardMessage{}
					if err := json.Unmarshal(channelMsg.Message, boardMsg); err != nil {
						break
					}
					r.updateBoard(boardMsg)
				case "lightning_executions_FX_BTC_JPY":
					executionMsg := &executionMessage{}
					if err := json.Unmarshal(channelMsg.Message, executionMsg); err != nil {
						break
					}
					r.updateExecutions(executionMsg)
				}
			}
		}
	}
}

func (r *RealtimeAPIClient) initBoard(msg *boardMessage) error {
	fmt.Println(msg)
	return nil
}

func (r *RealtimeAPIClient) updateBoard(msg *boardMessage) error {
	fmt.Println(msg)
	return nil
}

func (r *RealtimeAPIClient) updateExecutions(msg *executionMessage) error {
	fmt.Println(msg)
	return nil
}

func (r *RealtimeAPIClient) Close() error {
	return r.rpc.Close()
}
