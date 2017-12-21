package performgen

import (
	"bytes"
	"fmt"

	"github.com/ff14wed/performgen/encoding"
	"github.com/ff14wed/performgen/mml"
)

// Generate converts MML to a packet that can be injected
// It returns the sequence of blocks that conform with the FFXIV Network RPC
// for performing a sequence of notes. These can be injected into the client
// to perform an entire track of music.
// The durations in these blocks are important because while multiple blocks
// can be injected at once to produce a contiguous track of music, each
// character has only a limited sized buffer for perform data that has yet
// to be used. A user should make use of these durations to inject them
// slowly into the client to prevent filling the buffer faster than data can be
// consumed.
func Generate(input string) ([]encoding.PerformSegment, error) {
	r := bytes.NewReader([]byte(input))
	parser := mml.NewParser(r)
	ast, err := parser.Parse()
	if err != nil {
		return nil, err
	}
	state := new(mml.State)
	for i, cmd := range ast.Sequence {
		err = cmd.Execute(state)
		if err != nil {
			return nil, fmt.Errorf("execution error at %s: %s", ast.Positions[i], err)
		}
	}
	return state.Sequence.Segments(), nil
}
