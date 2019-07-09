package utils

var _conversions = map[int]int{
	0:  '0',
	1:  '1',
	2:  '2',
	3:  '3',
	4:  '4',
	5:  '5',
	6:  '6',
	7:  '7',
	8:  '8',
	9:  '9',
	10: 'a',
	11: 'b',
	12: 'c',
	13: 'd',
	14: 'e',
	15: 'f',
	16: 'g',
	17: 'h',
	18: 'i',
	19: 'j',
	20: 'k',
	21: 'l',
	22: 'm',
	23: 'n',
	24: 'o',
	25: 'p',
	26: 'q',
	27: 'r',
	28: 's',
	29: 't',
	30: 'u',
	31: 'v',
	32: 'w',
	33: 'x',
	34: 'y',
	35: 'z',
}

func IntInBase(value, base int, minWith int) string {
	if base < 2 || base > 36 {
		panic("int_in_base: Support positional systems from 2 to 36 only")
	}

	newNumString := ""
	current := value
	for {
		remainder := current % base
		if remainder < 0 {
			remainder = -remainder
		}
		newNumString = string(_conversions[remainder]) + newNumString
		current = current / base
		if current == 0 {
			break
		}
	}

	for len(newNumString) < minWith {
		newNumString = "0" + newNumString
	}
	if value < 0 {
		newNumString = "-" + newNumString
	}
	return newNumString
}
