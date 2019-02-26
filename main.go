package main

import (
	"fmt"
	"github.com/derseeger/ookfuck/dialects/brainfuck"
	"github.com/derseeger/ookfuck/dialects/ook"
	inter "github.com/derseeger/ookfuck/interpreter"
	"io/ioutil"
	"os"
)

func main() {
	exec := os.Args[0:1]
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Printf("Missing arguments!\nFormat is: %v [FILE] [--FLAGS]\n", exec[0])
	} else {
		file := args[0]
		if _, err := os.Stat(file); err == nil {
			// Read file
			handle, _ := os.Open(file)
			source_bytes, _ := ioutil.ReadAll(handle)

			interpreter := inter.NewEsotericInterpreter()

			var script inter.Script

			if len(args) > 1 && args[1] == "ook" {
				script = ook.NewOokScript()
			} else {
				script = brainfuck.NewBrainfuckScript()
			}
			script.SetSource(source_bytes)
			script.Execute(interpreter)

		} else if os.IsNotExist(err) {
			fmt.Printf("The file %v doesn't exist!\n", file)
		} else {
			fmt.Printf("There was an unexpected error:\n%v\n", err)
		}
	}
}
