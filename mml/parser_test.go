package mml_test

import (
	"bytes"

	"github.com/ff14wed/performgen/mml"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Parser", func() {
	var input *bytes.Reader
	Context("with a valid program", func() {
		BeforeEach(func() {
			input = bytes.NewReader([]byte("T120 L4 O0 > < Aa1 B# B3 C#4 C+4 C-4 D0d0 Rr5"))
		})
		It("generates a correct syntax tree and notes have -1 as default length", func() {
			parser := mml.NewParser(input)
			ast, err := parser.Parse()
			Expect(err).ToNot(HaveOccurred())
			Expect(ast.Sequence).To(Equal([]mml.Command{
				&mml.TempoCommand{Tempo: 120},
				&mml.LengthCommand{Length: 4},
				&mml.OctaveCommand{Octave: 0},
				&mml.OctaveUpCommand{},
				&mml.OctaveDownCommand{},
				&mml.NoteCommand{Note: "A", Length: -1},
				&mml.NoteCommand{Note: "a", Length: 1},
				&mml.NoteCommand{Note: "B", Modifier: "#", Length: -1},
				&mml.NoteCommand{Note: "B", Length: 3},
				&mml.NoteCommand{Note: "C", Modifier: "#", Length: 4},
				&mml.NoteCommand{Note: "C", Modifier: "+", Length: 4},
				&mml.NoteCommand{Note: "C", Modifier: "-", Length: 4},
				&mml.NoteCommand{Note: "D", Length: 0},
				&mml.NoteCommand{Note: "d", Length: 0},
				&mml.RestCommand{Length: -1},
				&mml.RestCommand{Length: 5},
			}))
		})
	})
	Context("when tempo command does not have a numeric argument", func() {
		BeforeEach(func() {
			input = bytes.NewReader([]byte("    T a"))
		})
		It("errors", func() {
			parser := mml.NewParser(input)
			_, err := parser.Parse()
			Expect(err).To(MatchError("Tempo command at line 1 col 5: expected numeric argument"))
		})
	})
	Context("when length command does not have a numeric argument", func() {
		BeforeEach(func() {
			input = bytes.NewReader([]byte("    L a"))
		})
		It("errors", func() {
			parser := mml.NewParser(input)
			_, err := parser.Parse()
			Expect(err).To(MatchError("Length command at line 1 col 5: expected numeric argument"))
		})
	})
	Context("when octave command does not have a numeric argument", func() {
		BeforeEach(func() {
			input = bytes.NewReader([]byte("    O a"))
		})
		It("errors", func() {
			parser := mml.NewParser(input)
			_, err := parser.Parse()
			Expect(err).To(MatchError("Octave command at line 1 col 5: expected numeric argument"))
		})
	})
	Context("when there is an unrecognized token at the beginning of the input", func() {
		BeforeEach(func() {
			input = bytes.NewReader([]byte("    HAaBbCcDd"))
		})
		It("errors", func() {
			parser := mml.NewParser(input)
			_, err := parser.Parse()
			Expect(err).To(MatchError("invalid token 'H' at line 1 col 5"))
		})
	})
	Context("when there is an unrecognized token in middle of the input", func() {
		BeforeEach(func() {
			input = bytes.NewReader([]byte("    AaBbHCcDd"))
		})
		It("errors", func() {
			parser := mml.NewParser(input)
			_, err := parser.Parse()
			Expect(err).To(MatchError("invalid token 'H' at line 1 col 9"))
		})
	})
	Context("when there is an unrecognized token after a note modifier", func() {
		BeforeEach(func() {
			input = bytes.NewReader([]byte("    AaBb+HCcDd"))
		})
		It("errors", func() {
			parser := mml.NewParser(input)
			_, err := parser.Parse()
			Expect(err).To(MatchError("invalid token 'H' at line 1 col 10"))
		})
	})
	Context("when there is a non-command token after a full command specificaiton", func() {
		BeforeEach(func() {
			input = bytes.NewReader([]byte("    O4+"))
		})
		It("errors", func() {
			parser := mml.NewParser(input)
			_, err := parser.Parse()
			Expect(err).To(MatchError("expected command, got '+' at line 1 col 7"))
		})
	})
	DescribeTable("errors when a numeric argument is too large",
		func(inputProg string) {
			input = bytes.NewReader([]byte(inputProg))
			parser := mml.NewParser(input)
			_, err := parser.Parse()
			Expect(err).To(MatchError(ContainSubstring("value out of range")))
		},
		Entry("Note", "    A9223372036854775808"),
		Entry("Rest", "    R9223372036854775808"),
		Entry("Tempo", "    T9223372036854775808"),
		Entry("Length", "    L9223372036854775808"),
		Entry("Octave", "    O9223372036854775808"),
	)
})
