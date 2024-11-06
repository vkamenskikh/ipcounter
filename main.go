package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"time"
)

// 2 arguments: 1 - file for parsing, 2 - optional, set to 1 to use single reading thread
func main() {
	args := os.Args
	if len(args) <= 1 {
		panic(errors.New("missing file name as a program argument"))
	}
	sourceFile := args[1]
	multyReadRun := true
	if len(args) > 2 && args[2] == "1" {
		multyReadRun = false
	}

	PrintMemUsage("Before")

	start := time.Now()
	topLayer := Layer5{}

	var rows int64
	if multyReadRun {
		rows = MultiRead(sourceFile, &topLayer)
	} else {
		rows = SimpleRead(sourceFile, &topLayer) // Single thread reading for testing
	}

	elapsed1 := time.Since(start)
	countUnique := topLayer.count()
	elapsed2 := time.Since(start)

	log.Printf("Reading took %s, all: %s\n", elapsed1, elapsed2)
	fmt.Printf("Result, rows: %v, unique IPs: %v\n", rows, countUnique)

	PrintMemUsage("After")
	fmt.Println("Done")
}

func SimpleRead(fileName string, topLayer LayerIf) int64 {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	var lineCount int64 = 0
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error reading file", err)
			}
			break
		}
		ipInt, err := parseToInt(line)
		if err != nil {
			panic(err)
		}
		topLayer.add(ipInt, 0)
		lineCount++
	}
	return lineCount
}

const mbDiv = 1024 * 1024

func PrintMemUsage(title string) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Memory %v, Alloc: %v MiB, TotalAlloc: %v MiB, Sys: %v MiB\n", title, m.Alloc/mbDiv, m.TotalAlloc/mbDiv, m.Sys/mbDiv)
}
