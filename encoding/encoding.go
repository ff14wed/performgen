package encoding

// Sequence is a series of steps and can possibly be encoded in a single or
// more than one Perform block
type Sequence []Step

// Step defines a single step (note or delay) in the wire format for a sequence
// for a Perform data block
type Step interface {
	Encode() []byte
}

// Note is a step that encodes the Perform key ID for a note
type Note byte

// Encode encodes a note to its wire format (PERFORM_KEY_ID)
func (n Note) Encode() []byte { return []byte{byte(n)} }

// Delay is a step that encodes the number of milliseconds of delay (1 - 250)
type Delay byte

// Encode encodes a delay to its wire format (0xFF DELAY)
func (r Delay) Encode() []byte { return []byte{0xFF, byte(r)} }

// Blocks returns the sequence of blocks that conform with the FFXIV RPC
// for performing a sequence of notes
func (s Sequence) Blocks() []*Perform {
	blocks := []*Perform{}

	buf := make([]byte, 0)
	for _, step := range s {
		stepBytes := step.Encode()
		if len(buf)+len(stepBytes) > 30 {
			blocks = append(blocks, createBlock(buf))
			buf = make([]byte, 0)
		}
		buf = append(buf, stepBytes...)
	}
	if len(buf) > 0 {
		blocks = append(blocks, createBlock(buf))
	}

	return blocks
}

func createBlock(buf []byte) *Perform {
	newBlock := &Perform{
		Length: byte(len(buf)),
	}
	copy(newBlock.Data[:], buf)
	return newBlock
}
