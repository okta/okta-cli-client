package iostream

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/mattn/go-isatty"
)

var (
	Input    = os.Stdin
	Output   = os.Stdout
	Messages = os.Stderr
)

func IsInputTerminal() bool {
	return isatty.IsTerminal(Input.Fd())
}

func IsOutputTerminal() bool {
	return isatty.IsTerminal(Output.Fd())
}

func PipedInput() []byte {
	if !IsInputTerminal() {
		reader := bufio.NewReader(Input)
		var pipedInput []byte

		for {
			input, err := reader.ReadBytes('\n')
			if err == io.EOF {
				break
			} else if err != nil {
				panic(fmt.Errorf("unable to read from pipe %v", err))
			}
			pipedInput = append(pipedInput, input...)
		}

		return pipedInput
	}

	return []byte{}
}
