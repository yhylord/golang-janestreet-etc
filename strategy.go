package main

type Strategy struct {
	// you can store state here, e.g. countBuyorder int
	xTopBuy, xLowSell, yTopBuy, yLowSell                     int
	xTopBuyCount, xLowSellCount, yTopBuyCount, yLowSellCount int
	// margin                                                   int
	fairValue    int
	xBook, yBook *Book
}

// ADR: VALE - Y
// underlying: VALBZ - X

func (self *Strategy) handle(message map[string]interface{}, orderId *int) (trades []Order) {
	//default message["type"] should be "book"
	var book *Book
	if message["type"] == "book" {
		book = BookFromMap(message)
	} else {
		return nil
	}

	if book.Symbol == "VALBZ" {
		self.xBook = book
	}

	if book.Symbol == "VALE" {
		self.yBook = book
	}

	if self.xBook == nil || self.yBook == nil {
		return nil
	}

	//
	//!! haven't dealt with oderID & margin
	//

	self.xTopBuy = self.xBook.Buy[0].price
	self.xLowSell = self.xBook.Sell[0].price

	// for _, b := range book.Buy {
	// 	if b.price > self.xTopBuy {
	// 		self.xTopBuy = b.price
	// 	}
	// }
	// for _, s := range book.Sell {
	// 	if s.price < self.xLowSell {
	// 		self.xLowSell = s.price
	// 	}
	// }

	//calculate fair value based on xTopBuy and xLowSell
	self.fairValue = (self.xTopBuy + self.xLowSell) / 2

	if book.Symbol == "VALE" {
		self.yTopBuy = book.Buy[0].price
		self.yLowSell = book.Sell[0].price
		// for _, b := range book.Buy {
		// 	if b.price > self.yTopBuy {
		// 		self.yTopBuy = b.price
		// 	}
		// }
		// for _, s := range book.Sell {
		// 	if s.price < self.xLowSell {
		// 		self.yLowSell = s.price
		// 	}
		// }
	}

	//if yTopBuy > fairValue, sell Y
	//if yLowSell < fairValue, buy Y

	if self.yTopBuy > self.fairValue {
		margin := (self.yTopBuy - self.fairValue) / 2
		*orderId++
		trades = append(trades, Order{
			Type:    "add",
			OrderId: *orderId,
			Symbol:  "VALE",
			Dir:     "SELL",
			Price:   self.fairValue + margin,
			Size:    10,
		})
		self.xBook = nil
		self.yBook = nil
	}

	if self.yLowSell < self.fairValue {
		*orderId++
		margin := (self.fairValue - self.yLowSell) / 2
		trades = append(trades, Order{
			Type:    "add",
			OrderId: *orderId,
			Symbol:  "VALE",
			Dir:     "BUY",
			Price:   self.fairValue - margin,
			Size:    10,
		})
		self.xBook = nil
		self.yBook = nil
	}
	return trades
}
