package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/stossmt/golox/lib"
)

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage golox [script]")
		os.Exit(64)
	} else if len(os.Args) == 2 {
		runFile(os.Args[1])
	} else {
		runPrompt()
	}
}

func runPrompt() {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter text: ")
		text, err := reader.ReadString('\n')
		if err != nil {
			panic("Can't read input")
		}
		run(text)
	}
}

func runFile(path string) {
	file, err := os.Open(path)
	if err != nil {
		msg := fmt.Sprintf("Unable to open file: %v", path)
		panic(msg)
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		msg := fmt.Sprintf("Unable to read file: %v", path)
		panic(msg)
	}
	run(string(b))
}

func run(input string) {
	reporter := lib.NewReporter()
	tokens := lib.Scan(input, reporter)
	for _, t := range tokens {
		fmt.Println(t)
	}
	if reporter.HadErr {
		os.Exit(65)
	}
}
