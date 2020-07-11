package main

import "log"

type BondStrategy struct {
}

func (self *BondStrategy) handle(message map[string]interface{}, orderId *int, bonds *int) (trades []Order) {
	if message["type"] == "fill" {
		var buy_filled, sell_filled int
		if message["dir"] == "BUY" {
			buy_filled = int(message["size"].(float64))
			*orderId++
			trades = append(trades, Order{
				Type:    "add",
				OrderId: *orderId,
				Symbol:  "BOND",
				Dir:     "BUY",
				Price:   999,
				Size:    buy_filled,
			})
			*bonds += buy_filled
		}
		if message["dir"] == "SELL" {
			sell_filled = int(message["size"].(float64))
			half := sell_filled / 2
			*orderId++
			trades = append(trades, Order{
				Type:    "add",
				OrderId: *orderId,
				Symbol:  "BOND",
				Dir:     "SELL",
				Price:   1001,
				Size:    half,
			})
			*orderId++
			trades = append(trades, Order{
				Type:    "add",
				OrderId: *orderId,
				Symbol:  "BOND",
				Dir:     "SELL",
				Price:   1002,
				Size:    sell_filled - half,
			})
			*bonds -= sell_filled
		}
		log.Printf("Buy filled: %v, Sell filled: %v, Currently holding: %v\n", buy_filled, sell_filled, *bonds)
	}

	if message["type"] == "trade" && message["symbol"] == "BOND" {
		log.Println(message["size"], "bonds traded at ", message["price"])
	}

	return trades
}
