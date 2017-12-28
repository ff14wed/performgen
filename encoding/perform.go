package encoding

import "encoding/hex"

// Perform defines the struct for a perform block
type Perform struct {
	Length byte
	Data   [30]byte
	U1     byte
}

// String returns the hexadecimal string representation of the 32 bytes that
// comprise the Perform block.
func (p *Perform) String() string {
	buf := make([]byte, 64)
	hex.Encode(buf[0:2], []byte{p.Length})
	hex.Encode(buf[2:62], p.Data[:])
	hex.Encode(buf[62:64], []byte{p.U1})
	return string(buf)
}
