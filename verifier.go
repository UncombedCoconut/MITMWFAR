package main

type dwfa struct {
	states      int
	symbols     int
	startState  wfaState
	transitions map[wfaState]map[symbol]struct {
		wfaState
		weight
	}
}
type wfaState int
type weight int
type specialSets struct {
	nonNegative map[wfaState]struct{}
	nonPositive map[wfaState]struct{}
}

type acceptSet struct{}

type turingMachine struct {
	states      int
	symbols     int
	transitions map[tmState]map[symbol]struct {
		symbol
		direction
		tmState
	}
}
type symbol int
type direction bool

const L direction = true
const LEFT direction = true
const R direction = false
const RIGHT direction = false

func (d direction) String() string {
	if d {
		return "L"
	}
	return "R"
}

type tmState int

const A tmState = 0
const B tmState = 1
const C tmState = 2
const D tmState = 3
const E tmState = 4
const F tmState = 5
const Z tmState = -1

const HALTSTATESTRING = "[HALT]"

func (tms tmState) String() string {
	if tms < 0 {
		return HALTSTATESTRING
	}
	return string(byte(tms) + 'A')
}

func MITMWFARverifier(tm turingMachine, leftWFA, rightWFA dwfa, leftSpecialSets, rightSpecialSets specialSets, acceptSet acceptSet) bool {
	return verifyValidTM(tm) &&
		verifyDeterministicWFA(leftWFA) &&
		verifyDeterministicWFA(rightWFA) &&
		verifySymbolCompatibility(tm, leftWFA, rightWFA) &&
		verifyLeadingBlankInvariant(leftWFA) &&
		verifyLeadingBlankInvariant(rightWFA) &&
		verifySpecialSets(leftWFA, leftSpecialSets) &&
		verifySpecialSets(rightWFA, rightSpecialSets) &&
		verifyStartConfigAccept(acceptSet) &&
		verifyNoHaltingConfigAccepted(tm, acceptSet) &&
		verifyForwardClosed(tm, leftWFA, rightWFA, leftSpecialSets, rightSpecialSets, acceptSet)
}

func verifyValidTM(tm turingMachine) bool {
	if tm.states <= 0 || tm.symbols <= 0 {
		return false
	}
	return true
}

func verifyDeterministicWFA(wfa dwfa) bool {
	if wfa.states <= 0 || wfa.symbols <= 0 {
		return false
	}
	if wfa.startState < 0 || int(wfa.startState) >= wfa.states {
		return false
	}
	for i := 0; i < wfa.states; i++ {
		for j := 0; j < wfa.symbols; j++ {
			transition, ok := wfa.transitions[wfaState(i)][symbol(j)]
			if !ok {
				return false
			}
			if transition.wfaState < 0 || int(transition.wfaState) >= wfa.states {
				return false
			}
		}
	}
	return true
}

func verifySymbolCompatibility(tm turingMachine, leftWFA, rightWFA dwfa) bool {
	return tm.symbols == leftWFA.symbols && tm.symbols == rightWFA.symbols
}

func verifyLeadingBlankInvariant(wfa dwfa) bool {
	state := wfa.startState
	transition := wfa.transitions[state][0]
	return transition.wfaState == state && transition.weight == 0
}

func verifySpecialSets(wfa dwfa, specialSets specialSets) bool {
	for i := 0; i < wfa.states; i++ {
		for j := 0; j < wfa.symbols; j++ {
			transition := wfa.transitions[wfaState(i)][symbol(j)]
			if !transitionRetainsSpecialSets(wfaState(i), transition.wfaState, transition.weight, specialSets) {
				return false
			}
		}
	}
	return true
}

func transitionRetainsSpecialSets(startState, endState wfaState, weight weight, specialSets specialSets) bool {

	_, endNonPositive := specialSets.nonPositive[endState]
	_, startNonPositive := specialSets.nonPositive[startState]
	if endNonPositive {
		if !startNonPositive {
			return false
		}
		if weight > 0 {
			return false
		}
	}

	_, endNonNegative := specialSets.nonNegative[endState]
	_, startNonNegative := specialSets.nonNegative[startState]
	if endNonNegative {
		if !startNonNegative {
			return false
		}
		if weight < 0 {
			return false
		}
	}

	return true
}

func verifyStartConfigAccept(acceptSet acceptSet) bool {
	panic("unimplemented")
}

func verifyNoHaltingConfigAccepted(tm turingMachine, acceptSet acceptSet) bool {
	panic("unimplemented")
}

func verifyForwardClosed(tm turingMachine, leftWFA, rightWFA dwfa, leftSpecialSets, rightSpecialSets specialSets, acceptSet acceptSet) bool {
	panic("unimplemented")
}
