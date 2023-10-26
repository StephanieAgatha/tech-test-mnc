package main

import (
	"log"
	"mnc-test/delivery"
	"mnc-test/util/helper"
)

func main() {
	helper.PrintAscii()
	s, err := delivery.NewServer()
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}
	s.Run()
}
