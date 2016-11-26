package datasplit

import (
	"strconv"
	"strings"
)

type buffer []string

// sort buffer
func (b buffer) Len() int {
	return len(b)
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
	oneA := strings.Split(one, "\t")
	twoA := strings.Split(two, "\t")

	if cmp := compareValues(oneA[0], twoA[0]); cmp != 0 {
		return cmp
	}

	if lenA, lenB := len(oneA), len(twoA); lenA > 0 && lenB > 0 {
		return compareByFirstOrNextValue(oneA[1], twoA[1])
	} else if lenA == 0 && lenB == 0 {
		return 0
	} else if lenA > 0 {
		return -1
	} else { // lenB > 0
		return 1
	}
}

func lessByFirstOrNextValue(one, two string) bool {
	return compareByFirstOrNextValue(one, two) < 0
}

func (b buffer) Less(i, j int) bool {
	return compareByFirstOrNextValue(b[i], b[j]) < 0
}

func (b buffer) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}