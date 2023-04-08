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

func TestFindClosure(t *testing.T) {
	t.Run("Incomplete", func(t *testing.T) {
		tm := turingMachine{
			states:  2,
			symbols: 2,
			transitions: map[tmState]map[symbol]tmTransition{
				A: {0: {1, R, B},
					1: {1, L, B}},
				B: {0: {1, L, A},
					1: {1, R, Z}},
			},
		}
		leftWFA := dwfa{
			states:     3,
			symbols:    2,
			startState: 0,
			transitions: map[wfaState]map[symbol]wfaTransition{
				0: {0: {0, 0},
					1: {2, 0}},
				1: {0: {1, 0},
					1: {1, 0}},
				2: {0: {1, 0},
					1: {2, 0}}},
		}
		rightWFA := dwfa{
			states:     2,
			symbols:    2,
			startState: 0,
			transitions: map[wfaState]map[symbol]wfaTransition{
				0: {0: {0, 0},
					1: {1, 0}},
				1: {0: {1, 0},
					1: {1, 0}}},
		}
		result, dir, state, symbol := findClosure(tm, leftWFA, rightWFA)
		if result != false || dir != RIGHT || state != 0 || symbol != 1 {
			t.Fail()
		}
	})
	t.Run("Closed", func(t *testing.T) {
		tm := turingMachine{
			states:  2,
			symbols: 2,
			transitions: map[tmState]map[symbol]tmTransition{
				A: {0: {1, R, B},
					1: {1, L, B}},
				B: {0: {1, L, A},
					1: {1, R, Z}},
			},
		}
		leftWFA := dwfa{
			states:     2,
			symbols:    2,
			startState: 0,
			transitions: map[wfaState]map[symbol]wfaTransition{
				0: {0: {0, 0},
					1: {0, 0}},
				1: {0: {1, 0},
					1: {1, 0}},
				2: {0: {1, 0},
					1: {2, 0}}},
		}
		rightWFA := dwfa{
			states:     3,
			symbols:    2,
			startState: 0,
			transitions: map[wfaState]map[symbol]wfaTransition{
				0: {0: {0, 0},
					1: {2, 0}},
				1: {0: {1, 0},
					1: {1, 0}},
				2: {0: {1, 0},
					1: {2, 0}}},
		}
		result, _, _, _ := findClosure(tm, leftWFA, rightWFA)
		if !result {
			t.Fail()
		}
	})
}

func TestMITMWFARdecider(t *testing.T) {
	tm := turingMachine{
		states:  5,
		symbols: 2,
		transitions: map[tmState]map[symbol]tmTransition{
			A: {0: {1, R, B}},
			B: {0: {0, R, C},
				1: {1, R, C}},
			C: {0: {1, R, D},
				1: {1, R, B}},
			D: {0: {1, L, E},
				1: {1, L, D}},
			E: {0: {0, R, A},
				1: {0, L, E}},
		},
	}
	if !MITMWFARdecider(tm, 11, 4, 4, 1) {
		t.Fail()
	}
}

func TestMITMDFAdecider(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		tm := turingMachine{
			states:  5,
			symbols: 2,
			transitions: map[tmState]map[symbol]tmTransition{
				A: {0: {1, R, B},
					1: {0, L, D}},
				B: {0: {1, R, C}},
				C: {0: {1, L, D},
					1: {0, R, E}},
				D: {0: {1, L, E},
					1: {1, L, A}},
				E: {0: {0, L, A},
					1: {0, L, A}},
			},
		}
		if !MITMDFAdecider(tm, 5) {
			t.Fail()
		}
	})
	t.Run("Failure", func(t *testing.T) {
		tm := turingMachine{
			states:  5,
			symbols: 2,
			transitions: map[tmState]map[symbol]tmTransition{
				A: {0: {1, R, B},
					1: {1, R, E}},
				B: {0: {1, L, C},
					1: {1, R, B}},
				C: {0: {0, R, A},
					1: {0, L, D}},
				D: {0: {1, L, B},
					1: {1, L, D}},
				E: {1: {0, R, A}},
			},
		}
		if MITMDFAdecider(tm, 5) {
			t.Fail()
		}
	})
}
