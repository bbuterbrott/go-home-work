package main

import (
	"fmt"
	"log"
	"time"

	"github.com/beevik/ntp"
)

const (
	ntpServer string = "ntp3.stratum2.ru"
)

func main() {
	t := time.Now()
	fmt.Printf("current time: %v\n", t.Round(0))
	ntpTime, err := ntp.Time(ntpServer)
	if err != nil {
		log.Fatalf("Got error while getting time from ntp server %v: %v", ntpServer, err)
	}
	fmt.Printf("exact time: %v\n", ntpTime.Round(0))
}
