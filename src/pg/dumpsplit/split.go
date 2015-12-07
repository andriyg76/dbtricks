package dumpsplit
import (
	"regexp"
	"orders"
	"pg/datasplit"
	"writer"
)

/**
 * Created by andriy on 04/12/15.
 */

type Splitter interface {
	Flush() error
	Close()
	HandleLine(line string) error
	Error() error
}

var copy_re, data_comment_re, constraint_comment_re *regexp.Regexp;

func init() {
	copy_re = regexp.MustCompile(`^COPY .*? \((.*?)\) FROM stdin;$`)
	data_comment_re = regexp.MustCompile(
		"^-- Data for Name: (?P<table>.*?); Type: TABLE DATA; Schema: (?P<schema>.*?);$")
	constraint_comment_re = regexp.MustCompile("^-- Name: .*; Type: (.*CONSTRAINT|INDEX); Schema: ")
}

func match_to_copy(line string) bool {
	return copy_re.MatchString(line)
}

func match_to_data_comment(line string) (match bool, table string, schema string) {
	if !data_comment_re.MatchString(line) {
		match = false
		table = ""
		schema = ""
		return
	}

	match = true
	match_values := data_comment_re.FindStringSubmatch(line)
	table = match_values[1]
	schema = match_values[2]
	return
}

func match_to_constraint_comment(line string) bool {
	return constraint_comment_re.MatchString(line)
}

const eot_line = "."

type splitter struct {
	counter        int
	previous_table string
	dumper         writer.Writer
	table_name     string
	epilogue       bool
	data_handler   datasplit.DataSplitter
	chunk_size     int
	orders         orders.Orders
	err            error
}

func NewSplitter(orders orders.Orders, chunk_size int) (Splitter, error) {
	dumper, err := writer.NewWriter("0000_prologue.sql")
	if err != nil {
		return nil, err
	}
	return &splitter{
		counter: 0,
		dumper: dumper,
		previous_table: "",
		table_name: "",
		orders: orders,
		chunk_size: chunk_size,
	}, nil
}

func (i *splitter) HandleLine(line string) error {
	if i.data_handler != nil {
		if line == eot_line {
			i.data_handler.FlushData(i.dumper)
			i.data_handler = nil
			i.previous_table = i.table_name
		} else {
			i.data_handler.AddLine(line)
		}
	} else {
		if i.epilogue || line == "" || line == "--" {
			i.dumper.AddLines(line)
		} else if match_to_constraint_comment(line) {
			backup := append(i.dumper.PopLastLine(), line)
			i.dumper.ResetOutput("zzzz_epilogue.sql")
			i.dumper.AddLines(backup...)
			i.epilogue = true
		} else {
			match, table, schema := match_to_data_comment(line)
			if match {
				table_name := schema + "." + table
				table_order := i.orders.GetTableOrder(table_name)

				backup := append(i.dumper.PopLastLine(), line)
				i.dumper.ResetOutput(_36_base_int(table_order) + "_" + table_name + ".sql")
				i.dumper.AddLines(backup...)
			} else if i.table_name != "" && match_to_copy(line) {
				i.data_handler = datasplit.NewDataSplitter(i.chunk_size, line, i.table_name, i.counter)
			} else {
				i.dumper.AddLines(line)
			}
		}
	}
	return nil
}

var _conversions = map[int]int{
	0: '0',
	1: '1',
	2: '2',
	3: '3',
	4: '4',
	5: '5',
	6: '6',
	7: '7',
	8: '8',
	9: '9',
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

func _36_base_int(value int) string {
	return int_in_base(value, 36, 4)
}

func int_in_base(value, base int, min_with int) string {
	if base < 2 || base > 36 {
		panic("int_in_base: Support positional systems from 2 to 36 only")
	}

	new_num_string := ""
	current := value
	for {
		remainder := current % base
		if remainder < 0 {
			remainder = -remainder
		}
		new_num_string = string(_conversions[remainder]) + new_num_string
		current = current / base
		if current == 0 {
			break;
		}
	}

	for len(new_num_string) < min_with {
		new_num_string = "0" + new_num_string
	}
	if value < 0 {
		new_num_string = "-" + new_num_string
	}
	return new_num_string
}

func (i *splitter) Flush() error {
	return nil
}

func (i *splitter) Close() {

}

func (i *splitter) Error() error {
	return i.err
}
