# bitflyer

[![GoDoc](https://godoc.org/github.com/sodefrin/bitflyer?status.svg)](https://godoc.org/github.com/sodefrin/bitflyer)
[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/sodefrin/bitflyer/master/LICENSE)

bitflyer api for trading bot.

this liibrary contains

- realtime api
- private api
- public api

## Install

```
$ go get -u github.com/sodefrin/bitflyer
```

requirements: go1.13

## Usage

### realtime api

#### using ticker

```go
package main

import (
	"context"
	"log"
	"time"

	"github.com/sodefrin/bitflyer"
)

func main() {
	bf := bitflyer.NewBitflyer()
	realtime, err := bf.GetRealtimeAPIClient()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	go realtime.Subscribe(ctx)

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		<-ticker.C
		mid, bids, ask := realtime.GetBoard()
		// you can recieve board data.
		// ...

		exs := realtime.GetExecutions(time.Second)
		// you can recieve execution data within 1 second.
		// ...
	}
}
```

#### using callback

```go
package main

import (
	"context"
	"log"

	"github.com/sodefrin/bitflyer"
)

func main() {
	bf := bitflyer.NewBitflyer()
	realtime, err := bf.GetRealtimeAPIClient()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	realtime.AddOnBoardCallback(ctx, func(mid float64, bids, asks []*bitflyer.Price) {
		// you can recieve board diff data by callback.
		// ...
	})
	realtime.AddOnExecutionCallback(ctx, func(exs []*bitflyer.Execution) {
		// you can recieve execution diff data by callback.
		// ...
	})

	// blocking.
	if err := realtime.Subscribe(ctx); err != nil {
		log.Fatal(err)
	}
}
```

## parivate api

### create order

```go
package main

import (
	"log"

	"github.com/sodefrin/bitflyer"
)

func main() {
	bf := bitflyer.NewBitflyer()
	private, err := bf.PrivateAPIClient("your api key", "your api secret")
	if err != nil {
		log.Fatal(err)
	}

	if _, err = private.CreateOrder("BUY", 910000, 0.01, "LIMIT"); err != nil {
		log.Fatal(err)
	}
}
```
