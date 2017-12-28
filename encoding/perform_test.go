package encoding_test

import (
	"github.com/ff14wed/performgen/encoding"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Perform", func() {
	Describe("String", func() {
		It("serializes the perform block to a hex string", func() {
			p := encoding.Perform{
				Length: 0xa1,
				Data: [30]byte{
					1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
				},
				U1: 0xb2,
			}
			pString := p.String()
			Expect(pString).To(Equal("a10102030405060708090a0b0c0d0e0f101112131400000000000000000000b2"))
			Expect(len(pString)).To(Equal(64))
		})
	})
})
