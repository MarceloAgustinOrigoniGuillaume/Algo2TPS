package xxh3

import (
	"bytes"
	"encoding/binary"
	"math/bits"
	"unsafe"
)

var key = ptr(&[...]u8{
	0xb8, 0xfe, 0x6c, 0x39, 0x23, 0xa4, 0x4b, 0xbe /* 8   */, 0x7c, 0x01, 0x81, 0x2c, 0xf7, 0x21, 0xad, 0x1c, /* 16  */
	0xde, 0xd4, 0x6d, 0xe9, 0x83, 0x90, 0x97, 0xdb /* 24  */, 0x72, 0x40, 0xa4, 0xa4, 0xb7, 0xb3, 0x67, 0x1f, /* 32  */
	0xcb, 0x79, 0xe6, 0x4e, 0xcc, 0xc0, 0xe5, 0x78 /* 40  */, 0x82, 0x5a, 0xd0, 0x7d, 0xcc, 0xff, 0x72, 0x21, /* 48  */
	0xb8, 0x08, 0x46, 0x74, 0xf7, 0x43, 0x24, 0x8e /* 56  */, 0xe0, 0x35, 0x90, 0xe6, 0x81, 0x3a, 0x26, 0x4c, /* 64  */
	0x3c, 0x28, 0x52, 0xbb, 0x91, 0xc3, 0x00, 0xcb /* 72  */, 0x88, 0xd0, 0x65, 0x8b, 0x1b, 0x53, 0x2e, 0xa3, /* 80  */
	0x71, 0x64, 0x48, 0x97, 0xa2, 0x0d, 0xf9, 0x4e /* 88  */, 0x38, 0x19, 0xef, 0x46, 0xa9, 0xde, 0xac, 0xd8, /* 96  */
	0xa8, 0xfa, 0x76, 0x3f, 0xe3, 0x9c, 0x34, 0x3f /* 104 */, 0xf9, 0xdc, 0xbb, 0xc7, 0xc7, 0x0b, 0x4f, 0x1d, /* 112 */
	0x8a, 0x51, 0xe0, 0x4b, 0xcd, 0xb4, 0x59, 0x31 /* 120 */, 0xc8, 0x9f, 0x7e, 0xc9, 0xd9, 0x78, 0x73, 0x64, /* 128 */
	0xea, 0xc5, 0xac, 0x83, 0x34, 0xd3, 0xeb, 0xc3 /* 136 */, 0xc5, 0x81, 0xa0, 0xff, 0xfa, 0x13, 0x63, 0xeb, /* 144 */
	0x17, 0x0d, 0xdd, 0x51, 0xb7, 0xf0, 0xda, 0x49 /* 152 */, 0xd3, 0x16, 0x55, 0x26, 0x29, 0xd4, 0x68, 0x9e, /* 160 */
	0x2b, 0x16, 0xbe, 0x58, 0x7d, 0x47, 0xa1, 0xfc /* 168 */, 0x8f, 0xf8, 0xb8, 0xd1, 0x7a, 0xd0, 0x31, 0xce, /* 176 */
	0x45, 0xcb, 0x3a, 0x8f, 0x95, 0x16, 0x04, 0x28 /* 184 */, 0xaf, 0xd7, 0xfb, 0xca, 0xbb, 0x4b, 0x40, 0x7e, /* 192 */
})

//utils

type (
	ptr = unsafe.Pointer
	u64 = uint64
	u8  = uint8
	u32 = uint32
	ui  = uintptr
)
type GeneralBytes struct {
	p ptr
	u uint
}

// constants
const (
	prime64_1 = 11400714785074694791
	prime64_2 = 14029467366897019727
	prime64_3 = 1609587929392839161

	key32_000 u32 = 0xbe4ba423
	key32_004 u32 = 0x396cfeb8
	key64_000 u64 = 0xbe4ba423396cfeb8
	key64_008 u64 = 0x1cad21f72c81017c
	key64_016 u64 = 0xdb979083e96dd4de
	key64_024 u64 = 0x1f67b3b7a4a44072
	key64_032 u64 = 0x78e5c0cc4ee679cb
	key64_040 u64 = 0x2172ffcc7dd05a82
	key64_048 u64 = 0x8e2443f7744608b8
	key64_056 u64 = 0x4c263a81e69035e0
	key64_064 u64 = 0xcb00c391bb52283c
	key64_072 u64 = 0xa32e531b8b65d088
	key64_080 u64 = 0x4ef90da297486471
	key64_088 u64 = 0xd8acdea946ef1938
	key64_096 u64 = 0x3f349ce33f76faa8
	key64_104 u64 = 0x1d4f0bc7c7bbdcf9
	key64_112 u64 = 0x3159b4cd4be0518a
	key64_120 u64 = 0x647378d9c97e9fc8

	key64_119 u64 = 0x7378d9c97e9fc831
	key64_127 u64 = 0xebd33483acc5ea64
)

func readU8(p ptr, o ui) uint8 {
	return *(*uint8)(ptr(ui(p) + o))
}

func readU16(p ptr, o ui) uint16 {
	b := (*[2]byte)(ptr(ui(p) + o))
	return uint16(b[0]) | uint16(b[1])<<8
}

func readU32(p ptr, o ui) uint32 {
	b := (*[4]byte)(ptr(ui(p) + o))
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
}

func readU64(p ptr, o ui) uint64 {
	b := (*[8]byte)(ptr(ui(p) + o))
	return uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 |
		uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56
}

func Xxh3Avalanche(x u64) u64 {
	x ^= x >> 37
	x *= 0x165667919e3779f9
	x ^= x >> 32
	return x
}

