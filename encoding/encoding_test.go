package encoding_test

import (
	"github.com/ff14wed/performgen/encoding"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Encoding", func() {
	Describe("Note", func() {
		It("encodes to a single byte", func() {
			n := encoding.Note(12)
			Expect(n.Encode()).To(ConsistOf(byte(12)))
		})
	})
	Describe("Delay", func() {
		It("encodes to a two bytes", func() {
			d := encoding.Delay(128)
			Expect(d.Encode()).To(Equal([]byte{0xFF, 128}))
		})
	})
	Describe("Sequence", func() {
		It("encodes to multiple perform data blocks", func() {
			s := encoding.Sequence{
				encoding.Note(1), encoding.Delay(0x80),
				encoding.Note(2), encoding.Delay(0x80),
				encoding.Note(3), encoding.Delay(0x80),
				encoding.Note(4), encoding.Delay(0x80),
				encoding.Note(5), encoding.Delay(0x80),
				encoding.Delay(0xFF), encoding.Delay(0xFF), encoding.Delay(0xFF), encoding.Delay(0xFF),
				encoding.Note(6), encoding.Delay(0x80),
				encoding.Note(7), encoding.Delay(0x80),
				encoding.Note(8), encoding.Delay(0x80),
				encoding.Note(9), encoding.Delay(0x80),
				encoding.Note(10), encoding.Delay(0x80),
			}
			Expect(s.Blocks()).To(Equal([]*encoding.Perform{
				&encoding.Perform{
					Length: 30,
					Data: [30]byte{
						1, 0xFF, 0x80, 2, 0xFF, 0x80, 3, 0xFF, 0x80, 4, 0xFF, 0x80, 5, 0xFF, 0x80,
						0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
						6, 0xFF, 0x80, 7, 0xFF, 0x80, 8,
					},
				},
				&encoding.Perform{
					Length: 8,
					Data:   [30]byte{0xFF, 0x80, 9, 0xFF, 0x80, 10, 0xFF, 0x80},
				},
			}))
		})
	})
})
