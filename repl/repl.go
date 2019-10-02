package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/shavit/go-interpreter/lexer"
	"github.com/shavit/go-interpreter/token"
)

const PROMPT = ">> "

// Start starts teh repl
// It starts to scan the input line by line to
//   create tokens, then print it out
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)

		for tkn := l.NextToken(); tkn.Type != token.EOF; tkn = l.NextToken() {
			fmt.Printf(`%+v\n`, tkn)
		}
	}
}
