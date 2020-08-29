package main

import (
	"fmt"
	"log"
	"time"

	"github.com/beevik/ntp"
)

const (
	// Если не указывать string (хоть он тут и не нужен), то не работает Monkey.Patch на ntp.Time() в тесте ¯\_(ツ)_/¯
	ntpServer string = "ntp3.stratum2.ru"
)

func main() {
	t := time.Now()
	fmt.Printf("current time: %v\n", t.Round(0))
	ntpTime, err := ntp.Time(ntpServer)
	if err == nil {
		fmt.Printf("exact time: %v\n", ntpTime.Round(0))
	} else {
		log.Fatalf("Got error while getting time from ntp server %v: %v", "ntp3.stratum2.ru", err)
	}
}
