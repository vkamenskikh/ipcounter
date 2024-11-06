package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

// const blockSizeMin int64 = 100
// const blockSizeMax int64 = 100

const blockSizeMin int64 = 32 * 1024 // Default buff reader size = 4096, do not set less
const blockSizeMax int64 = 100 * 1024 * 1024

func MultiRead(fileName string, topLayer LayerIf) int64 {

	fileInfo, err := os.Stat(fileName)
	if err != nil {
		panic(err)
	}
	numCPU := 4 //runtime.NumCPU()
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

	ipMerger := make(chan uint32, 1000)
	mergeResult := make(chan int64, 1)

	go func() {
		var rowCounter int64 = 0
		for ipInt := range ipMerger {
			rowCounter++
			topLayer.add(ipInt, 0)
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

func readBlockChain(fileName string, blockSize int64, blockNumbers <-chan int64, ipMerger chan<- uint32, wg *sync.WaitGroup) {
	defer wg.Done()

	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	for blockNumber := range blockNumbers {
		fmt.Println("Start block: ", blockNumber)

		_, err = file.Seek(blockSize*blockNumber, 0)
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		reader := bufio.NewReader(file)
		var readCount int64 = 0
		invalidCount := 0
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
					invalidCount++
					if len(strings.TrimSpace(line)) > 0 {
						fmt.Printf("Invalid IP: %v, ignoring\n", strings.TrimSpace(line))
					}
					if invalidCount > 5 {
						panic(errors.New("There are more than 5 ip parsing errors, last: " + line))
					}
				} else {
					ipMerger <- ipInt
				}
			}
			readCount += int64(len(line))
			if readCount > blockSize {
				break
			}
		}
	}
}
