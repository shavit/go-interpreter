package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/shavit/go-interpreter/repl"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	user, err := user.Current()
	checkErr(err)

	fmt.Println("   ____ PROGRAMMING LANGUAGE\n\n", user.Username, ", press Ctrl+C to exit")
	repl.Start(os.Stdin, os.Stdout)
}
