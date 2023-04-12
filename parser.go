package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func parseFullCertificate(input *bufio.Scanner, workTokens chan struct{}, printMode int) {
	for input.Scan() {
		tm, err := parseTM(input.Text())
		if err != nil {
			if input.Text() != "" {
				fmt.Fprintln(os.Stderr, err)
			}
			continue
		}
		input.Scan()
		leftWFA, err := parseWFA(input.Text())
		if err != nil {
			if input.Text() != "" {
				fmt.Fprintln(os.Stderr, err)
			}
			continue
		}
		input.Scan()
		rightWFA, err := parseWFA(input.Text())
		if err != nil {
			if input.Text() != "" {
				fmt.Fprintln(os.Stderr, err)
			}
			continue
		}
		input.Scan()
		leftSpecialSets, err := parseSpecialSets(input.Text())
		if err != nil {
			if input.Text() != "" {
				fmt.Fprintln(os.Stderr, err)
			}
			continue
		}
		input.Scan()
		rightSpecialSets, err := parseSpecialSets(input.Text())
		if err != nil {
			if input.Text() != "" {
				fmt.Fprintln(os.Stderr, err)
			}
			continue
		}
		input.Scan()
		acceptSet, err := parseAcceptSet(input.Text())
		if err != nil {
			if input.Text() != "" {
				fmt.Fprintln(os.Stderr, err)
			}
			continue
		}
		_ = <-workTokens
		go func() {
			MITMWFARverifier(tm, leftWFA, rightWFA, leftSpecialSets, rightSpecialSets, acceptSet, printMode)
			workTokens <- struct{}{}
		}()
	}
}

func parseShortCertificate(input *bufio.Scanner, workTokens chan struct{}, printMode int) {
	for input.Scan() {
		tm, err := parseTM(input.Text())
		if err != nil {
			if input.Text() != "" {
				fmt.Fprintln(os.Stderr, err)
			}
			continue
		}
		input.Scan()
		leftWFA, err := parseWFA(input.Text())
		if err != nil {
			if input.Text() != "" {
				fmt.Fprintln(os.Stderr, err)
			}
			continue
		}
		input.Scan()
		rightWFA, err := parseWFA(input.Text())
		if err != nil {
			if input.Text() != "" {
				fmt.Fprintln(os.Stderr, err)
			}
			continue
		}
		leftSpecialSets := deriveSpecialSets(leftWFA)
		rightSpecialSets := deriveSpecialSets(rightWFA)
		acceptSet := findAcceptSet(tm, leftWFA, rightWFA, leftSpecialSets, rightSpecialSets)
		_ = <-workTokens
		go func() {
			MITMWFARverifier(tm, leftWFA, rightWFA, leftSpecialSets, rightSpecialSets, acceptSet, printMode)
			workTokens <- struct{}{}
		}()
	}
}

func runSpecificValues(input *bufio.Scanner, workTokens chan struct{}, printMode, maxTransitions, maxLeftStates, maxRightStates, maxWeightPairs, addedMemory int) {
	for input.Scan() {
		tm, err := parseTM(input.Text())
		if err != nil {
			if input.Text() != "" {
				fmt.Fprintln(os.Stderr, err)
			}
			continue
		}
		_ = <-workTokens
		go func() {
			MITMWFARdecider(tm, maxTransitions, maxLeftStates, maxRightStates, maxWeightPairs, addedMemory, printMode)
			workTokens <- struct{}{}
		}()
	}
}

func runWeightedScan(input *bufio.Scanner, workTokens chan struct{}, printMode, maxTransitions, maxWeightPairs, addedMemory int) {
	for input.Scan() {
		tm, err := parseTM(input.Text())
		if err != nil {
			if input.Text() != "" {
				fmt.Fprintln(os.Stderr, err)
			}
			continue
		}
		_ = <-workTokens
		go func() {
			for transitions := 2; transitions <= maxTransitions; transitions++ {
				if MITMWFARdecider(tm, transitions, maxTransitions, maxTransitions, maxWeightPairs, addedMemory, printMode) {
					break
				}
			}
			workTokens <- struct{}{}
		}()
	}
}

func runDFAScan(input *bufio.Scanner, workTokens chan struct{}, printMode, maxStates int) {
	for input.Scan() {
		tm, err := parseTM(input.Text())
		if err != nil {
			if input.Text() != "" {
				fmt.Fprintln(os.Stderr, err)
			}
			continue
		}
		_ = <-workTokens
		go func() {
			maxTransitions := tm.symbols * (maxStates - 1) * 2
			for transitions := 2; transitions <= maxTransitions; transitions++ {
				if MITMWFARdecider(tm, transitions, maxStates, maxStates, 0, 0, printMode) {
					break
				}
			}
			workTokens <- struct{}{}
		}()
	}
}

