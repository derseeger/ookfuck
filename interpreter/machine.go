package interpreter

type Script interface {
	Sanitize()
	HasRemaining() bool
	NextInstruction(currentByte byte) byte
	Execute(interpreter Interpreter)
	SetSource(source []byte)
}

type Interpreter interface {
	IncrementAddress()
	DecrementAddress()
	IncrementByte()
	DecrementByte()
	ReadByte() byte
	SetByte(b byte)
}

type esotericInterpreter struct {
	cells   map[int]byte
	address int
}

func (e *esotericInterpreter) IncrementAddress() {
	e.address++
}

func (e *esotericInterpreter) DecrementAddress() {
	e.address--
}

func (e *esotericInterpreter) IncrementByte() {
	e.cells[e.address]++
}

func (e *esotericInterpreter) DecrementByte() {
	e.cells[e.address]--
}

func (e *esotericInterpreter) ReadByte() byte {
	return e.cells[e.address]
}

func (e *esotericInterpreter) SetByte(b byte) {
	e.cells[e.address] = b
}

func NewEsotericInterpreter() *esotericInterpreter {
	return &esotericInterpreter{
		cells: map[int]byte{},
	}
}
