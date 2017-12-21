package performgen_test

import (
	"github.com/ff14wed/performgen"
	"github.com/ff14wed/performgen/encoding"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Perform Generator", func() {
	It("generates correct perform data blocks from the MML", func() {
		data, err := performgen.Generate("t88 b2al2b+.")
		Expect(err).ToNot(HaveOccurred())
		Expect(data).To(Equal([]*encoding.Perform{
			{
				Length: 29,
				Data:   [30]byte{24, 255, 250, 255, 250, 255, 250, 255, 250, 255, 250, 255, 113, 22, 255, 250, 255, 250, 255, 181, 25, 255, 250, 255, 250, 255, 250, 255, 250},
			},
			{
				Length: 10,
				Data:   [30]byte{255, 250, 255, 250, 255, 250, 255, 250, 255, 44},
				U1:     0,
			},
		}))
	})
	It("errors when invalid symbol is encountered", func() {
		_, err := performgen.Generate(" HABCD")
		Expect(err).To(MatchError("invalid token 'H' at line 1 col 2"))
	})
	It("errors when a runtime error has occurred", func() {
		_, err := performgen.Generate(" ABCDo7")
		Expect(err).To(MatchError("cannot set octave to anything other than 3, 4, 5, or 6"))
	})
})
