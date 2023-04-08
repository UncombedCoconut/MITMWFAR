package main

import "fmt"

func ShortCertCompletion(tm turingMachine, leftWFA, rightWFA dwfa) (turingMachine, dwfa, dwfa, specialSets, specialSets, acceptSet) {
	leftSpecialSets := deriveSpecialSets(leftWFA)
	rightSpecialSets := deriveSpecialSets(rightWFA)
	acceptSet := findAcceptSet(tm, leftWFA, rightWFA, leftSpecialSets, rightSpecialSets)
	return tm, leftWFA, rightWFA, leftSpecialSets, rightSpecialSets, acceptSet
}

func deriveSpecialSets(wfa dwfa) specialSets {
	possibleNegative := map[wfaState]struct{}{}
	possiblePositive := map[wfaState]struct{}{}
	for _, tmp := range wfa.transitions {
		for _, transition := range tmp {
			if transition.weight < 0 {
				possibleNegative[transition.wfaState] = struct{}{}
			}
			if transition.weight > 0 {
				possiblePositive[transition.wfaState] = struct{}{}
			}
		}
	}
	completeClosure(possibleNegative, wfa)
	completeClosure(possiblePositive, wfa)

	specialSets := specialSets{
		nonNegative: map[wfaState]struct{}{},
		nonPositive: map[wfaState]struct{}{},
	}
	for i := 0; i < wfa.states; i++ {
		if _, ok := possibleNegative[wfaState(i)]; !ok {
			specialSets.nonNegative[wfaState(i)] = struct{}{}
		}
		if _, ok := possiblePositive[wfaState(i)]; !ok {
			specialSets.nonPositive[wfaState(i)] = struct{}{}
		}
	}
	return specialSets
}

func completeClosure(set map[wfaState]struct{}, wfa dwfa) {
	todo := []wfaState{}
	for initialState := range set {
		todo = append(todo, initialState)
	}
	for len(todo) > 0 {
		currentState := todo[0]
		todo = todo[1:]
		for _, transition := range wfa.transitions[currentState] {
			nextState := transition.wfaState
			if _, ok := set[nextState]; !ok {
				set[nextState] = struct{}{}
				todo = append(todo, nextState)
			}
		}
	}
}

func findAcceptSet(tm turingMachine, leftWFA, rightWFA dwfa, leftSpecialSets, rightSpecialSets specialSets) acceptSet {
	initialConfig := config{TMSTARTSTATE, TMSTARTSYMBOL, leftWFA.startState, rightWFA.startState}
	initialBounds := bounds{LOWER: 0, UPPER: 0}
	todo := []config{initialConfig}
	acceptSet := acceptSet{initialConfig: initialBounds}

	for len(todo) > 0 {
		currentConfig := todo[0]
		currentBounds := acceptSet[currentConfig]
		todo = todo[1:]

		for _, nextConfigWithWeightChange := range nextConfigsWithWeightChange(currentConfig, tm, leftWFA, rightWFA) {
			if changeAcceptSetToContainNextConfigWithWeightChange(nextConfigWithWeightChange, currentBounds, leftSpecialSets, rightSpecialSets, acceptSet) {
				todo = append(todo, nextConfigWithWeightChange.config)
			}
		}

	}
	return acceptSet
}

func changeAcceptSetToContainNextConfigWithWeightChange(nextConfigWithWeightChange configWithWeight, bounds bounds, leftSpecialSets, rightSpecialSets specialSets, acceptSet acceptSet) bool {
	nextConfig := nextConfigWithWeightChange.config
	lowerbound, lowerExists := bounds[LOWER]
	upperbound, upperExists := bounds[UPPER]

	//adjust bounds according to the change
	if lowerExists {
		lowerbound += nextConfigWithWeightChange.weight
	}
	if upperExists {
		upperbound += nextConfigWithWeightChange.weight
	}

	hardLower := false
	//adjust bounds according to the special sets
	_, leftStateNonNegative := leftSpecialSets.nonNegative[nextConfig.leftState]
	_, rightStateNonNegative := rightSpecialSets.nonNegative[nextConfig.rightState]
	if leftStateNonNegative && rightStateNonNegative {
		hardLower = true
		if !lowerExists || lowerbound < 0 {
			lowerExists = true
			lowerbound = 0
		}
	}
	hardUpper := false
	_, leftStateNonPositive := leftSpecialSets.nonPositive[nextConfig.leftState]
	_, rightStateNonPositive := rightSpecialSets.nonPositive[nextConfig.rightState]
	if leftStateNonPositive && rightStateNonPositive {
		hardUpper = true
		if !upperExists || upperbound > 0 {
			upperExists = true
			upperbound = 0
		}
	}

	nextBounds := map[boundType]weight{}
	if lowerExists {
		nextBounds[LOWER] = lowerbound
	}
	if upperExists {
		nextBounds[UPPER] = upperbound
	}

	if upperExists && lowerExists && upperbound < lowerbound {
		return false
	}
	return ChangeAcceptSetToCountainConfigBounds(acceptSet, nextConfig, nextBounds, hardLower, hardUpper)
}

