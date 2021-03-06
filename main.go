package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"
)

/*
 These constants decide which exchange to use (various test ones or production).
*/
const TEST_EXCHANGE_INDEX = 0

/*
 These constants are fixed regardless of config.
*/
// TODO: Check these constants when the real thing happens
const TEAM_NAME = "THEGRANDLIKEKING"
const BASE_PORT = 25000
const TEST_HOST = "test-exch-THEGRANDLIKEKING"
const PROD_HOST = "production"

const RETRY = 10 * time.Millisecond

type Hello struct {
	Type string `json:"type"`
	Team string `json:"team"`
}

type Order struct {
	// use capitalized names to become public
	Type    string `json:"type"`
	OrderId int    `json:"order_id"`
	Symbol  string `json:"symbol"`
	Dir     string `json:"dir"`
	Price   int    `json:"price"`
	Size    int    `json:"size"`
}

type Pair struct {
	price, size int
}

type Book struct {
	Type   string   `json:"type"`
	Symbol string   `json:"symbol"`
	Buy    [][2]int `json:"buy"`
	Sell   [][2]int `json:"sell"`
}

func BookFromMap(message map[string]interface{}) *Book {
	book := new(Book)
	bytes, _ := json.Marshal(message)
	json.Unmarshal(bytes, book)
	return book
}

// func NewOrder(order_type, symbol, dir string, price, size int) *Order {

// 	return &Order{}
// }

func tcpConnect(host string) net.Conn {
	for {
		log.Println("Establishing connection to " + host)
		exchange, err := net.Dial("tcp", host)
		if err == nil {
			log.Println("Connection established.")
			return exchange
		} else {
			log.Println("Connection failed. ", err, "Retrying in "+RETRY.String())
			time.Sleep(RETRY)
		}
	}
}

func ReadFromExchange(exchange net.Conn, message interface{}) error {
	reader := json.NewDecoder(exchange)
	err := reader.Decode(message)
	return err
}

func WriteToExchange(exchange net.Conn, message interface{}) error {
	writer := json.NewEncoder(exchange)
	err := writer.Encode(message)
	exchange.Write([]byte("\n"))
	return err
}

func putTrades(exchange net.Conn, orders []interface{}) {
	for _, o := range orders {
		WriteToExchange(exchange, o)
	}
}

func main() {
	log.SetFlags(log.Ltime | log.Lshortfile)
	prod := flag.Bool("prod", false, "production mode")
	flag.Parse()
	log.Println("prodL: ", *prod)
	var host string
	if *prod {
		host = PROD_HOST + ":" + strconv.Itoa(BASE_PORT)
	} else {
		host = TEST_HOST + ":" + strconv.Itoa(BASE_PORT)
	}
	exchange := tcpConnect(host)
	err := WriteToExchange(exchange, Hello{
		Type: "hello",
		Team: TEAM_NAME,
	})
	if err != nil {
		log.Printf("Failed to hello, error: %v\n", err)
	}
	var message map[string]interface{}
	err_read := ReadFromExchange(exchange, &message)
	if err_read == nil {
		fmt.Println(message)
	}
	orderId := 0
	bonds := 0
	strategy := new(Strategy)
	etfStrategy := NewArbStrategy([]string{"BOND", "GS", "MS", "WFC"}, []int{3, 2, 3, 2}, "XLF")
	for {
		var err1, err2 error
		for i := 0; i < 6; i++ {
			orderId++
			err1 = WriteToExchange(exchange, Order{
				Type:    "add",
				OrderId: orderId,
				Symbol:  "BOND",
				Dir:     "BUY",
				Price:   999,
				Size:    5,
			})

			orderId++
			err2 = WriteToExchange(exchange, Order{
				Type:    "add",
				OrderId: orderId,
				Symbol:  "BOND",
				Dir:     "SELL",
				Price:   1002,
				Size:    5,
			})
		}
		if err1 == nil && err2 == nil {
			var message map[string]interface{}
			for {
				ReadFromExchange(exchange, &message)
				if message["type"] == "fill" {
					var buy_filled, sell_filled int
					if message["dir"] == "BUY" {
						buy_filled = int(message["size"].(float64))
						orderId++
						WriteToExchange(exchange, Order{
							Type:    "add",
							OrderId: orderId,
							Symbol:  "BOND",
							Dir:     "BUY",
							Price:   999,
							Size:    buy_filled,
						})
						bonds += buy_filled
					}
					if message["dir"] == "SELL" {
						sell_filled = int(message["size"].(float64))
						half := sell_filled / 2
						orderId++
						WriteToExchange(exchange, Order{
							Type:    "add",
							OrderId: orderId,
							Symbol:  "BOND",
							Dir:     "SELL",
							Price:   1001,
							Size:    half,
						})
						orderId++
						WriteToExchange(exchange, Order{
							Type:    "add",
							OrderId: orderId,
							Symbol:  "BOND",
							Dir:     "SELL",
							Price:   1002,
							Size:    sell_filled - half,
						})
						bonds -= sell_filled
					}
					log.Printf("Buy filled: %v, Sell filled: %v, Currently holding: %v\n", buy_filled, sell_filled, bonds)
				}

				if message["type"] == "trade" && message["symbol"] == "BOND" {
					log.Println(message["size"], "bonds traded at ", message["price"])
				}

				putTrades(exchange, strategy.handle(message, &orderId))
				putTrades(exchange, etfStrategy.handle(message, &orderId))

				//bondStrategy.handle(message, &orderId)

				time.Sleep(RETRY / 2)
			}
		} else {
			log.Println("Error for buy order: ", err1)
			log.Println("Error for sell order: ", err2)
		}
	}
}
