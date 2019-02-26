package brainfuck

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/derseeger/ookfuck/interpreter"
	"log"
	"os"
	"regexp"
)

const clear_string_regexp = `[\+\-\<\>\[\]\,\.]*`

type BrainfuckScript struct {
	content []byte
	pos     int
}

func NewBrainfuckScript() *BrainfuckScript {
	return &BrainfuckScript{}
}

func (b *BrainfuckScript) SetSource(source []byte) {
	b.content = source
}

func (b *BrainfuckScript) Sanitize() {
	r := regexp.MustCompile(clear_string_regexp)
	b.content = bytes.Join(r.FindAll(b.content, -1), nil)
}

func (b *BrainfuckScript) HasRemaining() bool {
	return b.pos < len(b.content)
}

func (b *BrainfuckScript) NextInstruction(currentByte byte) byte {
	token := b.content[b.pos]

	if token == '[' && currentByte == 0x00 {
		b.getClosingBracket()
	}

	if token == ']' && currentByte != 0x00 {
		b.getOpeningBracket()
	}

	b.pos++
	return token
}

func (b *BrainfuckScript) Execute(interpreter interpreter.Interpreter) {
	b.Sanitize()
	input := bufio.NewReader(os.Stdin)
	for b.HasRemaining() {
		switch b.NextInstruction(interpreter.ReadByte()) {
		case '>':
			interpreter.IncrementAddress()
		case '<':
			interpreter.DecrementAddress()
		case '+':
			interpreter.IncrementByte()
		case '-':
			interpreter.DecrementByte()
		case '.':
			fmt.Printf("%c", interpreter.ReadByte())
		case ',':
			b, err := input.ReadByte()
			if err != nil {
				log.Fatalln(err)
			}
			interpreter.SetByte(b)
		}
	}
}

func (b *BrainfuckScript) getClosingBracket() {
	for depth := 1; depth > 0; {
		b.pos++
		switch b.content[b.pos] {
		case '[':
			depth++
		case ']':
			depth--
		}
	}
}

func (b *BrainfuckScript) getOpeningBracket() {
	for depth := 1; depth > 0; {
		b.pos--
		switch b.content[b.pos] {
		case ']':
			depth++
		case '[':
			depth--
		}
	}
}
