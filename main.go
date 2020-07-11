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
	const RETRY = 10 * time.Millisecond
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

/*func ReadFromExchange(exchange net.Conn) {
	reader := json.NewDecoder(exchange)
	err := reader.Decode()
}
*/
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
	WriteToExchange(exchange, Hello{
		Type: "hello",
		Team: TEAM_NAME,
	})
}
