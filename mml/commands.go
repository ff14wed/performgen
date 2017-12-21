package mml

// NoteCommand emits a note with a certain length
type NoteCommand struct {
	Note     string
	Modifier string
	Length   int
	Dot      bool
}

// Execute emits a note on the state
func (n *NoteCommand) Execute(e Executor) error {
	return e.EmitNote(n.Note, n.Modifier, n.Length, n.Dot)
}

// RestCommand emits a rest with a certain length
type RestCommand struct {
	Length int
	Dot    bool
}

// Execute emits a rest on the state
func (r *RestCommand) Execute(e Executor) error {
	return e.EmitRest(r.Length, r.Dot)
}

// TempoCommand sets the tempo
type TempoCommand struct {
	Tempo int
}

// Execute sets the tempo on the state
func (t *TempoCommand) Execute(e Executor) error {
	return e.SetTempo(t.Tempo)
}

// LengthCommand sets the default length
type LengthCommand struct {
	Length int
	Dot    bool
}

// Execute sets the default length on the state
func (l *LengthCommand) Execute(e Executor) error {
	return e.SetDefaultLength(l.Length, l.Dot)
}

// OctaveCommand sets the octave
type OctaveCommand struct {
	Octave int
}

// Execute sets the octave on the state
func (o *OctaveCommand) Execute(e Executor) error {
	return e.SetOctave(o.Octave)
}

// OctaveUpCommand increments the octave
type OctaveUpCommand struct{}

// Execute increments the octave on the state
func (o *OctaveUpCommand) Execute(e Executor) error {
	return e.SetOctave(e.CurrentOctave() + 1)
}

// OctaveDownCommand decrements the octave
type OctaveDownCommand struct{}

// Execute decrements the octave on the state
func (o *OctaveDownCommand) Execute(e Executor) error {
	return e.SetOctave(e.CurrentOctave() - 1)
}
