package main

import (
	"bufio"
	"flag"
	"os"
	"runtime"
)

func main() {

	fullcert := flag.Bool("fc", false, "reads the full certificate for TMs from stdin")
	shortcert := flag.Bool("sc", false, "reads a short certificate for TMs from stdin")
	transitions := flag.Int("t", 12, "maximum number of non-dead transitions in the combined WFAs")
	leftStates := flag.Int("l", 4, "maximum number of states in the left WFA")
	rightStates := flag.Int("r", 4, "maximum number of states in the left WFA")
	weightPairs := flag.Int("w", 1, "maximum number of weighted transitions in each WFA")
	printMode := flag.Int("pm", 0, "what to print: 0 -> solved TMs, 1 -> short certificates, 2 -> full certificates")
	cores := flag.Int("cores", 0, "maximum number of TMs to work on in parallel")

	flag.Parse()

	if *cores <= 0 {
		*cores = runtime.GOMAXPROCS(0)
	}
	workTokens := make(chan struct{}, *cores)
	for i := 0; i < *cores; i++ {
		workTokens <- struct{}{}
	}
	input := bufio.NewScanner(os.Stdin)
	switch {
	case *fullcert:
		parseFullCertificate(input, workTokens, *printMode)
	case *shortcert:
		parseShortCertificate(input, workTokens, *printMode)
	default:
		parseTmStandardFormat(input, workTokens, *printMode, *transitions, *leftStates, *rightStates, *weightPairs)
	}

	//make sure all the work is finished
	for i := 0; i < *cores; i++ {
		_ = <-workTokens
	}
}