func ChangeAcceptSetToCountainConfigBounds(acceptSet acceptSet, nextConfig config, nextBounds map[boundType]weight, hardLower, hardUpper bool) bool {
	acceptBounds, ok := acceptSet[nextConfig]
	if !ok {
		acceptSet[nextConfig] = nextBounds
		return true
	}
	change := false
	acceptedLower, acceptedLowerExists := acceptBounds[LOWER]
	nextLower, nextLowerExists := nextBounds[LOWER]
	if acceptedLowerExists && (!nextLowerExists || acceptedLower > nextLower) {
		delete(acceptSet[nextConfig], LOWER)
		if hardLower {
			acceptSet[nextConfig][LOWER] = 0
		}
		change = true
	}

	acceptedUpper, acceptedUpperExists := acceptBounds[UPPER]
	nextUpper, nextUpperExists := nextBounds[UPPER]
	if acceptedUpperExists && (!nextUpperExists || acceptedUpper < nextUpper) {
		delete(acceptSet[nextConfig], UPPER)
		if hardUpper {
			acceptSet[nextConfig][UPPER] = 0
		}
		change = true
	}
	return change
}

//------------------------------------------------------------------------------------------------

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
func MITMWFARdecider(tm turingMachine, maxTransitions, maxStatesLeft, maxStatesRight, maxWeightPairs int) bool {
	leftWFA := dwfa{
		states:      2,
		symbols:     tm.symbols,
		startState:  0,
		transitions: map[wfaState]map[symbol]wfaTransition{},
	}
	rightWFA := dwfa{
		states:      2,
		symbols:     tm.symbols,
		startState:  0,
		transitions: map[wfaState]map[symbol]wfaTransition{},
	}
	leftWFA.transitions[0] = map[symbol]wfaTransition{}
	leftWFA.transitions[1] = map[symbol]wfaTransition{}
	rightWFA.transitions[0] = map[symbol]wfaTransition{}
	rightWFA.transitions[1] = map[symbol]wfaTransition{}
	for i := 0; i < tm.symbols; i++ {
		//1 is deadstate. Transitions to 1 are default and don't count towards currentTransitions
		leftWFA.transitions[0][symbol(i)] = wfaTransition{1, 0}
		leftWFA.transitions[1][symbol(i)] = wfaTransition{1, 0}
		rightWFA.transitions[0][symbol(i)] = wfaTransition{1, 0}
		rightWFA.transitions[1][symbol(i)] = wfaTransition{1, 0}
	}
	leftWFA.transitions[0][0] = wfaTransition{0, 0}
	rightWFA.transitions[0][0] = wfaTransition{0, 0}
	return recursiveDecider(tm, leftWFA, rightWFA, 2, 2, 2, maxTransitions, maxStatesLeft, maxStatesRight, maxWeightPairs)
}

func recursiveDecider(tm turingMachine, leftWFA, rightWFA dwfa, currentTransitions, currentStatesLeft, currentStatesRight, maxTransitions, maxStatesLeft, maxStatesRight, maxWeightPairs int) bool {
	closed, breakingSide, breakingState, breakingSymbol := findClosure(tm, leftWFA, rightWFA)
	if closed {
		return recursiveWeightAdder(tm, leftWFA, rightWFA, 0, maxWeightPairs)
	}
	if currentTransitions >= maxTransitions {
		return false
	}
	switch breakingSide {
	case LEFT:
		if currentStatesLeft < maxStatesLeft {
			newWFA := copyWFA(leftWFA)
			newState := wfaState(newWFA.states)
			newWFA.states += 1
			newWFA.transitions[newState] = map[symbol]wfaTransition{}
			for i := 0; i < newWFA.symbols; i++ {
				newWFA.transitions[newState][symbol(i)] = wfaTransition{1, 0}
			}
			newWFA.transitions[breakingState][breakingSymbol] = wfaTransition{newState, 0}
			if recursiveDecider(tm, newWFA, rightWFA, currentTransitions+1, currentStatesLeft+1, currentStatesRight, maxTransitions, maxStatesLeft, maxStatesRight, maxWeightPairs) {
				return true
			}
		}
		for i := 0; i < leftWFA.states; i++ {
			if i == 1 {
				continue
			}
			newWFA := copyWFA(leftWFA)
			newWFA.transitions[breakingState][breakingSymbol] = wfaTransition{wfaState(i), 0}
			if recursiveDecider(tm, newWFA, rightWFA, currentTransitions+1, currentStatesLeft, currentStatesRight, maxTransitions, maxStatesLeft, maxStatesRight, maxWeightPairs) {
				return true
			}
		}
	case RIGHT:
		if currentStatesRight < maxStatesRight {
			newWFA := copyWFA(rightWFA)
			newState := wfaState(newWFA.states)
			newWFA.states += 1
			newWFA.transitions[newState] = map[symbol]wfaTransition{}
			for i := 0; i < newWFA.symbols; i++ {
				newWFA.transitions[newState][symbol(i)] = wfaTransition{1, 0}
			}
			newWFA.transitions[breakingState][breakingSymbol] = wfaTransition{newState, 0}
			if recursiveDecider(tm, leftWFA, newWFA, currentTransitions+1, currentStatesLeft, currentStatesRight+1, maxTransitions, maxStatesLeft, maxStatesRight, maxWeightPairs) {
				return true
			}
		}
		for i := 0; i < rightWFA.states; i++ {
			if i == 1 {
				continue
			}
			newWFA := copyWFA(rightWFA)
			newWFA.transitions[breakingState][breakingSymbol] = wfaTransition{wfaState(i), 0}
			if recursiveDecider(tm, leftWFA, newWFA, currentTransitions+1, currentStatesLeft, currentStatesRight, maxTransitions, maxStatesLeft, maxStatesRight, maxWeightPairs) {
				return true
			}
		}
	}
	return false
}

