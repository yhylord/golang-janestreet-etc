package main

import "log"

type ArbStrategy struct {
	underlying []string
	weights    []int
	composite  string
	gotBooks   map[string]bool
	// you can store state here, e.g. countBuyorder int
	TopBuy, LowSell map[string]int

	// margin                                                   int
	fairValue      int
	books          map[string]*Book
	books_obtained int
}

func NewArbStrategy(underlying []string, weights []int, composite string) *ArbStrategy {
	self := new(ArbStrategy)
	self.underlying = underlying
	self.weights = weights
	self.composite = composite
	sefl.books := make(map[string]*Book)
	for _, stock := range underlying {
		self.books[stock] = nil
	}
	self.books_obtained = 0
	return self
}

// 1 ETF(XLF) contains a basket of followings
// 3 BOND
// 2 GS
// 3 MS
// 2 WFC

func (self *ArbStrategy) handle(message map[string]interface{}, orderId *int) (trades []interface{}) {
	var book *Book
	if message["type"] == "book" {
		book = BookFromMap(message)
	} else {
		return nil
	}

	if stringInSlice(book.Symbol, self.underlying) {
		if self.books[book.Symbol] == nil {
			self.books_obtained++
		}
		self.books[book.Symbol] = book
	}

	if book.Symbol == self.composite {
		self.books[self.composite] = book
	}

	if self.books_obtained < len(self.underlying) || self.books[self.composite] == nil {
		return nil
	}

	for i, stock := range self.underlying {
		self.TopBuy[stock] = self.books[stock].Buy[0][0]
		self.LowSell[stock] = self.books[stock].Sell[0][0]
		underlying_fair_value := (self.TopBuy[stock] + self.LowSell[stock]) / 2
		if stock == "BOND" {
			underlying_fair_value = 1000
		}
		self.fairValue += underlying_fair_value * self.weights[i]
	}
	self.fairValue = self.fairValue / 10

	self.TopBuy[self.composite] = self.books[self.composite].Buy[0][0]
	self.TopBuy[self.composite] = self.books[self.composite].Sell[0][0]

	if self.TopBuy[self.composite] > self.fairValue {
		*orderId++
		margin := (self.TopBuy[self.composite] - self.fairValue) / 4 * 3
		trades = append(trades, Order{
			Type:    "add",
			OrderId: *orderId,
			Symbol:  self.composite,
			Dir:     "SELL",
			Price:   self.fairValue + margin,
			Size:    10,
		})
	}

	if self.LowSell[self.composite] < self.fairValue {
		*orderId++
		margin := (self.fairValue - self.LowSell[self.composite]) / 4 * 3
		trades = append(trades, Order{
			Type:    "add",
			OrderId: *orderId,
			Symbol:  self.composite,
			Dir:     "BUY",
			Price:   self.fairValue - margin,
			Size:    10,
		})
	}

	if len(trades) != 0 {
		log.Println("ETF trades: ", trades)
	}
	return trades
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
