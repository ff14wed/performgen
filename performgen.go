package performgen

import (
	"bytes"
	"fmt"

	"github.com/ff14wed/performgen/encoding"
	"github.com/ff14wed/performgen/mml"
)

// Generate converts MML to a packet that can be injected
func Generate(input string) ([]*encoding.Perform, error) {
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
	return state.Sequence.Blocks(), nil
}
