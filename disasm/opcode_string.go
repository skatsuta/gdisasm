// generated by stringer -type=Opcode; DO NOT EDIT

package disasm

import "fmt"

const _Opcode_name = "inc"

var _Opcode_index = [...]uint8{0, 3}

func (i Opcode) String() string {
	i -= 1
	if i < 0 || i+1 >= Opcode(len(_Opcode_index)) {
		return fmt.Sprintf("Opcode(%d)", i+1)
	}
	return _Opcode_name[_Opcode_index[i]:_Opcode_index[i+1]]
}