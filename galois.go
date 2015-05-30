package galois

import "errors"

// MaxGF is the maximum GF(x)
// default 24, up to 31 if you have petabytes of ram :p
const MaxGF = 24

var primPoly = []uint32 {
				0,
	/*  1 */    1,
	/*  2 */    07,
	/*  3 */    013,
	/*  4 */    023,
	/*  5 */    045,
	/*  6 */    0103,
	/*  7 */    0211,
	/*  8 */    0435,
	/*  9 */    01021,
	/* 10 */    02011,
	/* 11 */    04005,
	/* 12 */    010123,
	/* 13 */    020033,
	/* 14 */    042103,
	/* 15 */    0100003,
	/* 16 */    0210013,
	/* 17 */    0400011,
	/* 18 */    01000201,
	/* 19 */    02000047,
	/* 20 */    04000011,
	/* 21 */    010000005,
	/* 22 */    020000003,
	/* 23 */    040000041,
	/* 24 */    0100000207,
	/* 25 */    0200000011,
	/* 26 */    0400000107,
	/* 27 */    01000000047,
	/* 28 */    02000000011,
	/* 29 */    04000000005,
	/* 30 */    010040000007,
	/* 31 */    020000000011,
	/* 32 */    /* 040020000007, overflow */
}


// GfPoly is Polynomial struct
type GfPoly struct {
	base uint8
	NW uint32
	gflog []uint32
	gfilog []uint32
}

var gfPolyInstance = make([]*GfPoly, MaxGF + 1)

func newGF(base uint8) (*GfPoly, error) {

	var b, log uint32
	var poly *GfPoly

	poly = new(GfPoly)
	if base < 2 || base > MaxGF {
		return nil, errors.New("Prim polynomial out of range")
	}
	poly.base = base
	poly.NW = 1 << base
	poly.gflog = make([]uint32, poly.NW)
	poly.gfilog = make([]uint32, poly.NW)
	b = 1

	for log = 0; log < poly.NW; log++ {
		poly.gflog[b] = log
		poly.gfilog[log] = b
		b = b << 1
		if b & poly.NW != 0 {
			b = b ^ primPoly[base]
		}
	}
	return poly, nil
}

// GF is a singleton getter for new GfPoly struct
func GF(base uint8) (*GfPoly, error) {
	var error error

	if base < 2 || base > MaxGF {
		return nil, errors.New("Prim polynomial out of range")
	}
	if gfPolyInstance[base] == nil {
		gfPolyInstance[base], error = newGF(base)
	}
	return gfPolyInstance[base], error
}

// Mul is GfPoly struct method for multiplication
func (table *GfPoly) Mul(a, b uint32) (uint32, error) {
	if a == 0 || b == 0 {
		return 0, nil
	}
	if a >= table.NW || b >= table.NW {
		return 0, errors.New("mul: polynomial out of range")
	}
	sumLog := table.gflog[a] + table.gflog[b]
	if sumLog >= table.NW - 1 {
		sumLog -= table.NW - 1
	}
	return table.gfilog[sumLog], nil
}

// Div is GfPoly struct method for division
func (table *GfPoly) Div(a, b uint32) (uint32, error) {
	var diffLog int64
	if b == 0 {
		return 0, errors.New("div: division by 0 :/")
	}
	if a == 0 {
		return 0, nil
	}
	if a >= table.NW || b >= table.NW {
		return 0, errors.New("div: polynomial out of range")
	}
	diffLog = int64(table.gflog[a]) - int64(table.gflog[b])
	if diffLog < 0 {
		diffLog += int64(table.NW) - 1
	}
	return table.gfilog[diffLog], nil
}

// Expon is GfPoly struct method for exponential
func (table *GfPoly) Expon(a, e uint32) (uint32, error) {

	var err error
	var i uint32

	b := a
	if e == 0 {
		return 1, nil
	}
	if e == 1 {
		return a, nil
	}
	for i = 1; i < e; i++ {
		b, err = table.Mul(b, a)
		if err != nil {
			return 0, err
		}
	}
	return b, nil
}
