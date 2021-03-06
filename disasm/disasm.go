package disasm

import (
	"bufio"
	"fmt"
	"io"
)

// maxLenFolInstCod is the maximum length of bytes of an insruction code
// that follows the opcode.
const maxLenFolInstCod = 3

// numBytesPeeked is the number of bytes that are peeked to be interpreted.
const numBytesPeeked = 2

var (
	// 8-bit registers
	reg8 = [...]string{"al", "cl", "dl", "bl", "ah", "ch", "dh", "bh"}
	// 16-bit registers
	reg16 = [...]string{"ax", "cx", "dx", "bx", "sp", "bp", "si", "di"}
	// segment registers
	sreg = [...]string{"es", "cs", "ss", "ds"}
	// effective addresses
	regm = [...]string{"bx+si", "bx+di", "bp+si", "bp+di", "si", "di", "bp", "bx"}
)

/*
type command struct {
	c    byte
	bs   []byte
	mnem Mnemonic
	l    int
	d    byte
	w    byte
	reg  byte
}
*/

// Disasm is a disassembler.
type Disasm struct {
	rdr    *bufio.Reader
	wtr    io.Writer
	offset int // offset
	cmd    *command
}

// New returns a new Disasm.
func New(r *bufio.Reader, w io.Writer) *Disasm {
	return &Disasm{
		rdr:    r,
		wtr:    w,
		offset: 0,
		cmd: &command{
			bs: make([]byte, 1, maxLenFolInstCod),
		},
	}
}

// modrm interprets [mod *** r/m] byte immediately following the opcode.
func modrm(bs []byte) (string, error) {
	if len(bs) < 1 || len(bs) > maxLenFolInstCod {
		return "", fmt.Errorf("the length of %X is invalid", bs)
	}

	b := bs[0]

	mod := b >> 6 // [00]000000: upper two bits
	rm := b & 0x7 // 00000[000]: lower three bits

	switch mod {
	case 0x0: // mod = 00
		if rm == 0x6 { // rm = 110 ==> b = 00***110
			if len(bs) != maxLenFolInstCod {
				return "", modrmErr(rm, bs, maxLenFolInstCod)
			}
			s := fmt.Sprintf("[0x%02x%02x]", bs[2], bs[1])
			return s, nil
		}
		// the length of bs following 00****** (except 00***110) should be 1
		if len(bs) != 1 {
			return "", modrmErr(rm, bs, 1)
		}
		return fmt.Sprintf("[%v]", regm[rm]), nil
	case 0x1: // mod = 01
		if len(bs) != maxLenFolInstCod-1 {
			return "", modrmErr(rm, bs, maxLenFolInstCod-1)
		}
		s := fmt.Sprintf("[%v%+#x]", regm[rm], int8(bs[1]))
		return s, nil
	case 0x2: // mod = 10
		if len(bs) != maxLenFolInstCod {
			return "", modrmErr(rm, bs, maxLenFolInstCod)
		}
		// little endian
		disp := (int16(bs[2]) << 8) | int16(bs[1])
		s := fmt.Sprintf("[%v%+#x]", regm[rm], disp)
		return s, nil
	case 0x3: // mod = 11
		return reg16[rm], nil
	default:
		return "", fmt.Errorf("either mod = %v or r/m = %v is invalid", mod, rm)
	}
}

// modrmErr returns an error that may occur in modrm() function.
func modrmErr(rm byte, bs []byte, l int) error {
	return fmt.Errorf("r/m is %#x but %X does not have length %v", rm, bs, l)
}

// cmdStr returns an disassembled code.
func cmdStr(off int, bs []byte, mnem Mnemonic, w, opr1, opr2 string) string {
	return fmt.Sprintf("%08X  %X   %s %s%s%s", off, bs, mnem.String(), w, opr1, opr2)
}

// Parse parses a set of opcode and operand to an assembly operation.
func (d *Disasm) Parse() (string, error) {
	bs, err := d.rdr.Peek(numBytesPeeked)
	if err == io.EOF {
		return "", err
	}

	d.cmd.bs = bs

	//return d.parse(d.cmd.bs)
	return "", nil
}

func (d *Disasm) parse(b byte) (string, error) {
	switch {
	case b>>1 == 0x7F: // 1111111w
		d.offset += 4
		d.cmd.bs = d.cmd.bs[:3]
		if _, e := d.rdr.Read(d.cmd.bs); e != nil {
			return "", fmt.Errorf("Read() failed: %v", e)
		}
		opr, err := modrm(d.cmd.bs)
		if err != nil {
			return "", fmt.Errorf("modrm(%v) failed: %v", d.cmd.bs, err)
		}
		return cmdStr(d.offset, d.cmd.bs, inc, "word ", opr, ""), nil
	case b>>3 == 0x8: // 01000reg
		d.offset++
		reg := b & 0x7
		return cmdStr(d.offset, nil, inc, "", reg16[reg], ""), nil
	}
	d.offset++
	return "", nil
}

/*
type comProp struct {
	op  Mnemonic
	l   int
	d   byte
	w   byte
	reg byte
}

func (c *command) parseOpcode(bs []byte) error {
	c.init()

	if len(bs) != 2 {
		return errors.New("parseOp: the length of argument must be 2")
	}

	b := bs[0]

	switch {
	// add
	case b>>2 == 0x0:
		c.mnem = add
		c.l = 2
		c.d = getd(b)
		c.w = getw(b)
	case b>>1 == 0x2:
		c.mnem = add
		c.w = getw(b)
		c.l = int(c.w + 1)
	// push
	case b&0xE7 == 0x6:
		c.mnem = push
		c.l = 1
		c.reg = b >> 3 & 0x3
	// pop
	case b&0xE7 == 0x7:
		c.mnem = pop
		c.l = 1
		c.reg = b >> 3 & 0x3
	// or
	case b>>2 == 0x2:
		c.mnem = or
		c.l = 2
		c.d = getd(b)
		c.w = getw(b)
	case b>>1 == 0x6:
		c.mnem = or
		c.l = int(c.w + 1)
		c.w = getw(b)
	// adc
	case b>>2 == 0x4:
		c.mnem = adc
		c.l = 2
		c.d = getd(b)
		c.w = getw(b)
	// sbb
	case b>>2 == 0x6:
		c.mnem = sbb
		c.l = 2
		c.d = getd(b)
		c.w = getd(b)
	case b>>1 == 0x7:
		c.mnem = sbb
		c.l = int(c.w + 1)
		c.w = getw(b)
	// sub
	case b>>2 == 0xA:
		c.mnem = sub
		c.l = 2
		c.d = getd(b)
		c.w = getd(b)

	}
	return nil
}

func getd(b byte) byte {
	return (b >> 1) & 0x1
}

func getw(b byte) byte {
	return b & 0x1
}

func (c *command) init() {
	c.mnem = 0
	c.l = 0
	c.d = 0
	c.w = 0
	c.reg = 0
}
*/
