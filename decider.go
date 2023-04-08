package main

func shortCertCompletion(tm turingMachine, leftWFA, rightWFA dwfa) (turingMachine, dwfa, dwfa, specialSets, specialSets, acceptSet) {
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
