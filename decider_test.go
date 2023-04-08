package main

import (
	"reflect"
	"testing"
)

func TestDeriveSpecialSets(t *testing.T) {
	wfa := dwfa{
		states:     4,
		symbols:    2,
		startState: 0,
		transitions: map[wfaState]map[symbol]wfaTransition{
			0: {0: {0, 0},
				1: {1, 0}},
			1: {0: {2, 0},
				1: {1, 0}},
			2: {0: {2, 1},
				1: {3, 0}},
			3: {0: {3, -1},
				1: {3, 0}},
		},
	}
	expectedSets := specialSets{
		nonNegative: map[wfaState]struct{}{0: {}, 1: {}, 2: {}},
		nonPositive: map[wfaState]struct{}{0: {}, 1: {}},
	}

	specialSets := deriveSpecialSets(wfa)
	if !reflect.DeepEqual(expectedSets, specialSets) {
		t.Fail()
	}
}

func TestFindAcceptSet(t *testing.T) {

	tm := turingMachine{
		states:  2,
		symbols: 2,
		transitions: map[tmState]map[symbol]tmTransition{
			A: {0: {1, R, B},
				1: {1, L, A}},
			B: {0: {0, L, A},
				1: {0, R, B}},
		},
	}
	leftWFA := dwfa{
		states:     1,
		symbols:    2,
		startState: 0,
		transitions: map[wfaState]map[symbol]wfaTransition{
			0: {0: {0, 0},
				1: {0, 1}},
		},
	}
	rightWFA := dwfa{
		states:     3,
		symbols:    2,
		startState: 0,
		transitions: map[wfaState]map[symbol]wfaTransition{
			0: {0: {0, 0},
				1: {1, 0}},
			1: {0: {2, 0},
				1: {1, 1}},
			2: {0: {2, 0},
				1: {2, 0}},
		},
	}
	leftSpecialSets := specialSets{
		nonNegative: map[wfaState]struct{}{0: {}},
		nonPositive: map[wfaState]struct{}{},
	}
	rightSpecialSets := specialSets{
		nonNegative: map[wfaState]struct{}{0: {}, 1: {}, 2: {}},
		nonPositive: map[wfaState]struct{}{0: {}},
	}

	expectedResult := acceptSet{
		{A, 0, 0, 0}: {LOWER: 0},
		{A, 1, 0, 0}: {LOWER: 0},
		{A, 0, 0, 1}: {LOWER: 0},
		{A, 1, 0, 1}: {LOWER: 0},
		{B, 0, 0, 0}: {LOWER: 0},
		{B, 1, 0, 0}: {LOWER: 0},
		{B, 1, 0, 1}: {LOWER: 0},
	}
	result := findAcceptSet(tm, leftWFA, rightWFA, leftSpecialSets, rightSpecialSets)

	if !reflect.DeepEqual(expectedResult, result) {
		t.Fail()
	}
}
