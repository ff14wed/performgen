package mml_test

import (
	"github.com/ff14wed/performgen/encoding"
	"github.com/ff14wed/performgen/mml"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("State", func() {
	var s *mml.State
	BeforeEach(func() {
		s = new(mml.State)
	})
	Describe("EmitNote", func() {
		DescribeTable("emits the correct note",
			func(note, modifier string, octave, expectedID int) {
				s.SetOctave(octave + 3)
				s.EmitNote(note, modifier, 0)
				Expect(s.Sequence).To(ConsistOf(encoding.Note(byte(expectedID)), encoding.Delay(20)))
			},
			Entry("C (-1)", "C", "", -1, 1),
			Entry("C# (-1)", "C", "#", -1, 2),
			Entry("C# (-1)", "C", "+", -1, 2),
			Entry("Db (-1)", "D", "-", -1, 2),
			Entry("D (-1)", "D", "", -1, 3),
			Entry("D# (-1)", "D", "#", -1, 4),
			Entry("D# (-1)", "D", "+", -1, 4),
			Entry("Eb (-1)", "E", "-", -1, 4),
			Entry("E (-1)", "E", "", -1, 5),
			Entry("F (-1)", "F", "", -1, 6),
			Entry("F# (-1)", "F", "#", -1, 7),
			Entry("F# (-1)", "F", "+", -1, 7),
			Entry("Gb (-1)", "G", "-", -1, 7),
			Entry("G (-1)", "G", "", -1, 8),
			Entry("G# (-1)", "G", "#", -1, 9),
			Entry("G# (-1)", "G", "+", -1, 9),
			Entry("Ab (-1)", "A", "-", -1, 9),
			Entry("A (-1)", "A", "", -1, 10),
			Entry("A# (-1)", "A", "+", -1, 11),
			Entry("A# (-1)", "A", "+", -1, 11),
			Entry("Bb (-1)", "B", "-", -1, 11),
			Entry("B (-1)", "B", "", -1, 12),

			Entry("C", "C", "", 0, 13),
			Entry("C#", "C", "#", 0, 14),
			Entry("C#", "C", "+", 0, 14),
			Entry("Db", "D", "-", 0, 14),
			Entry("D", "D", "", 0, 15),
			Entry("D#", "D", "#", 0, 16),
			Entry("D#", "D", "+", 0, 16),
			Entry("Eb", "E", "-", 0, 16),
			Entry("E", "E", "", 0, 17),
			Entry("F", "F", "", 0, 18),
			Entry("F#", "F", "#", 0, 19),
			Entry("F#", "F", "+", 0, 19),
			Entry("Gb", "G", "-", 0, 19),
			Entry("G", "G", "", 0, 20),
			Entry("G#", "G", "#", 0, 21),
			Entry("G#", "G", "+", 0, 21),
			Entry("Ab", "A", "-", 0, 21),
			Entry("A", "A", "", 0, 22),
			Entry("A#", "A", "+", 0, 23),
			Entry("A#", "A", "+", 0, 23),
			Entry("Bb", "B", "-", 0, 23),
			Entry("B", "B", "", 0, 24),

			Entry("C (+1)", "C", "", 1, 25),
			Entry("C# (+1)", "C", "#", 1, 26),
			Entry("C# (+1)", "C", "+", 1, 26),
			Entry("Db (+1)", "D", "-", 1, 26),
			Entry("D (+1)", "D", "", 1, 27),
			Entry("D# (+1)", "D", "#", 1, 28),
			Entry("D# (+1)", "D", "+", 1, 28),
			Entry("Eb (+1)", "E", "-", 1, 28),
			Entry("E (+1)", "E", "", 1, 29),
			Entry("F (+1)", "F", "", 1, 30),
			Entry("F# (+1)", "F", "#", 1, 31),
			Entry("F# (+1)", "F", "+", 1, 31),
			Entry("Gb (+1)", "G", "-", 1, 31),
			Entry("G (+1)", "G", "", 1, 32),
			Entry("G# (+1)", "G", "#", 1, 33),
			Entry("G# (+1)", "G", "+", 1, 33),
			Entry("Ab (+1)", "A", "-", 1, 33),
			Entry("A (+1)", "A", "", 1, 34),
			Entry("A# (+1)", "A", "+", 1, 35),
			Entry("A# (+1)", "A", "+", 1, 35),
			Entry("Bb (+1)", "B", "-", 1, 35),
			Entry("B (+1)", "B", "", 1, 36),

			Entry("C (+2)", "C", "", 2, 37),
		)
		It("accepts lowercase notes", func() {
			Expect(s.EmitNote("a", "", 0)).To(Succeed())
			Expect(s.Sequence).To(ConsistOf(encoding.Note(22), encoding.Delay(20)))
		})
		It("errors if an invalid length is provided", func() {
			Expect(s.EmitNote("C", "", -2)).To(MatchError("invalid length: -2"))
		})
		It("errors if an invalid note is provided", func() {
			Expect(s.EmitNote("H", "#", 1)).To(MatchError("invalid note: H#"))
		})
		It("errors if a given note is out of range", func() {
			s.SetOctave(5)
			Expect(s.EmitNote("D", "+", 1)).To(MatchError("invalid note: D+ at octave 5"))
		})
	})
	Describe("EmitRest", func() {
		It("emits a default (quarter note) rest at 120bpm", func() {
			Expect(s.EmitRest(-1)).To(Succeed())
			Expect(s.Sequence).To(ConsistOf(encoding.Delay(250), encoding.Delay(250)))
		})
		It("emits a half note rest at 120bpm", func() {
			Expect(s.EmitRest(2)).To(Succeed())
			Expect(s.Sequence).To(ConsistOf(
				encoding.Delay(250),
				encoding.Delay(250),
				encoding.Delay(250),
				encoding.Delay(250),
			))
		})
		It("errors if an invalid length is provided", func() {
			Expect(s.EmitRest(-2)).To(MatchError("invalid length: -2"))
		})
		Context("with the default tempo and a default length set", func() {
			BeforeEach(func() {
				s.SetDefaultLength(32)
			})
			It("emits a thirty-secondth note rest at 120bpm", func() {
				Expect(s.EmitRest(-1)).To(Succeed())
				Expect(s.Sequence).To(ConsistOf(encoding.Delay(62)))
			})
		})
		Context("with provided tempo and provided length", func() {
			BeforeEach(func() {
				s.SetTempo(60)
			})
			It("emits a half note rest at 60bpm", func() {
				Expect(s.EmitRest(2)).To(Succeed())
				Expect(s.Sequence).To(ConsistOf(
					encoding.Delay(250),
					encoding.Delay(250),
					encoding.Delay(250),
					encoding.Delay(250),
					encoding.Delay(250),
					encoding.Delay(250),
					encoding.Delay(250),
					encoding.Delay(250),
				))
			})
		})
	})
	Describe("SetTempo", func() {
		It("sets the tempo on the state", func() {
			Expect(s.SetTempo(900)).To(Succeed())
			Expect(s.Tempo).To(Equal(900))
		})
		It("errors if the tempo is less than 1", func() {
			Expect(s.SetTempo(0)).To(MatchError("cannot set tempo to lower than 1"))
			Expect(s.Tempo).To(Equal(0))
		})
		It("errors if the tempo is greater than 900", func() {
			// Arbitrary cutoff point, but a BPM too fast doesn't really produce
			// meaningful output anyways
			Expect(s.SetTempo(901)).To(MatchError("cannot set tempo to greater than 900"))
			Expect(s.Tempo).To(Equal(0))
		})
	})
	Describe("SetDefaultLength", func() {
		It("sets the default length on the state", func() {
			Expect(s.SetDefaultLength(64)).To(Succeed())
			Expect(s.Length).To(Equal(64))
		})
		It("errors if the default length is set to 0", func() {
			Expect(s.SetDefaultLength(0)).To(MatchError("cannot set default length to 0"))
			Expect(s.Length).To(Equal(0))
		})
		It("errors if the default length is set to a negative number", func() {
			Expect(s.SetDefaultLength(-1)).To(MatchError("cannot set default length to less than 0"))
			Expect(s.Length).To(Equal(0))
		})
		It("errors if the default length is set to greater than 64", func() {
			Expect(s.SetDefaultLength(65)).To(MatchError("cannot set default length to greater than 64"))
			Expect(s.Length).To(Equal(0))
		})
	})
	Describe("SetOctave", func() {
		DescribeTable("allowed octaves",
			func(octave int) {
				Expect(s.SetOctave(octave)).To(Succeed())
				Expect(s.Octave).To(Equal(octave))
			},
			Entry("2", 2),
			Entry("3", 3),
			Entry("4", 4),
			Entry("5", 5),
		)
		It("errors if the octave is set to less than 2", func() {
			Expect(s.SetOctave(1)).To(MatchError("cannot set octave to anything other than 2, 3, 4, or 5"))
			Expect(s.Length).To(Equal(0))
		})
		It("errors if the octave is set to greater than 5", func() {
			Expect(s.SetOctave(6)).To(MatchError("cannot set octave to anything other than 2, 3, 4, or 5"))
			Expect(s.Length).To(Equal(0))
		})
	})
	Describe("CurrentOctave", func() {
		It("returns the current octave", func() {
			s.Octave = 9000
			Expect(s.CurrentOctave()).To(Equal(9000))
		})
	})
})
