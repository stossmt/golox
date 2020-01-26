package lib

import (
	"fmt"
	"os"
)

type Reporter struct {
	HadErr bool
}

func NewReporter() *Reporter {
	return &Reporter{}
}

func (r *Reporter) Report(line int, where string, msg string) {
	s := "[line %v] Error %v: message %v"
	fmt.Fprintf(os.Stderr, s, line, where, msg)
	r.HadErr = true
}
