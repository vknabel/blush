package op

const (
	_ Opcode = iota
	Const
	ConstTrue
	ConstFalse
	Pop

	Array
	Dict

	// does not consume, just assert top value's type
	AssertType

	Jump
	JumpTrue
	JumpFalse

	Negate
	Invert

	Add
	Sub
	Mul
	Div
	Mod

	Equal
	NotEqual
	GreaterThan
	GreaterThanOrEqual
	LessThan
	LessThanOrEqual

	Call
	Return

	// Serves as instruction to optionally pause on breakpoints.
	// Will not be compiled for non debugging sessions.
	Debug
)

var definitions = map[Opcode]*Definition{
	Const:      {"const", []int{2}}, // const id
	ConstTrue:  {"consttrue", []int{}},
	ConstFalse: {"constfalse", []int{}},
	Pop:        {"pop", []int{}},

	Array: {"array", []int{}},
	Dict:  {"dict", []int{}},

	AssertType: {"asserttype", []int{2}}, // type id

	Jump:      {"jump", []int{2}},      // address
	JumpTrue:  {"jumptrue", []int{2}},  // address
	JumpFalse: {"jumpfalse", []int{2}}, // address

	Negate: {"negate", []int{}},
	Invert: {"invert", []int{}},

	Add: {"add", []int{}},
	Sub: {"sub", []int{}},
	Mul: {"mul", []int{}},
	Div: {"div", []int{}},

	Equal:              {"eq", []int{}},
	NotEqual:           {"neq", []int{}},
	GreaterThan:        {"gt", []int{}},
	GreaterThanOrEqual: {"gte", []int{}},
	LessThan:           {"lt", []int{}},
	LessThanOrEqual:    {"lte", []int{}},

	Call:   {"call", []int{2}},   // arg count
	Return: {"return", []int{2}}, // return value

	Debug: {"debug", []int{}},
}
