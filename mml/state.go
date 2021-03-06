package mml

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ff14wed/performgen/encoding"
)

//go:generate counterfeiter -o mmlfakes/executor.go --fake-name Executor . Executor

// Executor is an interface on which state changes can be executed
type Executor interface {
	EmitNote(note string, modifier string, length int, dot bool) error
	EmitRest(length int, dot bool) error

	SetTempo(t int) error
	SetDefaultLength(l int, dot bool) error
	SetOctave(o int) error
	CurrentOctave() int
}

// State describes the state machine that consumes state changes and emits
// either notes or rests
type State struct {
	Sequence encoding.Sequence
	Tempo    int
	Length   int
	Octave   int

	dottedLength bool
}

var _ Executor = new(State)

var noteMappings = map[string]int{
	"C": 1,
	"D": 3,
	"E": 5,
	"F": 6,
	"G": 8,
	"A": 10,
	"B": 12,
}

// EmitNote emits a music note to the sequence.
// Modifiers can be one of:
// `+` or `#` - Makes this note a sharp note
// `-` - Makes this note a flat note
// Length is the denominator of 1/x, where the note will be spaced from the
// next note by 1/x of a beat.
// If length is -1 (empty length code), the default length will be used.
// If length is 0 (explicit length code of 0), the length will be set to a
// very small value (20 milliseconds).
// If an octave was not specified previously, it will default to octave 3
func (s *State) EmitNote(note string, modifier string, length int, dot bool) error {
	shift := (s.CurrentOctave() - 3) * 12
	noteMap, ok := noteMappings[strings.ToUpper(note)]
	if !ok {
		return fmt.Errorf("invalid note: %s%s", note, modifier)
	}
	pos := byte(noteMap + shift)
	if modifier == "#" || modifier == "+" {
		pos++
	} else if modifier == "-" {
		pos--
	}
	if pos < 1 || pos > 37 {
		return fmt.Errorf("invalid note: %s%s at octave %d", note, modifier, s.CurrentOctave())
	}
	s.Sequence = append(s.Sequence, encoding.Note(pos))
	return s.EmitRest(length, dot)
}

// EmitRest emits a rest note to the sequence. The length is the same as
// the length defined by EmitNote.
func (s *State) EmitRest(length int, dot bool) error {
	ml, err := s.lengthInMs(length)
	if err != nil {
		return err
	}
	if length == -1 && s.dottedLength {
		ml = ml + ml/2
	}
	if dot {
		ml = ml + ml/2
	}
	s.emitDelay(ml)
	return nil
}

func (s *State) emitDelay(ml uint16) {
	for ml > 0 {
		if ml >= 250 {
			s.Sequence = append(s.Sequence, encoding.Delay(250))
			ml -= 250
		} else {
			s.Sequence = append(s.Sequence, encoding.Delay(byte(ml)))
			ml = 0
		}
	}
}

// SetTempo sets the tempo (in BPM) on the state. If the Tempo is not set,
// it is assumed the tempo is 120 bpm.
func (s *State) SetTempo(t int) error {
	if t < 1 {
		return errors.New("cannot set tempo to lower than 1")
	} else if t > 900 {
		return errors.New("cannot set tempo to greater than 900")
	}
	s.Tempo = t
	return nil
}

// SetDefaultLength sets the default length on the state. If the default length
// is not set, then it is assumed the length is 1/4th of a beat.
func (s *State) SetDefaultLength(l int, dot bool) error {
	switch {
	case l < 0:
		return errors.New("cannot set default length to less than 0")
	case l == 0:
		return errors.New("cannot set default length to 0")
	case l > 64:
		return errors.New("cannot set default length to greater than 64")
	}
	s.Length = l
	s.dottedLength = dot
	return nil
}

// SetOctave sets the octave on the state
func (s *State) SetOctave(o int) error {
	if o < 3 || o > 6 {
		return errors.New("cannot set octave to anything other than 3, 4, 5, or 6")
	}
	s.Octave = o
	return nil
}

// CurrentOctave returns the current octave on the state
func (s *State) CurrentOctave() int {
	if s.Octave == 0 {
		s.Octave = 4
	}
	return s.Octave
}

// lengthInMs calculates the amount of delay required to achieve a length
// given a certain tempo
func (s *State) lengthInMs(lengthDenom int) (uint16, error) {
	switch {
	case lengthDenom == 0:
		return 20, nil
	case lengthDenom < -1:
		return 0, fmt.Errorf("invalid length: %d", lengthDenom)
	case lengthDenom == -1:
		lengthDenom = 4
		if s.Length != 0 {
			lengthDenom = s.Length
		}
	}
	tempo := 120
	if s.Tempo != 0 {
		tempo = s.Tempo
	}
	msPerBeat := 60000.0 / float64(tempo)
	length := 4 * msPerBeat / float64(lengthDenom)
	return uint16(length), nil
}
