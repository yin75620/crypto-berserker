package main

import (
	"log"
	"strconv"
	"time"
)

func main() {
	nanos := time.Now().UnixNano() / 1000000
	ts := strconv.FormatInt(nanos, 10)
	log.Println(ts)
}
