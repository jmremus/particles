package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/tarm/serial"
)

const PORT = "8080"
const DEV = "/dev/ttyUSB0"

// this looks weird, but it's to get 80 samples in a 90 minute period
// 80 samples == assume terminal width at least 80
// 90 minutes == a decent timeframe
const sampleSize = 80
const SleepSecs = 67

// let's make a simple ring buffer here
// instead of a queue...
var samples [sampleSize]string
var sampleIdx = 0
var buf string
var mu sync.Mutex

var c = &serial.Config{Name: DEV, Baud: 115200}
var s, porterr = serial.OpenPort(c)

func parseAndAddToRing(s string) {
	mu.Lock()
	split := strings.Split(s, ",")

	if (len(split) >= 3) && (split[2] == "PC0.3") {
		samples[sampleIdx] = split[3]
	}

	// Roll over our circular buffer if needed
	if sampleIdx == (sampleSize - 1) {
		sampleIdx = 0
	} else {
		sampleIdx += 1
	}

	mu.Unlock()
}

func pollSerial() {
	scanner := bufio.NewScanner(s)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		mu.Lock()
		buf = scanner.Text()
		mu.Unlock()
		parseAndAddToRing(buf)
		time.Sleep(SleepSecs * time.Second)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
}

func writeParticles(w http.ResponseWriter, _ *http.Request) {
	mu.Lock()
	io.WriteString(w, buf)

	for i := sampleIdx; i < sampleSize; i++ {
		io.WriteString(w, samples[i]+",")
	}

	for i := 0; i < sampleIdx-1; i++ {
		io.WriteString(w, samples[i]+",")
	}

	// write the last line without a commma
	io.WriteString(w, samples[sampleIdx-1])

	mu.Unlock()
}

func main() {
	// init ring
	for i := 0; i < sampleSize; i++ {
		samples[i] = "0"
	}

	if porterr != nil {
		log.Fatal(porterr)
	}

	go pollSerial()

	http.HandleFunc("/", writeParticles)
	log.Fatal(http.ListenAndServe(":"+PORT, nil))
}
