package main

type Strategy struct {
	// you can store state here, e.g. countBuyorder int
	xTopBuy, xLowSell, yTopBuy, yLowSell                     int
	xTopBuyCount, xLowSellCount, yTopBuyCount, yLowSellCount int
	margin                                                   int
	fairValue                                                int
	xMessagesSize, yMessagesSize                             int
	/*	xMessages                                                []*Book
		yMessages                                                []*Book
	*/
}

// ADR: VALE - Y
// underlying: VALBZ - X

func (self *Strategy) handle(message map[string]interface{}, orderId *int) (trades []Order) {
	//default message["type"] should be "book"
	book := BookFromMap(message)

	//
	//!! haven't dealt with oderID & margin
	//

	for _, b := range book.Buy {
		if b.price > self.xTopBuy {
			self.xTopBuy = b.price
		}
	}
	for _, s := range book.Sell {
		if s.price < self.xLowSell {
			self.xLowSell = s.price
		}
	}

	if book.Symbol == "VALE" {
		self.yTopBuy = book.Buy[0].price
		self.yLowSell = book.Sell[0].price
		for _, b := range book.Buy {
			if b.price > self.yTopBuy {
				self.yTopBuy = b.price
			}
		}
		for _, s := range book.Sell {
			if s.price < self.xLowSell {
				self.yLowSell = s.price
			}
		}
	}
	self.fairValue = (self.xTopBuy + self.xLowSell) / 2
	self.margin = (self.yTopBuy - self.fairValue) / 2

	//if yTopBuy > fairValue, sell Y
	//if yLowSell < fairValue, buy Y

	if self.yTopBuy > self.fairValue {
		*orderId++
		trades = append(trades, Order{
			Type:    "add",
			OrderId: *orderId,
			Symbol:  "VALE",
			Dir:     "SELL",
			Price:   self.fairValue + self.margin,
			Size:    10,
		})
	}

	if self.yLowSell < self.fairValue {
		*orderId++
		trades = append(trades, Order{
			Type:    "add",
			OrderId: *orderId,
			Symbol:  "VALE",
			Dir:     "BUY",
			Price:   self.fairValue - self.margin,
			Size:    10,
		})
	}
	return trades
}
