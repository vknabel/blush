package op

const (
	_ Opcode = iota
	Const
	ConstTrue
	ConstFalse
	Pop

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

	// Serves as instruction to optionally pause on breakpoints.
	// Will not be compiled for non debugging sessions.
	Debug
)

var definitions = map[Opcode]*Definition{
	Const:      {"const", []int{2}},
	ConstTrue:  {"consttrue", []int{}},
	ConstFalse: {"constfalse", []int{}},
	Pop:        {"pop", []int{}},

	AssertType: {"asserttype", []int{2}},

	Jump:      {"jump", []int{2}},
	JumpTrue:  {"jumptrue", []int{2}},
	JumpFalse: {"jumpfalse", []int{2}},

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

	Debug: {"debug", []int{}},
}
