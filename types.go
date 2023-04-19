package main

import "fmt"

type dwfa struct {
	states      int
	symbols     int
	startState  wfaState
	transitions map[wfaState]map[symbol]wfaTransition
}
type wfaState int
type wfaTransition struct {
	wfaState
	weight
}
type weight int

const MAXINT = int(^uint(0) >> 2)
const MININT = -MAXINT - 1

func check(w weight) {
	if int(w) > MAXINT || int(w) < MININT {
		panic("possible integer overflow detected")
	}
}

type specialSets struct {
	nonNegative set[wfaState]
	nonPositive set[wfaState]
}

type turingMachine struct {
	states      int
	symbols     int
	transitions map[tmState]map[symbol]tmTransition
}
type tmTransition struct {
	symbol
	direction
	tmState
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

const TMSTARTSTATE tmState = 0
const TMSTARTSYMBOL symbol = 0

const HALTSTATESTRING = "[HALT]"

func (tms tmState) String() string {
	if tms < 0 {
		return HALTSTATESTRING
	}
	return string(byte(tms) + 'A')
}

type acceptSet map[config]bounds

type config struct {
	tmState    tmState
	tmSymbol   symbol
	leftState  wfaState
	rightState wfaState
}

type bounds map[boundType]weight

type boundType bool

const LOWER boundType = false
const UPPER boundType = true

type configWithWeight struct {
	config
	weight
}

type set[T comparable] map[T]struct{}

func (s set[T]) contains(elem T) bool {
	_, exists := s[elem]
	return exists
}
func (s set[T]) add(elem T) {
	s[elem] = struct{}{}
}
func (s set[T]) remove(elem T) {
	delete(s, elem)
}

func copyWFA(oldWFA dwfa) dwfa {
	newWFA := dwfa{
		states:      oldWFA.states,
		symbols:     oldWFA.symbols,
		startState:  oldWFA.startState,
		transitions: map[wfaState]map[symbol]wfaTransition{},
	}
	for state, tmp := range oldWFA.transitions {
		newWFA.transitions[state] = map[symbol]wfaTransition{}
		for symbol, transition := range tmp {
			newWFA.transitions[state][symbol] = transition
		}
	}
	return newWFA
}

func (wfa dwfa) String() string {
	if wfa.states == 0 {
		return ""
	}
	result := "_"
	for i := 0; i < wfa.states; i++ {
		for j := 0; j < wfa.symbols; j++ {
			transition, ok := wfa.transitions[wfaState(i)][symbol(j)]
			if !ok {
				result += "-,-;"
				continue
			}
			result += fmt.Sprintf("%v,%v;", transition.wfaState, transition.weight)
		}
		result = result[:len(result)-1] + "_"
	}
	return result[1 : len(result)-1]
}

func (tm turingMachine) String() string {
	if tm.states == 0 {
		return ""
	}
	result := ""
	for i := 0; i < tm.states; i++ {
		result += "_"
		for j := 0; j < tm.symbols; j++ {
			transition, ok := tm.transitions[tmState(i)][symbol(j)]
			if !ok {
				result += "---"
				continue
			}
			result += fmt.Sprintf("%v%v%v", transition.symbol, transition.direction, transition.tmState)
		}
	}
	return result[1:]
}

func (s set[T]) String() string {
	if len(s) == 0 {
		return ""
	}
	result := ""
	for elem := range s {
		result += fmt.Sprintf(",%v", elem)
	}
	return result[1:]
}

func (s specialSets) String() string {
	return fmt.Sprintf("%v_%v", s.nonNegative, s.nonPositive)
}

func (as acceptSet) String() string {
	if len(as) == 0 {
		return ""
	}
	result := ""
	for config, bounds := range as {
		result += fmt.Sprintf("_%v", config)
		for _, bound := range []boundType{LOWER, UPPER} {
			if value, ok := bounds[bound]; ok {
				result += fmt.Sprintf(",%v", value)
			} else {
				result += ",-"
			}
		}
	}
	return result[1:]
}

func (c config) String() string {
	return fmt.Sprintf("%v,%v,%v,%v", c.tmState, c.tmSymbol, c.leftState, c.rightState)
}
