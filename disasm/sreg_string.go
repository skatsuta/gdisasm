// generated by stringer -type=Sreg; DO NOT EDIT

package disasm

import "fmt"

const _Sreg_name = "escsssds"

var _Sreg_index = [...]uint8{0, 2, 4, 6, 8}

func (i Sreg) String() string {
	if i < 0 || i+1 >= Sreg(len(_Sreg_index)) {
		return fmt.Sprintf("Sreg(%d)", i)
	}
	return _Sreg_name[_Sreg_index[i]:_Sreg_index[i+1]]
}