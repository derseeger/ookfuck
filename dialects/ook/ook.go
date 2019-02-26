package ook

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/derseeger/ookfuck/interpreter"
	"log"
	"os"
	"regexp"
	"strings"
)

const clear_string_regexp = `[\.\?\!]*`

var translation map[string]string

func init() {
	translation = make(map[string]string)
	translation[".."] = "+"
	translation["!!"] = "-"
	translation[".?"] = ">"
	translation["?."] = "<"
	translation["!?"] = "["
	translation["?!"] = "]"
	translation["!."] = "."
	translation[".!"] = ","
}

type OokScript struct {
	content []byte
	pos     int
}

func NewOokScript() *OokScript {
	return &OokScript{}
}

func (b *OokScript) SetSource(source []byte) {
	b.content = source
}

func (b *OokScript) Sanitize() {
	// Remove everything but the ook symbols (.!?)
	r := regexp.MustCompile(clear_string_regexp)
	b.content = bytes.Join(r.FindAll(b.content, -1), nil)
	ook := splitSubN(string(b.content), 2)
	bf := []string{}
	for _, o := range ook {
		for key, val := range translation {
			o = strings.Replace(o, key, val, -1)
		}
		bf = append(bf, o)
	}
	b.content = []byte(strings.Join(bf, ""))

}

func splitSubN(s string, n int) []string {
	sub := ""
	subs := []string{}

	runes := bytes.Runes([]byte(s))
	l := len(runes)
	for i, r := range runes {
		sub = sub + string(r)
		if (i+1)%n == 0 {
			subs = append(subs, sub)
			sub = ""
		} else if (i + 1) == l {
			subs = append(subs, sub)
		}
	}

	return subs
}

func (b *OokScript) HasRemaining() bool {
	return b.pos < len(b.content)
}

func (b *OokScript) NextInstruction(currentByte byte) byte {
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

func (b *OokScript) Execute(interpreter interpreter.Interpreter) {
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

func (b *OokScript) getClosingBracket() {
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

func (b *OokScript) getOpeningBracket() {
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
