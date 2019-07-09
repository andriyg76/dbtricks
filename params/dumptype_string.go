// Code generated by "stringer -type DumpType"; DO NOT EDIT.

package params

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[None - -1]
	_ = x[Mysql-1]
	_ = x[Pgsql-2]
}

const (
	_Dumptype_name_0 = "None"
	_Dumptype_name_1 = "MysqlPgsql"
)

var (
	_Dumptype_index_1 = [...]uint8{0, 5, 10}
)

func (i DumpType) String() string {
	switch {
	case i == -1:
		return _Dumptype_name_0
	case 1 <= i && i <= 2:
		i -= 1
		return _Dumptype_name_1[_Dumptype_index_1[i]:_Dumptype_index_1[i+1]]
	default:
		return "DumpType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}