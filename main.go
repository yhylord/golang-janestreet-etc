package main

import (
	"encoding/json"
	"flag"
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

func tcpConnect(host string) net.Conn {
	for {
		log.Println("Establishing connection to " + host)
		exchange, err := net.Dial("tcp", host)
		if err == nil {
			log.Println("Connection established.")
			return exchange
		} else {
			log.Println("Connection failed. Retrying in " + RETRY.String())
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

func main() {
	log.SetFlags(log.Ltime | log.Lshortfile)
	prod := flag.Bool("production", false, "production mode")
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
	orderId := 0
	for {
		orderId++
		WriteToExchange(exchange, Order{
			Type:    "add",
			OrderId: orderId,
			Symbol:  "BOND",
			Dir:     "buy",
			Price:   999,
			Size:    10,
		})
		orderId++
		WriteToExchange(exchange, Order{
			Type:    "add",
			OrderId: orderId,
			Symbol:  "BOND",
			Dir:     "sell",
			Price:   1001,
			Size:    10,
		})
		var message map[string]interface{}
		filled := 0
		for filled < 2 {
			ReadFromExchange(exchange, &message)
			if message["type"] == "fill" {
				filled++
			}
			time.Sleep(time.Millisecond)
		}
		log.Println("Penny pinching filled!")
	}
}
