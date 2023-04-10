package main

import (
	"bufio"
	"strconv"
	"strings"
)

//parsing isn't robust. Might panic on bad input.
func parseFullCertificate(input *bufio.Scanner, workTokens chan struct{}, printMode int) {
	for input.Scan() {
		_ = <-workTokens
		tm := parseTM(input.Text())
		input.Scan()
		leftWFA := parseWFA(input.Text())
		input.Scan()
		rightWFA := parseWFA(input.Text())
		input.Scan()
		leftSpecialSets := parseSpecialSets(input.Text())
		input.Scan()
		rightSpecialSets := parseSpecialSets(input.Text())
		input.Scan()
		acceptSet := parseAcceptSet(input.Text())
		go func() {
			MITMWFARverifier(tm, leftWFA, rightWFA, leftSpecialSets, rightSpecialSets, acceptSet, printMode)
			workTokens <- struct{}{}
		}()
	}
}

//parsing isn't robust. Might panic on bad input.
func parseShortCertificate(input *bufio.Scanner, workTokens chan struct{}, printMode int) {
	for input.Scan() {
		_ = <-workTokens
		tm := parseTM(input.Text())
		input.Scan()
		leftWFA := parseWFA(input.Text())
		input.Scan()
		rightWFA := parseWFA(input.Text())
		leftSpecialSets := deriveSpecialSets(leftWFA)
		rightSpecialSets := deriveSpecialSets(rightWFA)
		acceptSet := findAcceptSet(tm, leftWFA, rightWFA, leftSpecialSets, rightSpecialSets)
		go func() {
			MITMWFARverifier(tm, leftWFA, rightWFA, leftSpecialSets, rightSpecialSets, acceptSet, printMode)
			workTokens <- struct{}{}
		}()
	}
}

//parsing isn't robust. Might panic on bad input.
func runSpecificValues(input *bufio.Scanner, workTokens chan struct{}, printMode, maxTransitions, maxLeftStates, maxRightStates, maxWeightPairs int) {
	for input.Scan() {
		_ = <-workTokens
		tm := parseTM(input.Text())
		go func() {
			MITMWFARdecider(tm, maxTransitions, maxLeftStates, maxRightStates, maxWeightPairs, printMode)
			workTokens <- struct{}{}
		}()
	}
}

//parsing isn't robust. Might panic on bad input.
func runWeightedScan(input *bufio.Scanner, workTokens chan struct{}, printMode, maxTransitions, maxWeightPairs int) {
	for input.Scan() {
		_ = <-workTokens
		tm := parseTM(input.Text())
		go func() {
			for transitions := 2; transitions <= maxTransitions; transitions++ {
				if MITMWFARdecider(tm, transitions, maxTransitions, maxTransitions, maxWeightPairs, printMode) {
					break
				}
			}
			workTokens <- struct{}{}
		}()
	}
}

//parsing isn't robust. Might panic on bad input.
func runDFAScan(input *bufio.Scanner, workTokens chan struct{}, printMode, maxStates int) {
	for input.Scan() {
		_ = <-workTokens
		tm := parseTM(input.Text())
		go func() {
			maxTransitions := 4*maxStates - 4
			for transitions := 2; transitions <= maxTransitions; transitions++ {
				if MITMWFARdecider(tm, transitions, maxStates, maxStates, 0, printMode) {
					break
				}
			}
			workTokens <- struct{}{}
		}()
	}
}

//standard text format
func parseTM(s string) turingMachine {
	stateStrings := strings.Split(s, "_")

	tm := turingMachine{
		states:      len(stateStrings),
		symbols:     len(stateStrings[0]) / 3,
		transitions: map[tmState]map[symbol]tmTransition{},
	}
	for i, stateString := range stateStrings {
		tm.transitions[tmState(i)] = map[symbol]tmTransition{}
		for j := 0; len(stateString) >= 3; j++ {
			symbolString := stateString[:3]
			stateString = stateString[3:]
			newTMState := tmState(symbolString[2] - 'A')
			if int(newTMState) < 0 || int(newTMState) >= tm.states {
				continue
			}
			newSymbol := symbol(symbolString[0] - '0')
			newDirection := L
			if symbolString[1] == 'R' {
				newDirection = R
			}
			tm.transitions[tmState(i)][symbol(j)] = tmTransition{newSymbol, newDirection, newTMState}
		}
	}
	return tm
}

//"0,0;1,0_1,1;0,0"
func parseWFA(s string) dwfa {
	stateStrings := strings.Split(s, "_")
	wfa := dwfa{
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
	return wfa
}

//"0,1,4,5_0,2"
func parseSpecialSets(s string) specialSets {
	setStrings := strings.Split(s, "_")
	return specialSets{
		nonNegative: parseStateSet(setStrings[0]),
		nonPositive: parseStateSet(setStrings[1]),
	}
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
func parseAcceptSet(s string) acceptSet {
	acceptSet := acceptSet{}
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
		acceptSet[newConfig] = newBounds
	}
	return acceptSet
}
