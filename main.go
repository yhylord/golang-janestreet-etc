package main

import (
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
const TEAM_NAME = "GRANDLIKEKING"
const BASE_PORT = 25000
const TEST_HOST = "test-exch"
const PROD_HOST = "production"

func tcpConnect(host string) *net.Conn {
	const RETRY = 10 * time.Millisecond
	for {
		log.Println("Establishing connection to " + host)
		c, err := net.Dial("tcp", host)
		if err == nil {
			log.Println("Connection established.")
			return &c
		} else {
			log.Println("Connection failed. Retrying in " + RETRY.String())
			time.Sleep(RETRY)
		}
	}
}

func main() {
	tcpConnect(":8080")
}