type errorString string

func (e errorString) Error() string {
	return string(e)
}

//standard text format
func parseTM(s string) (tm turingMachine, err error) {
	defer func() {
		if recover() != nil {
			err = errorString("Couldn't parse TM: \"" + s + "\"")
		}
	}()

	stateStrings := strings.Split(s, "_")
	if len(stateStrings[0])%3 != 0 {
		panic("")
	}
	tm = turingMachine{
		states:      len(stateStrings),
		symbols:     len(stateStrings[0]) / 3,
		transitions: map[tmState]map[symbol]tmTransition{},
	}
	if tm.states < 2 {
		panic("")
	}
	for i, stateString := range stateStrings {
		if len(stateString) != tm.symbols*3 {
			panic("")
		}
		tm.transitions[tmState(i)] = map[symbol]tmTransition{}
		for j := 0; len(stateString) >= 3; j++ {
			symbolString := stateString[:3]
			stateString = stateString[3:]
			newTMState := tmState(symbolString[2] - 'A')
			if int(newTMState) < 0 || int(newTMState) >= tm.states {
				continue
			}
			newSymbol := symbol(symbolString[0] - '0')
			if newSymbol < 0 || int(newSymbol) >= tm.symbols {
				panic("")
			}
			newDirection := L
			if symbolString[1] == 'R' {
				newDirection = R
			}
			tm.transitions[tmState(i)][symbol(j)] = tmTransition{newSymbol, newDirection, newTMState}
		}
	}
	return
}

//"0,0;1,0_1,1;0,0"
func parseWFA(s string) (wfa dwfa, err error) {
	defer func() {
		if recover() != nil {
			err = errorString("Couldn't parse WFA: \"" + s + "\"")
		}
	}()
	stateStrings := strings.Split(s, "_")
	wfa = dwfa{
		states:      len(stateStrings),
		startState:  0,
		transitions: map[wfaState]map[symbol]wfaTransition{},
	}
	for i, stateString := range stateStrings {
		symbolStrings := strings.Split(stateString, ";")
		wfa.symbols = len(symbolStrings)
		wfa.transitions[wfaState(i)] = map[symbol]wfaTransition{}
		for j, symbolString := range symbolStrings {
			values := strings.Split(symbolString, ",")
			targetState, _ := strconv.Atoi(values[0])
			addedWeight, _ := strconv.Atoi(values[1])
			wfa.transitions[wfaState(i)][symbol(j)] = wfaTransition{
				wfaState(targetState),
				weight(addedWeight),
			}
		}
	}
	return
}

//"0,1,4,5_0,2"
func parseSpecialSets(s string) (sets specialSets, err error) {
	defer func() {
		if recover() != nil {
			err = errorString("Couldn't parse special sets: \"" + s + "\"")
		}
	}()
	setStrings := strings.Split(s, "_")
	sets = specialSets{
		nonNegative: parseStateSet(setStrings[0]),
		nonPositive: parseStateSet(setStrings[1]),
	}
	return
}

func parseStateSet(s string) set[wfaState] {
	set := set[wfaState]{}
	if s == "" {
		return set
	}
	for _, stateString := range strings.Split(s, ",") {
		state, _ := strconv.Atoi(stateString)
		set.add(wfaState(state))
	}
	return set
}

//"A,0,0,0,-,-_B,1,0,2,2,-"
func parseAcceptSet(s string) (set acceptSet, err error) {
	defer func() {
		if recover() != nil {
			err = errorString("Couldn't parse accept set: \"" + s + "\"")
		}
	}()
	set = acceptSet{}
	for _, accepter := range strings.Split(s, "_") {
		values := strings.Split(accepter, ",")
		newTMState := tmState(values[0][0] - 'A')
		newSymbol := symbol(values[1][0] - '0')
		leftState, _ := strconv.Atoi(values[2])
		rightState, _ := strconv.Atoi(values[3])
		newConfig := config{newTMState, newSymbol, wfaState(leftState), wfaState(rightState)}
		newBounds := map[boundType]weight{}
		lowerbound, lowerExists := strconv.Atoi(values[4])
		if lowerExists == nil {
			newBounds[LOWER] = weight(lowerbound)
		}
		upperbound, upperExists := strconv.Atoi(values[5])
		if upperExists == nil {
			newBounds[UPPER] = weight(upperbound)
		}
		set[newConfig] = newBounds
	}
	return
}
