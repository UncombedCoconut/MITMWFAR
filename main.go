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

	flag.Parse()

	input := bufio.NewScanner(os.Stdin)
	switch {
	case *fullcert:
		result := MITMWFARverifier(parseFullCertificate(input))
		fmt.Println(result)
	case *shortcert:
		result := MITMWFARverifier(parseShortCertificate(input))
		fmt.Println(result)
	}
}
