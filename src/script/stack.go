package script
import (
	"bytes"
	"encoding/binary"
)

type Stack struct {
	Contents []Operand
}

// Top Return reference to the top item of the stack
func (s *Stack) Top() Operand {
	return s.Contents[len(s.Contents) - 1]
}

// Second Return reference to the second item from the stack top
func (s *Stack) Second() Operand {
	return s.Contents[len(s.Contents) - 2 ]
}

// Pop Pops one item of the stack
func (s *Stack) Pop() {
	s.Contents = s.Contents[:len(s.Contents) -1]
}

// PopTwo Pops two items of the stack
func (s *Stack) PopTwo() {
	s.Contents = s.Contents[:len(s.Contents) -2]
}

func (s *Stack) Push(item Operand) {
	s.Contents = append(s.Contents, item)
}

// DuplicateTop duplicates the top items of the
// stack and pushes it on top
func (s *Stack) DuplicateTop() {
	top := s.Top()
	switch v := top.(type) {
			case PUB_KEY_V1:
				s.Push(PUB_KEY_V1{ Key: append([]byte{}, v.Key...) })
			default:
				panic("Unexpected duplicating top stack item %T!\n")

	}
}

// Ser Serialise a Stack
func (s Stack) Ser() *bytes.Buffer {
	var ret bytes.Buffer
	for _, op := range s.Contents {
			binary.Write(&ret, littleEndian, op.AsByte())
			binary.Write(&ret, littleEndian, op.LenData())
			if op.LenData() > 0x00 {
				binary.Write(&ret, littleEndian, op.Data())
			}
	}
	return &ret
}

// ListTypes List the types in a stack
func (s * Stack) ListTypes() []string {
	ret := []string{}
	for _, op := range s.Contents {
		ret = append(ret, op.Name())
	}
	return ret
}

// consumeByte Consume a byte from a buffer
func consumeByte(b *bytes.Buffer) {
	var lenToRead = byte(0)
	binary.Read(b, littleEndian, lenToRead)
}

// Marshall Given a bytes.Buffer marshall to a Stack
// TODO: Test marshalling!
func Marshall(b *bytes.Buffer) Stack {
	ret := []Operand{}
	for b.Len() > 0 {
		var readByte byte
		binary.Read(b, littleEndian, &readByte)
		switch readByte {
		case OP_EQUALVERIFY{}.AsByte():
			consumeByte(b)
			ret = append(ret, OP_EQUALVERIFY{})
		case OP_CHECKSIG{}.AsByte():
			consumeByte(b)
			ret = append(ret, OP_CHECKSIG{})
		case OP_DUP{}.AsByte():
			consumeByte(b)
			ret = append(ret, OP_DUP{})
		case OP_HASH_160{}.AsByte():
			consumeByte(b)
			ret = append(ret, OP_HASH_160{})
		case PUB_KEY_HASH{}.AsByte():
			var lenToRead = byte(0)
			binary.Read(b, littleEndian, &lenToRead)
			var buffer []byte
			for i := byte(0); i < lenToRead; i++ {
				var read byte
				binary.Read(b, littleEndian, &read)
				buffer = append(buffer, read)
			}
			ret = append(ret, PUB_KEY_HASH{buffer})
		case SIG{}.AsByte():
			var lenToRead = byte(0)
			binary.Read(b, littleEndian, &lenToRead)
			var buffer []byte
			for i := byte(0); i < lenToRead; i++ {
				var read byte
				binary.Read(b, littleEndian, &read)
				buffer = append(buffer, read)
			}
			ret = append(ret, SIG{buffer})
		case PUB_KEY_V1{}.AsByte():
			var lenToRead = byte(0)
			binary.Read(b, littleEndian, &lenToRead)
			var buffer []byte
			for i := byte(0); i < lenToRead; i++ {
				var read byte
				binary.Read(b, littleEndian, &read)
				buffer = append(buffer, read)
			}
			ret = append(ret, PUB_KEY_V1{buffer})
		}
	}
	return Stack{ret}
}
