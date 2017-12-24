package mml_test

import (
	"errors"

	"github.com/ff14wed/performgen/mml"
	"github.com/ff14wed/performgen/mml/mmlfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var fooError = errors.New("foo")
var _ = Describe("Commands", func() {
	var fakeExecutor *mmlfakes.Executor
	BeforeEach(func() {
		fakeExecutor = new(mmlfakes.Executor)
	})
	Describe("NoteCommand", func() {
		var c *mml.NoteCommand
		BeforeEach(func() {
			c = &mml.NoteCommand{
				Note:     "A",
				Modifier: "#",
				Length:   64,
				Dot:      true,
			}
		})
		It("emits a note on the state", func() {
			Expect(c.Execute(fakeExecutor)).To(Succeed())
			Expect(fakeExecutor.EmitNoteCallCount()).To(Equal(1))
			note, modifier, length, dot := fakeExecutor.EmitNoteArgsForCall(0)
			Expect(note).To(Equal("A"))
			Expect(modifier).To(Equal("#"))
			Expect(length).To(Equal(64))
			Expect(dot).To(BeTrue())
		})
		Context("when the state emits an error", func() {
			BeforeEach(func() {
				fakeExecutor.EmitNoteReturns(fooError)
			})
			It("command returns the same error", func() {
				Expect(c.Execute(fakeExecutor)).To(MatchError(fooError))
			})
		})
	})
	Describe("RestCommand", func() {
		var c *mml.RestCommand
		BeforeEach(func() {
			c = &mml.RestCommand{
				Length: 64,
				Dot:    true,
			}
		})
		It("emits a rest on the state", func() {
			Expect(c.Execute(fakeExecutor)).To(Succeed())
			Expect(fakeExecutor.EmitRestCallCount()).To(Equal(1))
			length, dot := fakeExecutor.EmitRestArgsForCall(0)
			Expect(length).To(Equal(64))
			Expect(dot).To(BeTrue())
		})
		Context("when the state emits an error", func() {
			BeforeEach(func() {
				fakeExecutor.EmitRestReturns(fooError)
			})
			It("command returns the same error", func() {
				Expect(c.Execute(fakeExecutor)).To(MatchError(fooError))
			})
		})
	})
	Describe("TempoCommand", func() {
		var c *mml.TempoCommand
		BeforeEach(func() {
			c = &mml.TempoCommand{
				Tempo: 120,
			}
		})
		It("sets the tempo on the state", func() {
			Expect(c.Execute(fakeExecutor)).To(Succeed())
			Expect(fakeExecutor.SetTempoCallCount()).To(Equal(1))
			tempo := fakeExecutor.SetTempoArgsForCall(0)
			Expect(tempo).To(Equal(120))
		})
		Context("when the state emits an error", func() {
			BeforeEach(func() {
				fakeExecutor.SetTempoReturns(fooError)
			})
			It("command returns the same error", func() {
				Expect(c.Execute(fakeExecutor)).To(MatchError(fooError))
			})
		})
	})
	Describe("LengthCommand", func() {
		var c *mml.LengthCommand
		BeforeEach(func() {
			c = &mml.LengthCommand{
				Length: 1,
				Dot:    true,
			}
		})
		It("sets the default length on the state", func() {
			Expect(c.Execute(fakeExecutor)).To(Succeed())
			Expect(fakeExecutor.SetDefaultLengthCallCount()).To(Equal(1))
			length, dot := fakeExecutor.SetDefaultLengthArgsForCall(0)
			Expect(length).To(Equal(1))
			Expect(dot).To(BeTrue())
		})
		Context("when the state emits an error", func() {
			BeforeEach(func() {
				fakeExecutor.SetDefaultLengthReturns(fooError)
			})
			It("command returns the same error", func() {
				Expect(c.Execute(fakeExecutor)).To(MatchError(fooError))
			})
		})
	})
	Describe("OctaveCommand", func() {
		var c *mml.OctaveCommand
		BeforeEach(func() {
			c = &mml.OctaveCommand{
				Octave: -1,
			}
		})
		It("sets the octave on the state", func() {
			Expect(c.Execute(fakeExecutor)).To(Succeed())
			Expect(fakeExecutor.SetOctaveCallCount()).To(Equal(1))
			octave := fakeExecutor.SetOctaveArgsForCall(0)
			Expect(octave).To(Equal(-1))
		})
		Context("when the state emits an error", func() {
			BeforeEach(func() {
				fakeExecutor.SetOctaveReturns(fooError)
			})
			It("command returns the same error", func() {
				Expect(c.Execute(fakeExecutor)).To(MatchError(fooError))
			})
		})
	})
	Describe("OctaveUpCommand", func() {
		var c *mml.OctaveUpCommand
		BeforeEach(func() {
			c = &mml.OctaveUpCommand{}
			fakeExecutor.CurrentOctaveReturns(1)
		})
		It("sets the octave on the state", func() {
			Expect(c.Execute(fakeExecutor)).To(Succeed())
			Expect(fakeExecutor.SetOctaveCallCount()).To(Equal(1))
			octave := fakeExecutor.SetOctaveArgsForCall(0)
			Expect(octave).To(Equal(2))
		})
		Context("when the state emits an error", func() {
			BeforeEach(func() {
				fakeExecutor.SetOctaveReturns(fooError)
			})
			It("command returns the same error", func() {
				Expect(c.Execute(fakeExecutor)).To(MatchError(fooError))
			})
		})
	})
	Describe("OctaveDownCommand", func() {
		var c *mml.OctaveDownCommand
		BeforeEach(func() {
			c = &mml.OctaveDownCommand{}
			fakeExecutor.CurrentOctaveReturns(-1)
		})
		It("sets the octave on the state", func() {
			Expect(c.Execute(fakeExecutor)).To(Succeed())
			Expect(fakeExecutor.SetOctaveCallCount()).To(Equal(1))
			octave := fakeExecutor.SetOctaveArgsForCall(0)
			Expect(octave).To(Equal(-2))
		})
		Context("when the state emits an error", func() {
			BeforeEach(func() {
				fakeExecutor.SetOctaveReturns(fooError)
			})
			It("command returns the same error", func() {
				Expect(c.Execute(fakeExecutor)).To(MatchError(fooError))
			})
		})
	})
	Describe("NoOpCommand", func() {
		var n *mml.NoOpCommand
		BeforeEach(func() {
			n = &mml.NoOpCommand{}
		})
		It("does nothing", func() {
			Expect(n.Execute(fakeExecutor)).To(Succeed())
			Expect(fakeExecutor.Invocations()).To(BeEmpty())
		})
	})
})
