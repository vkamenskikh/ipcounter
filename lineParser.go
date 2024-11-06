package main

import "fmt"

func parse4(line string, ip *[4]uint8) error { // does not create new objects for garbage collection
	pos := 0
	var number uint8 = 0
	for _, char := range []byte(line) {
		if char == '.' {
			ip[pos] = number
			number = 0
			pos++
			continue
		}

		numChar := char - '0'
		if numChar > 9 {
			if char == '\n' {
				break
			}
			if char == ' ' || char == '\r' {
				continue
			}
			return fmt.Errorf("unsupported character: %v", char)
		}
		number = number*10 + numChar
	}
	if pos != 3 {
		return fmt.Errorf("invalid ip address: %v", line)
	}
	ip[pos] = number
	return nil
}

func parseToInt(line string) (uint32, error) { // does not create new objects for garbage collection
	var result uint32 = 0
	var number uint8 = 0
	pos := 0

	for _, char := range []byte(line) {
		if char == '.' {
			result = (result + uint32(number)) << 8
			number = 0
			pos++
			continue
		}

		numChar := char - '0'
		if numChar > 9 {
			if char == '\n' {
				break
			}
			if char == ' ' || char == '\r' {
				continue
			}
			return 0, fmt.Errorf("unsupported character: %v", char)
		}
		number = number*10 + numChar
	}
	if pos != 3 {
		return result, fmt.Errorf("invalid ip address: %v", line)
	}
	result += uint32(number)
	return result, nil
}
