package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	input := bufio.NewScanner(os.Stdin)
	result := MITMWFARverifier(parseFullCertificate(input))
	fmt.Println(result)
}
