package op

const (
	_ Opcode = iota
	CONST

	ADD
)

var definitions = map[Opcode]*Definition{
	CONST: {"CONST", []int{2}},
	ADD:   {"ADD", []int{}},
}
