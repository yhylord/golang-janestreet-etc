package main

import (
	"encoding/json"
	"log"
	"net"
	"time"
)

/*
 These constants decide which exchange to use (various test ones or production).
*/
const TEST_MODE = false
const TEST_EXCHANGE_INDEX = 0

/*
 These constants are fixed regardless of config.
*/
// TODO: Check these constants when the real thing happens
const TEAM_NAME = "THEGRANDLIKEKING"
const BASE_PORT = 25000
const TEST_HOST = "test-exch"
const PROD_HOST = "production"

type Order struct {
	// use capitalized names to become public
	Type    string `json:"type"`
	OrderId int    `json:"order_id"`
	Symbol  string `json:"symbol"`
	Dir     string `json:"dir"`
	Price   int    `json:"price"`
	Size    int    `json:"size"`
}

func tcpConnect(host string) *net.Conn {
	const RETRY = 10 * time.Millisecond
	for {
		log.Println("Establishing connection to " + host)
		exchange, err := net.Dial("tcp", host)
		if err == nil {
			log.Println("Connection established.")
			return &exchange
		} else {
			log.Println("Connection failed. Retrying in " + RETRY.String())
			time.Sleep(RETRY)
		}
	}
}

func WriteToExchange(exchange *net.Conn) {
	jsonWriter := json.NewEncoder(exchange)
	jsonWriter.Encode()
}

func main() {
	log.SetFlags(log.Ltime | log.Lshortfile)
	tcpConnect(":8080")
}
