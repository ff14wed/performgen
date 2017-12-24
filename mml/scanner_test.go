package mml_test

import (
	"bytes"

	"github.com/ff14wed/performgen/mml"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type testTok struct {
	typ     mml.TokenType
	ident   string
	lineNum int
	colNum  int
}

var _ = Describe("Scanner", func() {
	var input *bytes.Reader

	Context("with notes in the input", func() {
		BeforeEach(func() {
			input = bytes.NewReader([]byte("AaBbCcDdEeFfGg"))
		})
		It("scans notes and the EOF", func() {
			scanner := mml.NewScanner(input)

			tokens := []string{
				"A", "a", "B", "b", "C", "c", "D", "d", "E", "e", "F", "f", "G", "g",
			}
			for i, t := range tokens {
				token := scanner.Scan()
				Expect(token.Type()).To(Equal(mml.TNote))
				Expect(token.Ident()).To(Equal(t))
				Expect(token.Position()).To(Equal(mml.Position{Line: 1, Column: i + 1}))
			}
			token := scanner.Scan()
			Expect(token.Type()).To(Equal(mml.TEOF))
			Expect(token.Position()).To(Equal(mml.Position{Line: 1, Column: 15}))
		})
	})
	Context("with whitespaces interspersed between notes", func() {
		BeforeEach(func() {
			input = bytes.NewReader([]byte("a  A  \nb "))
		})
		It("scans notes with the correct line and column numbers", func() {
			scanner := mml.NewScanner(input)

			expectedTokens := []testTok{
				testTok{typ: mml.TNote, ident: "a", lineNum: 1, colNum: 1},
				testTok{typ: mml.TNote, ident: "A", lineNum: 1, colNum: 4},
				testTok{typ: mml.TNote, ident: "b", lineNum: 2, colNum: 1},
				testTok{typ: mml.TEOF, ident: string(rune(0)), lineNum: 2, colNum: 3},
			}
			for _, tok := range expectedTokens {
				token := scanner.Scan()
				Expect(token.Type()).To(Equal(tok.typ))
				Expect(token.Ident()).To(Equal(tok.ident))
				Expect(token.Position()).To(Equal(mml.Position{Line: tok.lineNum, Column: tok.colNum}))
			}
		})
	})
	Context("with numerics, whitespace, and other tokens", func() {
		BeforeEach(func() {
			input = bytes.NewReader([]byte("r R+ tT#123 1\n456 lL- oO . >< 2 v V7 a&c8"))
		})
		It("scans with the correct line and column numbers", func() {
			scanner := mml.NewScanner(input)

			expectedTokens := []testTok{
				testTok{typ: mml.TRest, ident: "r", lineNum: 1, colNum: 1},
				testTok{typ: mml.TRest, ident: "R", lineNum: 1, colNum: 3},
				testTok{typ: mml.TModifier, ident: "+", lineNum: 1, colNum: 4},
				testTok{typ: mml.TTempo, ident: "t", lineNum: 1, colNum: 6},
				testTok{typ: mml.TTempo, ident: "T", lineNum: 1, colNum: 7},
				testTok{typ: mml.TModifier, ident: "#", lineNum: 1, colNum: 8},
				testTok{typ: mml.TNumeric, ident: "123", lineNum: 1, colNum: 9},
				testTok{typ: mml.TNumeric, ident: "1", lineNum: 1, colNum: 13},
				testTok{typ: mml.TNumeric, ident: "456", lineNum: 2, colNum: 1},
				testTok{typ: mml.TLength, ident: "l", lineNum: 2, colNum: 5},
				testTok{typ: mml.TLength, ident: "L", lineNum: 2, colNum: 6},
				testTok{typ: mml.TModifier, ident: "-", lineNum: 2, colNum: 7},
				testTok{typ: mml.TOctave, ident: "o", lineNum: 2, colNum: 9},
				testTok{typ: mml.TOctave, ident: "O", lineNum: 2, colNum: 10},
				testTok{typ: mml.TDot, ident: ".", lineNum: 2, colNum: 12},
				testTok{typ: mml.TOctaveUp, ident: ">", lineNum: 2, colNum: 14},
				testTok{typ: mml.TOctaveDown, ident: "<", lineNum: 2, colNum: 15},
				testTok{typ: mml.TNumeric, ident: "2", lineNum: 2, colNum: 17},
				testTok{typ: mml.TVolume, ident: "v", lineNum: 2, colNum: 19},
				testTok{typ: mml.TVolume, ident: "V", lineNum: 2, colNum: 21},
				testTok{typ: mml.TNumeric, ident: "7", lineNum: 2, colNum: 22},
				testTok{typ: mml.TNote, ident: "a", lineNum: 2, colNum: 24},
				testTok{typ: mml.TExtend, ident: "&", lineNum: 2, colNum: 25},
				testTok{typ: mml.TNote, ident: "c", lineNum: 2, colNum: 26},
				testTok{typ: mml.TNumeric, ident: "8", lineNum: 2, colNum: 27},
				testTok{typ: mml.TEOF, ident: string(rune(0)), lineNum: 2, colNum: 28},
			}
			for _, tok := range expectedTokens {
				token := scanner.Scan()
				Expect(token.Type()).To(Equal(tok.typ))
				Expect(token.Ident()).To(Equal(tok.ident))
				Expect(token.Position()).To(Equal(mml.Position{Line: tok.lineNum, Column: tok.colNum}))
			}
		})
	})
	Context("with unrecognized tokens", func() {
		BeforeEach(func() {
			input = bytes.NewReader([]byte("   HABCD"))
		})
		It("returns TIllegal for the first attempt, but returns notes afterwards", func() {
			scanner := mml.NewScanner(input)

			expectedTokens := []testTok{
				testTok{typ: mml.TIllegal, ident: "H", lineNum: 1, colNum: 4},
				testTok{typ: mml.TNote, ident: "A", lineNum: 1, colNum: 5},
				testTok{typ: mml.TNote, ident: "B", lineNum: 1, colNum: 6},
				testTok{typ: mml.TNote, ident: "C", lineNum: 1, colNum: 7},
				testTok{typ: mml.TNote, ident: "D", lineNum: 1, colNum: 8},
			}
			for _, tok := range expectedTokens {
				token := scanner.Scan()
				Expect(token.Type()).To(BeAssignableToTypeOf(tok.typ))
				Expect(token.Ident()).To(Equal(tok.ident))
				Expect(token.Position()).To(Equal(mml.Position{Line: tok.lineNum, Column: tok.colNum}))
			}
		})
	})
})
