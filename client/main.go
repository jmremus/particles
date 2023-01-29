package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

const url = "http://192.168.1.1:8080/"
const stepSize = 10000

var levels = []string{"▁", "▂", "▃", "▄", "▅", "▆", "▇", "█"}

func spark(splits []string) {
	fmt.Println("PC0.3 particles in the last 90 minutes:")
	lenSplits := len(splits)

	for i := 30; i < lenSplits; i++ {
		splitInt, err := strconv.Atoi(splits[i])
		if err != nil {
			fmt.Println(err)
			break
		}

		level := splitInt / stepSize

		// Above 70,000 particles is the max step
		if level > 7 {
			level = 7
		}

		fmt.Print(levels[level])
	}
	fmt.Println()
}

func prettyPrint(splits []string) {
	for i := 0; i < 5; i++ {
		fmt.Print(splits[2*i], ": ", splits[2*i+1])
		if i < 4 {
			fmt.Print("  /  ")
		}
	}
	fmt.Println()
}

func main() {

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	splits := strings.Split(string(res), ",")
	spark(splits)
	prettyPrint(splits)
}
