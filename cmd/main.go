package main

import (
	"mnc-test/delivery"
	"mnc-test/util/helper"
)

func main() {
	helper.PrintAscii()
	delivery.NewServer().Run()
}
