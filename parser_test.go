package main

import (
	"bufio"
	"os"
	"strings"
	"testing"
)

func TestParser(t *testing.T) {
	inputString := `1RB1LA_0LA0RB
	0,0;0,1
	0,0;1,0_2,0;1,1_2,0;2,0
	0_
	0,1,2_0
	A,0,0,0,0,-_A,1,0,0,0,-_A,0,0,1,0,-_A,1,0,1,0,-_B,0,0,0,0,-_B,1,0,0,0,-_B,1,0,1,0,-`

	input := bufio.NewScanner(strings.NewReader(strings.ReplaceAll(inputString, "\t", "")))
	if !MITMWFARverifier(parseFullCertificate(input)) {
		t.Fail()
	}
}

func TestFullCert(t *testing.T) {
	file, err := os.Open("TestFullCertificate.txt")
	if err != nil {
		t.FailNow()
	}
	defer file.Close()
	input := bufio.NewScanner(file)
	if !MITMWFARverifier(parseFullCertificate(input)) {
		t.Fail()
	}
}

func TestShortCert(t *testing.T) {
	file, err := os.Open("TestShortCertificate.txt")
	if err != nil {
		t.FailNow()
	}
	defer file.Close()
	input := bufio.NewScanner(file)
	if !MITMWFARverifier(parseShortCertificate(input)) {
		t.Fail()
	}
}
