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
			input = bytes.NewReader([]byte("T120 L4 O0 > < Aa1 B# B3 C#4 \n C+4 C-4 D0d0 R E5. r5. L4. v127"))
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
				&mml.NoteCommand{Note: "E", Length: 5, Dot: true},
				&mml.RestCommand{Length: 5, Dot: true},
				&mml.LengthCommand{Length: 4, Dot: true},
				&mml.NoOpCommand{},
			}))
			Expect(ast.Positions).To(Equal([]mml.Position{
				{Line: 1, Column: 1},
				{Line: 1, Column: 6},
				{Line: 1, Column: 9},
				{Line: 1, Column: 12},
				{Line: 1, Column: 14},
				{Line: 1, Column: 16},
				{Line: 1, Column: 17},
				{Line: 1, Column: 20},
				{Line: 1, Column: 23},
				{Line: 1, Column: 26},
				{Line: 2, Column: 2},
				{Line: 2, Column: 6},
				{Line: 2, Column: 10},
				{Line: 2, Column: 12},
				{Line: 2, Column: 15},
				{Line: 2, Column: 17},
				{Line: 2, Column: 21},
				{Line: 2, Column: 25},
				{Line: 2, Column: 29},
			}))
		})
	})
	DescribeTable("commands with required numeric arguments should error when not provided a numeric argument",
		func(command, input string) {
			reader := bytes.NewReader([]byte(input))
			parser := mml.NewParser(reader)
			_, err := parser.Parse()
			Expect(err).To(MatchError(command + " command at line 1, column 5: expected numeric argument"))
		},
		Entry("Tempo Command", "Tempo", "    T a"),
		Entry("Length Command", "Length", "    L a"),
		Entry("Octave Command", "Octave", "    O a"),
		Entry("Volume Command", "Volume", "    V a"),
	)
	DescribeTable("unrecognized tokens in various places should error",
		func(input string, location mml.Position) {
			reader := bytes.NewReader([]byte(input))
			parser := mml.NewParser(reader)
			_, err := parser.Parse()
			Expect(err).To(MatchError("invalid token 'H' at " + location.String()))
		},
		Entry("at the beginning of the input", "    HAaBbCcDd", mml.Position{Line: 1, Column: 5}),
		Entry("in the middle of the input", "    AaBbHCcDd", mml.Position{Line: 1, Column: 9}),
		Entry("after a note modifier", "    AaBb+HCcDd", mml.Position{Line: 1, Column: 10}),
		Entry("after a note length", "    AaBb+3HCcDd", mml.Position{Line: 1, Column: 11}),
		Entry("after a note dot", "    AaBb+3.HCcDd", mml.Position{Line: 1, Column: 12}),
		Entry("after a rest dot", "    AaBr.HCcDd", mml.Position{Line: 1, Column: 10}),
		Entry("after a length dot", "    AaBL4.HCcDd", mml.Position{Line: 1, Column: 11}),
		Entry("after a length dot", "    AaBL4.HCcDd", mml.Position{Line: 1, Column: 11}),
	)
	Context("when there is a non-command token after a full command specificaiton", func() {
		BeforeEach(func() {
			input = bytes.NewReader([]byte("    O4+"))
		})
		It("errors", func() {
			parser := mml.NewParser(input)
			_, err := parser.Parse()
			Expect(err).To(MatchError("expected command, got '+' at line 1, column 7"))
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
		Entry("Volume", "    V9223372036854775808"),
	)
})
