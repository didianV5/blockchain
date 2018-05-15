package main

import(
	"fmt"
	"unsafe"
)

const Size = 16

const (
	chunk = 64
	init0 = 0x67452301
	init1 = 0xEFCDAB89
	init2 = 0x98BADCFE
	init3 = 0x10325476
)

// digest represents the partial evaluation of a checksum.
type digest struct {
	s   [4]uint32
	x   [chunk]byte
	nx  int
	len uint64
}

func (d *digest) Reset() {
	d.s[0] = init0
	d.s[1] = init1
	d.s[2] = init2
	d.s[3] = init3
	d.nx = 0
	d.len = 0
}

func (d *digest) Write(p []byte) (nn int, err error) {
	nn = len(p)
	d.len += uint64(nn)
	if d.nx > 0 {
		n := copy(d.x[d.nx:], p)
		d.nx += n
		if d.nx == chunk {
			blockGeneric(d, d.x[:])
			d.nx = 0
		}
		p = p[n:]
	}
	if len(p) >= chunk {
		n := len(p) &^ (chunk - 1)
		blockGeneric(d, p[:n])
		p = p[n:]
	}
	if len(p) > 0 {
		d.nx = copy(d.x[:], p)
	}

	fmt.Printf("%v\n",d)

	return
}

func (d *digest) checkSum() [Size]byte {
	// Padding. Add a 1 bit and 0 bits until 56 bytes mod 64.
	len := d.len
	var tmp [64]byte
	tmp[0] = 0x80
	if len%64 < 56 {
		d.Write(tmp[0 : 56-len%64])
		fmt.Printf("%v",d)
	} else {
		d.Write(tmp[0 : 64+56-len%64])
	}

	// Length in bits.
	len <<= 3
	for i := uint(0); i < 8; i++ {
		tmp[i] = byte(len >> (8 * i))
	}
	d.Write(tmp[0:8])

	if d.nx != 0 {
		panic("d.nx != 0")
	}

	var digest [Size]byte
	for i, s := range d.s {
		digest[i*4] = byte(s)
		digest[i*4+1] = byte(s >> 8)
		digest[i*4+2] = byte(s >> 16)
		digest[i*4+3] = byte(s >> 24)
	}

	return digest
}

