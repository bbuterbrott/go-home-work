package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/beevik/ntp"
)

type TimeError string

func (e TimeError) Error() string {
	return string(e)
}

// go test -gcflags=-l
func TestHelloNow(t *testing.T) {
	t.Run("test normal behavior", func(t *testing.T) {
		layout := "2 Jan 2006 15:04:05"

		monkey.Patch(time.Now, func() time.Time {
			nowTime, err := time.Parse(layout, "9 May 1945 10:03:00")
			if err != nil {
				t.Fatal(err)
			}
			return nowTime
		})

		monkey.Patch(ntp.Time, func(_ string) (time.Time, error) {
			ntpTime, err := time.Parse(layout, "9 May 1945 10:03:02")
			if err != nil {
				t.Fatal(err)
			}
			return ntpTime, nil
		})

		result, err := catchStdout(main)
		if err != nil {
			t.Fatal(err)
		}

		expected := `current time: 1945-05-09 10:03:00 +0000 UTC
exact time: 1945-05-09 10:03:02 +0000 UTC
`
		if string(result) != expected {
			t.Fatalf("invalid output:\n%v, expected:\n%s", string(result), expected)
		}
	})
}

func TestHelloNowError(t *testing.T) {
	t.Run("test error behavior", func(t *testing.T) {
		layout := "2 Jan 2006 15:04:05"

		monkey.Patch(time.Now, func() time.Time {
			nowTime, err := time.Parse(layout, "9 May 1945 10:03:00")
			if err != nil {
				t.Fatal(err)
			}
			return nowTime
		})

		monkey.Patch(ntp.Time, func(_ string) (time.Time, error) {
			ntpTime, err := time.Parse(layout, "9 May 1945 10:03:02")
			if err != nil {
				t.Fatal(err)
			}
			return ntpTime, TimeError("lookup ntp3.stratum2.ru on [::1]:53: read udp [::1]:55753->[::1]:53: read: connection refused")
		})

		// Если этого не делать то программа выйдет до того, как завершится тест. Наверное, exit code можно как-то более правильно проверить, но я только до такого решения додумался
		monkey.Patch(log.Fatalf, func(string, ...interface{}) {})

		result, err := catchStdout(main)
		if err != nil {
			t.Fatal(err)
		}

		expected := `current time: 1945-05-09 10:03:00 +0000 UTC
`
		errorMes := `Got error while getting time from ntp server`
		if strings.HasPrefix(string(result), expected) && strings.Contains(string(result), errorMes) {
			t.Fatalf("invalid output:\n%s, expected:\n%s", result, expected)
		}
	})
}

func catchStdout(runnable func()) (result []byte, err error) {
	realOut := os.Stdout
	defer func() { os.Stdout = realOut }()

	r, fakeOut, err := os.Pipe()
	if err != nil {
		return
	}

	os.Stdout = fakeOut
	runnable()
	if err = fakeOut.Close(); err != nil {
		return
	}

	result, err = ioutil.ReadAll(r)
	if err != nil {
		return
	}

	err = r.Close()
	return
}
