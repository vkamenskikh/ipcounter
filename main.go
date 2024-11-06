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

func PrintMemUsage(title string) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Println("======= Memory ", title, "=======")
	fmt.Printf("Alloc = %v k\n", bToMb(m.Alloc)) // KiB
	fmt.Printf("TotalAlloc = %v k\n", bToMb(m.TotalAlloc))
	fmt.Printf("Sys = %v k\n", bToMb(m.Sys))
	fmt.Println("------------------------")
}
func bToMb(b uint64) uint64 {
	return b / 1024
}
