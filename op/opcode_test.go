package op_test

import (
	"testing"

	. "github.com/vknabel/blush/op"
)

func TestMake(t *testing.T) {
	tests := []struct {
		op       Opcode
		operands []int
		want     []byte
	}{
		{Const, []int{65534}, []byte{byte(Const), 255, 254}},
		{GetGlobal, []int{65534}, []byte{byte(GetGlobal), 255, 254}},
		{SetGlobal, []int{65534}, []byte{byte(SetGlobal), 255, 254}},
	}
	for _, tt := range tests {
		instruction := Make(tt.op, tt.operands...)

		if len(instruction) != len(tt.want) {
			t.Errorf("instruction has wrong length. want=%d, got=%d",
				len(tt.want), len(instruction))
		}

		for i, b := range tt.want {
			if instruction[i] != tt.want[i] {
				t.Errorf("wrong byte at pos %d. want=%d, got=%d",
					i, b, instruction[i])
			}
		}
	}
}