func findClosure(tm turingMachine, leftWFA, rightWFA dwfa) (bool, direction, wfaState, symbol) {
	accept := set[config]{}
	initialConfig := config{TMSTARTSTATE, TMSTARTSYMBOL, leftWFA.startState, rightWFA.startState}
	accept.add(initialConfig)
	todo := []config{initialConfig}
	for len(todo) > 0 {
		currentConfig := todo[0]
		todo = todo[1:]
		for _, tmp := range nextConfigsWithWeightChange(currentConfig, tm, leftWFA, rightWFA) {
			nextConfig := tmp.config
			if accept.contains(nextConfig) {
				continue
			}

			if nextConfig.leftState == wfaState(1) {
				return false, LEFT, currentConfig.leftState, tm.transitions[currentConfig.tmState][currentConfig.tmSymbol].symbol
			}
			if nextConfig.rightState == wfaState(1) {
				return false, RIGHT, currentConfig.rightState, tm.transitions[currentConfig.tmState][currentConfig.tmSymbol].symbol
			}
			accept.add(nextConfig)
			todo = append(todo, nextConfig)
		}
	}
	return true, L, 0, 0
}

func recursiveWeightAdder(tm turingMachine, leftWFA, rightWFA dwfa, currenWeightPairs, maxWeightPairs int) bool {
	if MITMWFARverifier(ShortCertCompletion(tm, leftWFA, rightWFA)) {
		fmt.Println(leftWFA, rightWFA)
		return true
	}
	if currenWeightPairs >= maxWeightPairs {
		return false
	}
	weightPermutations := [][2]weight{{1, -1}}
	if currenWeightPairs > 0 {
		weightPermutations = append(weightPermutations, [2]weight{-1, 1})
	}
	for _, weights := range weightPermutations {
		for leftState, tmpLeft := range leftWFA.transitions {
			for leftSymbol, leftTransition := range tmpLeft {
				if leftTransition.wfaState == 1 {
					continue
				}
				newLeftWFA := copyWFA(leftWFA)
				newLeftWFA.transitions[leftState][leftSymbol] = wfaTransition{leftTransition.wfaState, leftTransition.weight + weights[0]}
				for rightState, tmpRight := range rightWFA.transitions {
					for rightSymbol, rightTransition := range tmpRight {
						if rightTransition.wfaState == 1 {
							continue
						}
						newRightWFA := copyWFA(rightWFA)
						newRightWFA.transitions[rightState][rightSymbol] = wfaTransition{rightTransition.wfaState, rightTransition.weight + weights[1]}
						if recursiveWeightAdder(tm, newLeftWFA, newRightWFA, currenWeightPairs+1, maxWeightPairs) {
							return true
						}
					}
				}
			}
		}
	}
	return false
}

func MITMDFAdecider(tm turingMachine, maxStates int) bool {
	return MITMWFARdecider(tm, maxStates*tm.symbols, maxStates, maxStates, 0)
}

func (wfa dwfa) String() string {
	result := "_"
	for i := 0; i < wfa.states; i++ {
		for j := 0; j < wfa.symbols; j++ {
			transition, ok := wfa.transitions[wfaState(i)][symbol(j)]
			if !ok {
				result += "-,-;"
				continue
			}
			result += fmt.Sprintf("%d,%d;", transition.wfaState, transition.weight)
		}
		result = result[:len(result)-1] + "_"
	}
	return result[1 : len(result)-1]
}
