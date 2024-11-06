package main

type LayerIf interface {
	add(input uint32, nestIndex int)
	count() uint32
}

type Layer5 [32]LayerIf
type LayerLast [64]uint64

var nestOffsets = [...]byte{5, 5, 5, 5, 12} // sum must be 32

func (i *Layer5) add(input uint32, nestIndex int) {
	nestIndex++
	index := input >> 27 // 32 - N
	inner := i[index]
	if inner == nil {
		switch nestOffsets[nestIndex] {
		case 5:
			inner = &Layer5{}
		default:
			inner = &LayerLast{}
		}
		i[index] = inner
	}
	inner.add(input<<5, nestIndex)
}

func (i *LayerLast) add(input uint32, nestIndex int) {
	arrIndex := input >> 26
	longindex := (input << 6) >> 26
	i[arrIndex] |= uint64(1) << longindex
}

func (i Layer5) count() uint32 {
	var sum uint32 = 0
	for _, inner := range i {
		if inner != nil {
			sum += inner.count()
		}
	}
	return sum
}

func (i LayerLast) count() uint32 {
	var sum uint32 = 0
	for _, arrValue := range i {
		for ii := uint64(1); ii > 0; ii <<= 1 {
			if arrValue&ii > 0 {
				sum++
			}
		}
	}
	return sum
}
