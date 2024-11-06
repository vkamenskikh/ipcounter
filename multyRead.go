package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"
)

const blockSizeMin int64 = 32 * 1024 // Default buff reader size = 4096, do not set less
const blockSizeMax int64 = 100 * 1024 * 1024
const maxErrorCount = 5
const accumLen = 512

var parserErrorCount = 0

// We get performance advantages with concurrent reading multiple file sections, if it is stored on SSD
// Old style HDD works slower with concurrent reads comparing with 1 read because the HDD head must move around.
func MultiRead(fileName string, topLayer LayerIf) int64 {

	fileInfo, err := os.Stat(fileName)
	if err != nil {
		panic(err)
	}
	numCPU := runtime.NumCPU()
	fileSize := fileInfo.Size()

	blockSize := 1 + fileSize/int64(numCPU)
	if blockSize < blockSizeMin {
		blockSize = blockSizeMin
	}
	if blockSize > blockSizeMax {
		blockSize = blockSizeMax
	}
	blockCount := 1 + fileSize/blockSize
	threadCount := int(blockCount)
	if threadCount > numCPU {
		threadCount = numCPU
	}
	fmt.Println("Number of CPUs: ", numCPU, " reading threads: ", threadCount)

	ipMerger := make(chan []uint32, 50)
	mergeResult := make(chan int64, 1)

	go func() {
		var rowCounter int64 = 0
		for ipAccum := range ipMerger {
			rowCounter += int64(len(ipAccum))
			for _, ipInt := range ipAccum {
				topLayer.add(ipInt, 0)
			}
		}
		mergeResult <- rowCounter
	}()

	var wg sync.WaitGroup
	wg.Add(threadCount)

	blockNumbers := getBlockSender(blockCount)

	for ii := 0; ii < threadCount; ii++ {
		go readBlockChain(fileName, blockSize, blockNumbers, ipMerger, &wg)
	}
	wg.Wait()
	close(ipMerger)

	rows := <-mergeResult
	close(mergeResult)
	return rows
}

func getBlockSender(blockCount int64) <-chan int64 {
	blockNumbers := make(chan int64, 10)
	go func() {
		for ii := int64(0); ii < blockCount; ii++ {
			blockNumbers <- ii
		}
		close(blockNumbers)
	}()
	return blockNumbers
}

func readBlockChain(fileName string, blockSize int64, blockNumbers <-chan int64, ipMerger chan<- []uint32, wg *sync.WaitGroup) {
	defer wg.Done()

	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	for blockNumber := range blockNumbers {
		//fmt.Println("Start block: ", blockNumber)

		_, err = file.Seek(blockSize*blockNumber, 0)
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		reader := bufio.NewReader(file)
		var readCount int64 = 0
		accum := make([]uint32, 0, accumLen)

		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					break
				}
				panic(err)
			}
			if readCount > 0 || blockNumber == 0 {
				ipInt, err := parseToInt(line) // ipInt
				if err != nil {
					procParseError(line, err)
				} else {
					accum = append(accum, ipInt)
					if len(accum) >= accumLen {
						ipMerger <- accum
						accum = make([]uint32, 0, accumLen)
					}
				}
			}
			readCount += int64(len(line))
			if readCount > blockSize {
				break
			}
		}
		if len(accum) > 0 {
			ipMerger <- accum
		}
	}
}

func procParseError(line string, err error) {
	parserErrorCount++
	if len(strings.TrimSpace(line)) > 0 {
		fmt.Printf(err.Error() + ", ignoring\n")
	}
	if parserErrorCount >= maxErrorCount {
		panic(errors.New("There are  5 ip parsing errors, last one: " + err.Error()))
	}

}
