package bitflyer

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"time"
)

var maxExecutions = 100000

func (r *RealtimeAPIClient) GetBoard() (float64, []*Price, []*Price) {
	bids := []*Price{}
	asks := []*Price{}

	r.boardMu.Lock()
	mid := r.board.midPrice
	for k, v := range r.board.bids {
		bids = append(bids, &Price{Price: k, Size: v})
	}
	for k, v := range r.board.asks {
		asks = append(asks, &Price{Price: k, Size: v})
	}
	r.boardMu.Unlock()

	sort.Slice(bids, func(i, j int) bool {
		return bids[i].Price > bids[j].Price
	})

	sort.Slice(asks, func(i, j int) bool {
		return asks[i].Price < asks[j].Price
	})

	return mid, asks, bids
}

func (r *RealtimeAPIClient) GetExecutions(duration time.Duration) []*Execution {
	exs := []*Execution{}
	start := time.Now().Add(-duration)
	r.executionMu.Lock()
	i := sort.Search(len(r.executions), func(i int) bool {
		return r.executions[i].Timestamp.After(start)
	})
	for j := i; j < len(r.executions); j++ {
		exs = append(exs, r.executions[j])
	}
	r.executionMu.Unlock()

	return exs
}

type Board struct {
	midPrice float64
	bids     map[float64]float64
	asks     map[float64]float64
}

type channelMessage struct {
	Channel string          `json:"channel"`
	Message json.RawMessage `json:"message,omitempty"`
}

type boardMessage struct {
	MidPrice float64 `json:"mid_price"`
	Bids     []Price `json:"bids"`
	Asks     []Price `json:"asks"`
}

type Price struct {
	Price float64 `json:"price"`
	Size  float64 `json:"size"`
}

type executionMessage []*Execution

type Execution struct {
	ID                         int64     `json:"id"`
	Side                       string    `json:"side"`
	Price                      float64   `json:"price"`
	Size                       float64   `json:"size"`
	ExecDate                   string    `json:"exec_date"`
	Timestamp                  time.Time `json:"timestamp"`
	BuyChildOrderAcceptanceID  string    `json:"buy_child_order_acceptance_id"`
	SellChildOrderAcceptanceID string    `json:"sell_child_order_acceptance_id"`
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
			if err := r.recv(); err != nil {
				fmt.Println(err)
			}
		}
	}
}

func (r *RealtimeAPIClient) recv() error {
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
		case "lightning_board_FX_BTC_JPY", "lightning_board_snapshot_FX_BTC_JPY":
			boardMsg := &boardMessage{}
			if err := json.Unmarshal(channelMsg.Message, boardMsg); err != nil {
				break
			}
			return r.updateBoard(boardMsg)
		case "lightning_executions_FX_BTC_JPY":
			executionMsg := executionMessage{}
			if err := json.Unmarshal(channelMsg.Message, &executionMsg); err != nil {
				fmt.Println(err)
				break
			}
			return r.updateExecutions(executionMsg)
		}
	}

	return nil
}

func (r *RealtimeAPIClient) updateBoard(msg *boardMessage) error {
	r.board.midPrice = msg.MidPrice

	r.boardMu.Lock()
	for _, v := range msg.Asks {
		if v.Size == 0 {
			delete(r.board.asks, v.Price)
			continue
		}
		r.board.asks[v.Price] = v.Size
	}

	for _, v := range msg.Bids {
		if v.Size == 0 {
			delete(r.board.bids, v.Price)
			continue
		}
		r.board.bids[v.Price] = v.Size
	}
	r.boardMu.Unlock()

	return nil
}

func (r *RealtimeAPIClient) updateExecutions(msg []*Execution) error {
	r.executionMu.Lock()
	for _, v := range msg {
		ts, err := parseBfTime(v.ExecDate)
		if err != nil {
			return err
		}
		v.Timestamp = ts
		r.executions = append(r.executions, v)
	}
	if len(r.executions) > maxExecutions {
		r.executions = r.executions[len(r.executions)-maxExecutions:]
	}
	r.executionMu.Unlock()
	return nil
}

func (r *RealtimeAPIClient) Close() error {
	return r.rpc.Close()
}

func parseBfTime(str string) (time.Time, error) {
	tmp, err := time.Parse("2006-01-02T15:04:05", str[:19])
	if err != nil {
		return tmp, err
	}
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	return tmp.In(jst), nil
}