func xxhAvalancheSmall(x u64) u64 {
	x ^= x >> 33
	x *= prime64_2
	x ^= x >> 29
	x *= prime64_3
	x ^= x >> 32
	return x
}

func mulFold64(x, y u64) u64 {
	hi, lo := bits.Mul64(x, y)
	return hi ^ lo
}

func Rrmxmx(h64 u64, len u64) u64 {
	h64 ^= bits.RotateLeft64(h64, 49) ^ bits.RotateLeft64(h64, 24)
	h64 *= 0x9fb21c651e98df25
	h64 ^= (h64 >> 35) + len
	h64 *= 0x9fb21c651e98df25
	h64 ^= (h64 >> 28)
	return h64
}

// Hash function

func hashAny(s GeneralBytes) (acc u64) {
	p, l := s.p, s.u

	switch {
	case l <= 16:
		switch {
		case l > 8: // 9-16
			inputlo := readU64(p, 0) ^ (key64_024 ^ key64_032)
			inputhi := readU64(p, ui(l)-8) ^ (key64_040 ^ key64_048)
			folded := mulFold64(inputlo, inputhi)
			return Xxh3Avalanche(u64(l) + bits.ReverseBytes64(inputlo) + inputhi + folded)

		case l > 3: // 4-8
			input1 := readU32(p, 0)
			input2 := readU32(p, ui(l)-4)
			input64 := u64(input2) + u64(input1)<<32
			keyed := input64 ^ (key64_008 ^ key64_016)
			return Rrmxmx(keyed, u64(l))

		case l == 3: // 3
			c12 := u64(readU16(p, 0))
			c3 := u64(readU8(p, 2))
			acc = c12<<16 + c3 + 3<<8

		case l > 1: // 2
			c12 := u64(readU16(p, 0))
			acc = c12*(1<<24+1)>>8 + 2<<8

		case l == 1: // 1
			c1 := u64(readU8(p, 0))
			acc = c1*(1<<24+1<<16+1) + 1<<8

		default: // 0
			return 0x2d06800538d394c2 // xxh_avalanche(key64_056 ^ key64_064)
		}

		acc ^= u64(key32_000 ^ key32_004)
		return xxhAvalancheSmall(acc)

	case l <= 128:
		acc = u64(l) * prime64_1

		if l > 32 {
			if l > 64 {
				if l > 96 {
					acc += mulFold64(readU64(p, 6*8)^key64_096, readU64(p, 7*8)^key64_104)
					acc += mulFold64(readU64(p, ui(l)-8*8)^key64_112, readU64(p, ui(l)-7*8)^key64_120)
				} // 96
				acc += mulFold64(readU64(p, 4*8)^key64_064, readU64(p, 5*8)^key64_072)
				acc += mulFold64(readU64(p, ui(l)-6*8)^key64_080, readU64(p, ui(l)-5*8)^key64_088)
			} // 64
			acc += mulFold64(readU64(p, 2*8)^key64_032, readU64(p, 3*8)^key64_040)
			acc += mulFold64(readU64(p, ui(l)-4*8)^key64_048, readU64(p, ui(l)-3*8)^key64_056)
		} // 32
		acc += mulFold64(readU64(p, 0*8)^key64_000, readU64(p, 1*8)^key64_008)
		acc += mulFold64(readU64(p, ui(l)-2*8)^key64_016, readU64(p, ui(l)-1*8)^key64_024)

		return Xxh3Avalanche(acc)

	default: //case l <= 240:
		acc = u64(l) * prime64_1

		acc += mulFold64(readU64(p, 0*16+0)^key64_000, readU64(p, 0*16+8)^key64_008)
		acc += mulFold64(readU64(p, 1*16+0)^key64_016, readU64(p, 1*16+8)^key64_024)
		acc += mulFold64(readU64(p, 2*16+0)^key64_032, readU64(p, 2*16+8)^key64_040)
		acc += mulFold64(readU64(p, 3*16+0)^key64_048, readU64(p, 3*16+8)^key64_056)
		acc += mulFold64(readU64(p, 4*16+0)^key64_064, readU64(p, 4*16+8)^key64_072)
		acc += mulFold64(readU64(p, 5*16+0)^key64_080, readU64(p, 5*16+8)^key64_088)
		acc += mulFold64(readU64(p, 6*16+0)^key64_096, readU64(p, 6*16+8)^key64_104)
		acc += mulFold64(readU64(p, 7*16+0)^key64_112, readU64(p, 7*16+8)^key64_120)

		// avalanche
		acc = Xxh3Avalanche(acc)

		// trailing groups after 128
		top := ui(l) &^ 15
		for i := ui(8 * 16); i < top; i += 16 {
			acc += mulFold64(readU64(p, i+0)^readU64(key, i-125), readU64(p, i+8)^readU64(key, i-117))
		}

		// last 16 bytes
		acc += mulFold64(readU64(p, ui(l)-16)^key64_119, readU64(p, ui(l)-8)^key64_127)

		return Xxh3Avalanche(acc)

	}
}


func Casteame[T any](elemento any) T{
	return *(*T)(ptr(&elemento))
}

func Hash(b []byte) uint64 {
	return hashAny(*(*GeneralBytes)(ptr(&b)))
}

func HashGeneral[T any](b T) uint64 {
	size := binary.Size(b)
	if size == -1 {
		panic("WAS NOT FIXED SIZE?")
	}

	var byteBuff bytes.Buffer

	binary.Write(&byteBuff, binary.BigEndian, b)
	return Hash(byteBuff.Bytes())
}

func HashShort(b []byte) uint64 {
	return hashAny(*(*GeneralBytes)(ptr(&b)))

}