func blockGeneric(dig *digest, p []byte) {
	a := dig.s[0]
	b := dig.s[1]
	c := dig.s[2]
	d := dig.s[3]
	var X *[16]uint32
	for len(p) >= chunk {
		aa, bb, cc, dd := a, b, c, d

		X = (*[16]uint32)(unsafe.Pointer(&p[0]))

		// Round 1.

		a += (((c ^ d) & b) ^ d) + X[0] + 3614090360
		a = a<<7 | a>>(32-7) + b

		d += (((b ^ c) & a) ^ c) + X[1] + 3905402710
		d = d<<12 | d>>(32-12) + a

		c += (((a ^ b) & d) ^ b) + X[2] + 606105819
		c = c<<17 | c>>(32-17) + d

		b += (((d ^ a) & c) ^ a) + X[3] + 3250441966
		b = b<<22 | b>>(32-22) + c

		a += (((c ^ d) & b) ^ d) + X[4] + 4118548399
		a = a<<7 | a>>(32-7) + b

		d += (((b ^ c) & a) ^ c) + X[5] + 1200080426
		d = d<<12 | d>>(32-12) + a

		c += (((a ^ b) & d) ^ b) + X[6] + 2821735955
		c = c<<17 | c>>(32-17) + d

		b += (((d ^ a) & c) ^ a) + X[7] + 4249261313
		b = b<<22 | b>>(32-22) + c

		a += (((c ^ d) & b) ^ d) + X[8] + 1770035416
		a = a<<7 | a>>(32-7) + b

		d += (((b ^ c) & a) ^ c) + X[9] + 2336552879
		d = d<<12 | d>>(32-12) + a

		c += (((a ^ b) & d) ^ b) + X[10] + 4294925233
		c = c<<17 | c>>(32-17) + d

		b += (((d ^ a) & c) ^ a) + X[11] + 2304563134
		b = b<<22 | b>>(32-22) + c

		a += (((c ^ d) & b) ^ d) + X[12] + 1804603682
		a = a<<7 | a>>(32-7) + b

		d += (((b ^ c) & a) ^ c) + X[13] + 4254626195
		d = d<<12 | d>>(32-12) + a

		c += (((a ^ b) & d) ^ b) + X[14] + 2792965006
		c = c<<17 | c>>(32-17) + d

		b += (((d ^ a) & c) ^ a) + X[15] + 1236535329
		b = b<<22 | b>>(32-22) + c

		// Round 2.

		a += (((b ^ c) & d) ^ c) + X[(1+5*0)&15] + 4129170786
		a = a<<5 | a>>(32-5) + b

		d += (((a ^ b) & c) ^ b) + X[(1+5*1)&15] + 3225465664
		d = d<<9 | d>>(32-9) + a

		c += (((d ^ a) & b) ^ a) + X[(1+5*2)&15] + 643717713
		c = c<<14 | c>>(32-14) + d

		b += (((c ^ d) & a) ^ d) + X[(1+5*3)&15] + 3921069994
		b = b<<20 | b>>(32-20) + c

		a += (((b ^ c) & d) ^ c) + X[(1+5*4)&15] + 3593408605
		a = a<<5 | a>>(32-5) + b

		d += (((a ^ b) & c) ^ b) + X[(1+5*5)&15] + 38016083
		d = d<<9 | d>>(32-9) + a

		c += (((d ^ a) & b) ^ a) + X[(1+5*6)&15] + 3634488961
		c = c<<14 | c>>(32-14) + d

		b += (((c ^ d) & a) ^ d) + X[(1+5*7)&15] + 3889429448
		b = b<<20 | b>>(32-20) + c

		a += (((b ^ c) & d) ^ c) + X[(1+5*8)&15] + 568446438
		a = a<<5 | a>>(32-5) + b

		d += (((a ^ b) & c) ^ b) + X[(1+5*9)&15] + 3275163606
		d = d<<9 | d>>(32-9) + a

		c += (((d ^ a) & b) ^ a) + X[(1+5*10)&15] + 4107603335
		c = c<<14 | c>>(32-14) + d

		b += (((c ^ d) & a) ^ d) + X[(1+5*11)&15] + 1163531501
		b = b<<20 | b>>(32-20) + c

		a += (((b ^ c) & d) ^ c) + X[(1+5*12)&15] + 2850285829
		a = a<<5 | a>>(32-5) + b

		d += (((a ^ b) & c) ^ b) + X[(1+5*13)&15] + 4243563512
		d = d<<9 | d>>(32-9) + a

		c += (((d ^ a) & b) ^ a) + X[(1+5*14)&15] + 1735328473
		c = c<<14 | c>>(32-14) + d

		b += (((c ^ d) & a) ^ d) + X[(1+5*15)&15] + 2368359562
		b = b<<20 | b>>(32-20) + c

		// Round 3.

		a += (b ^ c ^ d) + X[(5+3*0)&15] + 4294588738
		a = a<<4 | a>>(32-4) + b

		d += (a ^ b ^ c) + X[(5+3*1)&15] + 2272392833
		d = d<<11 | d>>(32-11) + a

		c += (d ^ a ^ b) + X[(5+3*2)&15] + 1839030562
		c = c<<16 | c>>(32-16) + d

		b += (c ^ d ^ a) + X[(5+3*3)&15] + 4259657740
		b = b<<23 | b>>(32-23) + c

		a += (b ^ c ^ d) + X[(5+3*4)&15] + 2763975236
		a = a<<4 | a>>(32-4) + b

		d += (a ^ b ^ c) + X[(5+3*5)&15] + 1272893353
		d = d<<11 | d>>(32-11) + a

		c += (d ^ a ^ b) + X[(5+3*6)&15] + 4139469664
		c = c<<16 | c>>(32-16) + d

		b += (c ^ d ^ a) + X[(5+3*7)&15] + 3200236656
		b = b<<23 | b>>(32-23) + c

		a += (b ^ c ^ d) + X[(5+3*8)&15] + 681279174
		a = a<<4 | a>>(32-4) + b

		d += (a ^ b ^ c) + X[(5+3*9)&15] + 3936430074
		d = d<<11 | d>>(32-11) + a

		c += (d ^ a ^ b) + X[(5+3*10)&15] + 3572445317
		c = c<<16 | c>>(32-16) + d

		b += (c ^ d ^ a) + X[(5+3*11)&15] + 76029189
		b = b<<23 | b>>(32-23) + c

		a += (b ^ c ^ d) + X[(5+3*12)&15] + 3654602809
		a = a<<4 | a>>(32-4) + b

		d += (a ^ b ^ c) + X[(5+3*13)&15] + 3873151461
		d = d<<11 | d>>(32-11) + a

		c += (d ^ a ^ b) + X[(5+3*14)&15] + 530742520
		c = c<<16 | c>>(32-16) + d

		b += (c ^ d ^ a) + X[(5+3*15)&15] + 3299628645
		b = b<<23 | b>>(32-23) + c

		// Round 4.

		a += (c ^ (b | ^d)) + X[(7*0)&15] + 4096336452
		a = a<<6 | a>>(32-6) + b

		d += (b ^ (a | ^c)) + X[(7*1)&15] + 1126891415
		d = d<<10 | d>>(32-10) + a

		c += (a ^ (d | ^b)) + X[(7*2)&15] + 2878612391
		c = c<<15 | c>>(32-15) + d

		b += (d ^ (c | ^a)) + X[(7*3)&15] + 4237533241
		b = b<<21 | b>>(32-21) + c

		a += (c ^ (b | ^d)) + X[(7*4)&15] + 1700485571
		a = a<<6 | a>>(32-6) + b

		d += (b ^ (a | ^c)) + X[(7*5)&15] + 2399980690
		d = d<<10 | d>>(32-10) + a

		c += (a ^ (d | ^b)) + X[(7*6)&15] + 4293915773
		c = c<<15 | c>>(32-15) + d

		b += (d ^ (c | ^a)) + X[(7*7)&15] + 2240044497
		b = b<<21 | b>>(32-21) + c

		a += (c ^ (b | ^d)) + X[(7*8)&15] + 1873313359
		a = a<<6 | a>>(32-6) + b

		d += (b ^ (a | ^c)) + X[(7*9)&15] + 4264355552
		d = d<<10 | d>>(32-10) + a

		c += (a ^ (d | ^b)) + X[(7*10)&15] + 2734768916
		c = c<<15 | c>>(32-15) + d

		b += (d ^ (c | ^a)) + X[(7*11)&15] + 1309151649
		b = b<<21 | b>>(32-21) + c

		a += (c ^ (b | ^d)) + X[(7*12)&15] + 4149444226
		a = a<<6 | a>>(32-6) + b

		d += (b ^ (a | ^c)) + X[(7*13)&15] + 3174756917
		d = d<<10 | d>>(32-10) + a

		c += (a ^ (d | ^b)) + X[(7*14)&15] + 718787259
		c = c<<15 | c>>(32-15) + d

		b += (d ^ (c | ^a)) + X[(7*15)&15] + 3951481745
		b = b<<21 | b>>(32-21) + c

		a += aa
		b += bb
		c += cc
		d += dd

		p = p[chunk:]
	}

	dig.s[0] = a
	dig.s[1] = b
	dig.s[2] = c
	dig.s[3] = d
}

func main() {
	var d digest
	d.Reset()
	data := []byte("hello world")

	d.Write(data)

	hashnum := d.checkSum()

	fmt.Printf("%x\n",hashnum)

}