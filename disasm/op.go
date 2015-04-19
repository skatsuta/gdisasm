//go:generate stringer -type=Opcode

package disasm

// Opcode is an operation code.
type Opcode int

// Opcodes.
const (
	_   Opcode = iota
	inc        // increment
)