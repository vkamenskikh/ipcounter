package main

import (
	"bufio"
	"fmt"
)

// Testing approach, not used
func Accumulate(scanner *bufio.Scanner) {

	data := make([][][][]bool, 256)
	lineNum := 0
	ip := [4]uint8{0, 0, 0, 0}
	for scanner.Scan() {
		line := scanner.Text()
		if err := parse4(line, &ip); err != nil {
			panic(err)
		}
		//fmt.Println("New", ip[0], ip[1], ip[2], ip[3])

		three := data[ip[0]]
		if three == nil {
			three = make([][][]bool, 256)
			data[ip[0]] = three
		}

		two := three[ip[1]]
		if two == nil {
			two = make([][]bool, 256)
			three[ip[1]] = two
		}

		one := two[ip[2]]
		if one == nil {
			one = make([]bool, 256)
			two[ip[2]] = one
		}
		one[ip[3]] = true
		lineNum++
	}
	fmt.Println("rows ", lineNum)
	countResult(data)
}

func countResult(data [][][][]bool) {
	counter := 0
	for i0, three := range data {
		for i1, two := range three {
			for i2, one := range two {
				for i3, val := range one {
					if val {
						fmt.Println(i0, i1, i2, i3, val)
						counter++
					}
				}
			}
		}
	}
	fmt.Println("Unique ips: ", counter)

}

// arr := strings.SplitN(line, ".", 4)
// val, err := strconv.Atoi(arr[0])
// if err != nil {
// 	fmt.Println("error in line ", lineNum)
// 	panic(err)
// }
// fmt.Printf("first: %v, ", val)
// fmt.Println("Old", arr[0], arr[1], arr[2], arr[3])
