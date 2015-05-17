package disasm

import "errors"

type command struct {
	bs   []byte
	mnem Mnemonic
	l    int
	d    byte
	w    byte
	reg  Reg
}

func (c *command) parseOpcode(bs []byte) error {
	c.init()

	if len(bs) != 2 {
		return errors.New("parseOpcode: the length of argument must be 2")
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
		c.reg = Sreg(b >> 3 & 0x3)

	// pop
	case b&0xE7 == 0x7:
		c.mnem = pop
		c.l = 1
		c.reg = Sreg(b >> 3 & 0x3)

	// or
	case b>>2 == 0x2:
		c.mnem = or
		c.l = 2
		c.d = getd(b)
		c.w = getw(b)
	case b>>1 == 0x6:
		c.mnem = or
		c.w = getw(b)
		c.l = int(c.w + 1)

	// adc
	case b>>2 == 0x4:
		c.mnem = adc
		c.l = 2
		c.d = getd(b)
		c.w = getw(b)
	case b>>1 == 0xA:
		c.mnem = adc
		c.w = getw(b)
		c.l = int(c.w + 1)

	// sbb
	case b>>2 == 0x6:
		c.mnem = sbb
		c.l = 2
		c.d = getd(b)
		c.w = getw(b)
	case b>>1 == 0xE:
		c.mnem = sbb
		c.w = getw(b)
		c.l = int(c.w + 1)

	// sub
	case b>>2 == 0xA:
		c.mnem = sub
		c.l = 2
		c.d = getd(b)
		c.w = getw(b)

	// and
	case b>>2 == 0x8:
		c.mnem = and
		c.l = 2
		c.d = getd(b)
		c.w = getw(b)
	case b>>1 == 0x12:
		c.mnem = and
		c.w = getw(b)
		c.l = int(c.w + 1)
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
	c.bs = nil
	c.mnem = 0
	c.l = 0
	c.d = 0
	c.w = 0
	c.reg = nil
}