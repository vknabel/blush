package op

const (
	_ Opcode = iota
	Const
	Pop

	Jump
	JumpTrue
	JumpFalse

	Add
	Sub
	Mul
	Div

	// Serves as instruction to optionally pause on breakpoints.
	// Will not be compiled for non debugging sessions.
	Debug
)

var definitions = map[Opcode]*Definition{
	Const: {"const", []int{2}},
	Pop:   {"pop", []int{}},

	Jump:      {"jump", []int{2}},
	JumpTrue:  {"jumptrue", []int{2}},
	JumpFalse: {"jumpfalse", []int{2}},

	Add: {"add", []int{}},
	Sub: {"sub", []int{}},
	Mul: {"mul", []int{}},
	Div: {"div", []int{}},

	Debug: {"debug", []int{}},
}
