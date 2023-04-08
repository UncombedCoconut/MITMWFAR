package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

func main() {

	fullcert := flag.Bool("fc", false, "reads the full certificate for TMs from stdin")
	shortcert := flag.Bool("sc", false, "reads a short certificate for TMs from stdin")
	transitions := flag.Int("t", 12, "maximum number of non-dead transitions in the combined WFAs")
	leftStates := flag.Int("l", 4, "maximum number of states in the left WFA")
	rightStates := flag.Int("r", 4, "maximum number of states in the left WFA")
	weightPairs := flag.Int("w", 1, "maximum number of weighted transitions in each WFA")

	flag.Parse()

	input := bufio.NewScanner(os.Stdin)
	switch {
	case *fullcert:
		result := MITMWFARverifier(parseFullCertificate(input))
		fmt.Println(result)
	case *shortcert:
		result := MITMWFARverifier(parseShortCertificate(input))
		fmt.Println(result)
	default:
		for input.Scan() {
			result := MITMWFARdecider(parseTM(input.Text()), *transitions, *leftStates, *rightStates, *weightPairs)
			fmt.Println(input.Text(), result)
		}
	}
}
