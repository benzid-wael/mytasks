package cli

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

func CaptureOutput(f func()) string {
	reader, writer, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	stdout := os.Stdout
	stderr := os.Stderr
	defer func() {
		os.Stdout = stdout
		os.Stderr = stderr
		log.SetOutput(os.Stderr)
	}()
	os.Stdout = writer
	os.Stderr = writer
	log.SetOutput(writer)
	out := make(chan string)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		var buf bytes.Buffer
		wg.Done()
		io.Copy(&buf, reader)
		out <- buf.String()
	}()
	wg.Wait()
	f()
	writer.Close()
	return <-out
}

func Difference(a, b time.Time) (year, month, day, hour, min, sec int) {
	if a.Location() != b.Location() {
		b = b.In(a.Location())
	}
	if a.After(b) {
		a, b = b, a
	}
	y1, M1, d1 := a.Date()
	y2, M2, d2 := b.Date()

	h1, m1, s1 := a.Clock()
	h2, m2, s2 := b.Clock()

	year = int(y2 - y1)
	month = int(M2 - M1)
	day = int(d2 - d1)
	hour = int(h2 - h1)
	min = int(m2 - m1)
	sec = int(s2 - s1)

	// Normalize negative values
	if sec < 0 {
		sec += 60
		min--
	}
	if min < 0 {
		min += 60
		hour--
	}
	if hour < 0 {
		hour += 24
		day--
	}
	if day < 0 {
		// days in month:
		t := time.Date(y1, M1, 32, 0, 0, 0, 0, time.UTC)
		day = t.Day() + day
		month--
	}
	if month < 0 {
		month += 12
		year--
	}

	return
}

func GetDurationText(a, b time.Time) string {
	year, month, day, hour, min, sec := Difference(a, b)
	duration := "now"
	if year > 0 {
		duration = fmt.Sprintf("%v years", year)
	} else if month > 0 {
		duration = fmt.Sprintf("%v months", month)
	} else if day > 0 {
		duration = fmt.Sprintf("%vd", day)
	} else if hour > 0 {
		duration = fmt.Sprintf("%vh", hour)
	} else if min > 0 {
		duration = fmt.Sprintf("%vm", min)
	} else if sec > 0 {
		duration = fmt.Sprintf("%vs", sec)
	}
	return duration
}
