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
	self.books = make(map[string]*Book)
	self.TopBuy = make(map[string]int)
	self.LowSell = make(map[string]int)
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

func vwap(xs [][2]int) int {
	sum := 0
	vol := 0
	for _, x := range xs {
		sum += x[0] * x[1]
		vol += x[1]
	}
	return sum / vol
}

func (self *ArbStrategy) handle(message map[string]interface{}, orderId *int) (trades []interface{}) {
	var book *Book
	if message["type"] == "book" {
		book = BookFromMap(message)
	} else {
		return nil
	}
	if (len(book.Buy) == 0 || len(book.Sell) == 0) {
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
		buy := vwap(self.books[stock].Buy)
		sell := vwap(self.books[stock].Sell)
		underlying_fair_value := (buy + sell) / 2
		if stock == "BOND" {
			underlying_fair_value = 1000
		}
		self.fairValue += underlying_fair_value * self.weights[i]
	}
	self.fairValue = self.fairValue / 10

	self.TopBuy[self.composite] = self.books[self.composite].Buy[0][0]
	self.LowSell[self.composite] = self.books[self.composite].Sell[0][0]

	log.Println("fairValue: ", self.fairValue, " ", vwap(self.books[self.composite].Buy), vwap(self.books[self.composite].Sell))

	if vwap(self.books[self.composite].Buy) > self.fairValue {
		*orderId++
		margin := (self.TopBuy[self.composite] - vwap(self.books[self.composite].Buy)) / 2
		trades = append(trades, Order{
			Type:    "add",
			OrderId: *orderId,
			Symbol:  self.composite,
			Dir:     "SELL",
			Price:   self.fairValue + margin,
			Size:    10,
		})
	}

	if vwap(self.books[self.composite].Sell) < self.fairValue {
		*orderId++
		margin := (self.fairValue - vwap(self.books[self.composite].Sell)) / 2
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
	for k, _ := range self.books {
		self.books[k] = nil
	}
	self.fairValue = 0
	self.books_obtained = 0
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
