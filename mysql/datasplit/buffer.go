package datasplit

import (
	"errors"
	"strconv"
	"strings"
)

type buffer []string

// sort buffer
func (b buffer) Len() int {
	return len(b)
}

func cleanStartSpaces(line string) string {
	for line != "" && (line[0] == ' ' || line[0] == '\t') {
		line = line[1:]
	}
	return line
}

func headAndTail(line string, pos int) (string, string) {
	parts := strings.SplitN(line[pos:], ",", 2)
	if len(parts) < 2 {
		return parts[0], ""
	} else {
		return parts[0], cleanStartSpaces(parts[1])
	}
}

func firstValueAndOther(line string) (error, string, string) {
	if line == "" {
		return nil, "", ""
	}

	if line[0] == '\'' || line[0] == '"' {
		last := 1
		for {
			pos := strings.IndexByte(line[last:], line[0])
			if pos < 0 {
				return errors.New("Can't find ending quote in: " + line), "", ""
			} else if pos == 1 {
				_, tail := headAndTail(line, pos + last + 1)
				return nil, "", tail
			} else if line[last + pos - 1] == '\\' {
				last = last + pos + 1
				continue
			} else {
				value := line[1:last + pos]
				_, tail := headAndTail(line, pos + last + 1)
				return nil, value, tail
			}
		}
	} else {
		head, tail := headAndTail(line, 0)
		return nil, head, tail
	}
}

func compareValues(one, two string) int {
	if _1, _e1 := strconv.ParseFloat(one, 64); _e1 == nil {
		if _2, _e2 := strconv.ParseFloat(two, 64); _e2 == nil {
			if d := _1 - _2; d < 0 {
				return -1
			} else if d > 0 {
				return 1
			}
			return 0
		}
	}

	if one < two {
		return -1
	} else if one > two {
		return 1
	}
	return 0
}

func compareByFirstOrNextValue(one, two string) int {
	errA, oneA, tailA := firstValueAndOther(one)
	if errA != nil {
		return -1
	}
	errB, oneB, tailB := firstValueAndOther(two)
	if errB != nil {
		return 1
	}

	if cmp := compareValues(oneA, oneB); cmp != 0 {
		return cmp
	}

	if tailA == "" || tailB == "" {
		return strings.Compare(tailA, tailB)
	} else {
		return compareByFirstOrNextValue(tailA, tailB)
	}

}

func lessByFirstOrNextValue(one, two string) bool {
	return compareByFirstOrNextValue(one, two) < 0
}

// Implementation of Sort interface
func (b buffer) Less(i, j int) bool {
	return compareByFirstOrNextValue(b[i], b[j]) < 0
}

func (b buffer) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}