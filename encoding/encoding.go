package encoding

import "time"

// Sequence is a series of steps and can possibly be encoded in a single or
// more than one Perform block
type Sequence []Step

// Step defines a single step (note or delay) in the wire format for a sequence
// for a Perform data block
type Step interface {
	Encode() []byte
	Length() time.Duration
}

// Note is a step that encodes the Perform key ID for a note
type Note byte

// Encode encodes a note to its wire format (PERFORM_KEY_ID)
func (n Note) Encode() []byte { return []byte{byte(n)} }

// Length determines the length in time of the encoded note
func (n Note) Length() time.Duration { return 0 }

// Delay is a step that encodes the number of milliseconds of delay (1 - 250)
type Delay byte

// Encode encodes a delay to its wire format (0xFF DELAY)
func (r Delay) Encode() []byte { return []byte{0xFF, byte(r)} }

// Length determines the length in time of the encoded delay
func (r Delay) Length() time.Duration { return time.Duration(r) * time.Millisecond }

// PerformSegment encapsulates a single block of a performance. It's not a
// measure. It only encapsulates what can fit in a single packet of data.
type PerformSegment struct {
	Block  *Perform
	Length time.Duration
}

// Segments returns the sequence of segments containing blocks that conform
// with the FFXIV RPC for performing a sequence of notes.
func (s Sequence) Segments() []PerformSegment {
	blocks := []PerformSegment{}

	buf := make([]byte, 0)
	length := time.Duration(0)
	for _, step := range s {
		stepBytes := step.Encode()
		if len(buf)+len(stepBytes) > 30 {
			blocks = append(blocks, PerformSegment{
				Block:  createBlock(buf),
				Length: length,
			})
			buf = make([]byte, 0)
			length = 0
		}
		buf = append(buf, stepBytes...)
		length = length + step.Length()
	}
	if len(buf) > 0 {
		blocks = append(blocks, PerformSegment{
			Block:  createBlock(buf),
			Length: length,
		})
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
