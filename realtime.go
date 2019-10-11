package bitflyer

import (
	"context"
	"encoding/json"
	"sort"
	"sync"
	"time"

	"github.com/sodefrin/wsjsonrpc"
	"golang.org/x/sync/errgroup"
)

type RealtimeAPIClient struct {
	rpc                 *wsjsonrpc.JsonRPC
	boardMu             *sync.Mutex
	board               *Board
	executionMu         *sync.Mutex
	executions          []*Execution
	onBoardCallback     func(float64, []*Price, []*Price)
	onExecutionCallback func([]*Execution)
}

var maxExecutions = 100000
var timeoutInterval = time.Minute * 1

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
	MidPrice float64  `json:"mid_price"`
	Bids     []*Price `json:"bids"`
	Asks     []*Price `json:"asks"`
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

	timer := time.NewTimer(timeoutInterval)

	eg := errgroup.Group{}
	eg.Go(func() error {
		for {
			if err := r.recv(); err != nil {
				return err
			}
			timer.Reset(timeoutInterval)
		}
	})

	eg.Go(func() error {
		select {
		case <-timer.C:
			return r.Close()
		case <-ctx.Done():
			return r.Close()
		}
	})

	return eg.Wait()
}

func (r *RealtimeAPIClient) AddOnBoardCallback(ctx context.Context, callback func(mid float64, bids []*Price, asks []*Price)) {
	r.onBoardCallback = callback
}

func (r *RealtimeAPIClient) AddOnExecutionCallback(ctx context.Context, callback func([]*Execution)) {
	r.onExecutionCallback = callback
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
			if err := r.updateBoard(boardMsg); err != nil {
				return err
			}
			if r.onBoardCallback != nil {
				r.onBoardCallback(boardMsg.MidPrice, boardMsg.Bids, boardMsg.Asks)
			}
			return nil
		case "lightning_executions_FX_BTC_JPY":
			executionMsg := executionMessage{}
			if err := json.Unmarshal(channelMsg.Message, &executionMsg); err != nil {
				break
			}
			if err := r.updateExecutions(executionMsg); err != nil {
				return err
			}
			if r.onExecutionCallback != nil {
				r.onExecutionCallback(executionMsg)
			}
		}
	}
	return nil
}

func (r *RealtimeAPIClient) updateBoard(msg *boardMessage) error {
	r.boardMu.Lock()
	r.board.midPrice = msg.MidPrice
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
		ts, err := parseTimeString(v.ExecDate)
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
