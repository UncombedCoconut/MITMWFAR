package main

import "fmt"

func MITMWFARverifier(tm turingMachine, leftWFA, rightWFA dwfa, leftSpecialSets, rightSpecialSets specialSets, acceptSet acceptSet, printMode int) bool {
	result := verifyCoherentDefinitions(tm, leftWFA, rightWFA, leftSpecialSets, rightSpecialSets, acceptSet) &&
		verifyLeadingBlankInvariant(leftWFA) &&
		verifyLeadingBlankInvariant(rightWFA) &&
		verifySpecialSetsHaveClaimedProperty(leftWFA, leftSpecialSets) &&
		verifySpecialSetsHaveClaimedProperty(rightWFA, rightSpecialSets) &&
		verifyStartConfigAccept(leftWFA, rightWFA, acceptSet) &&
		verifyNoHaltingConfigAccepted(tm, acceptSet) &&
		verifyForwardClosed(tm, leftWFA, rightWFA, leftSpecialSets, rightSpecialSets, acceptSet)
	if result {
		if printMode >= 0 {
			cert := fmt.Sprintln(tm)
			if printMode >= 1 {
				cert += fmt.Sprintln(leftWFA)
				cert += fmt.Sprintln(rightWFA)
			}
			if printMode >= 2 {
				cert += fmt.Sprintln(leftSpecialSets)
				cert += fmt.Sprintln(rightSpecialSets)
				cert += fmt.Sprintln(acceptSet)
			}
			fmt.Print(cert)
		}
	}
	return result
}

func verifyCoherentDefinitions(tm turingMachine, leftWFA, rightWFA dwfa, leftSpecialSets, rightSpecialSets specialSets, acceptSet acceptSet) bool {
	return verifyValidTM(tm) &&
		verifyDeterministicWFA(leftWFA) &&
		verifyDeterministicWFA(rightWFA) &&
		verifySymbolCompatibility(tm, leftWFA, rightWFA) &&
		verifySpecialSetsAreSubsets(leftWFA, leftSpecialSets) &&
		verifySpecialSetsAreSubsets(rightWFA, rightSpecialSets) &&
		verifyAcceptSetIsValid(tm, leftWFA, rightWFA, acceptSet)
}

func verifyValidTM(tm turingMachine) bool {
	if tm.states <= 0 || tm.symbols <= 0 {
		return false
	}
	for state, symbolTransitions := range tm.transitions {
		if int(state) < 0 || int(state) >= tm.states {
			return false
		}
		for symbol, transition := range symbolTransitions {
			if int(symbol) < 0 || int(symbol) >= tm.symbols {
				return false
			}
			writeSymbol := transition.symbol
			if int(writeSymbol) < 0 || int(writeSymbol) >= tm.symbols {
				return false
			}
		}
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
	for state, symbolTransitions := range wfa.transitions {
		if int(state) < 0 || int(state) >= wfa.states {
			return false
		}
		for symbol, transition := range symbolTransitions {
			if int(symbol) < 0 || int(symbol) >= wfa.symbols {
				return false
			}
			targetState := transition.wfaState
			if int(targetState) < 0 || int(targetState) >= wfa.states {
				return false
			}
			check(transition.weight)
		}
	}
	return true
}

func verifySymbolCompatibility(tm turingMachine, leftWFA, rightWFA dwfa) bool {
	return tm.symbols == leftWFA.symbols && tm.symbols == rightWFA.symbols
}

func verifySpecialSetsAreSubsets(wfa dwfa, specialSets specialSets) bool {
	for state := range specialSets.nonNegative {
		if int(state) < 0 || int(state) >= wfa.states {
			return false
		}
	}
	for state := range specialSets.nonPositive {
		if int(state) < 0 || int(state) >= wfa.states {
			return false
		}
	}
	return true
}

func verifyAcceptSetIsValid(tm turingMachine, leftWFA, rightWFA dwfa, acceptSet acceptSet) bool {
	for config, bounds := range acceptSet {
		if int(config.tmState) < 0 || int(config.tmState) >= tm.states {
			return false
		}
		if int(config.tmSymbol) < 0 || int(config.tmSymbol) >= tm.symbols {
			return false
		}
		if int(config.leftState) < 0 || int(config.leftState) >= leftWFA.states {
			return false
		}
		if int(config.rightState) < 0 || int(config.rightState) >= rightWFA.states {
			return false
		}
		lowerbound, lowerExists := bounds[LOWER]
		if lowerExists {
			check(lowerbound)
		}
		upperbound, upperExists := bounds[UPPER]
		if upperExists {
			check(upperbound)
		}
		if lowerExists && upperExists && lowerbound > upperbound {
			return false
		}
	}
	return true
}

func verifyLeadingBlankInvariant(wfa dwfa) bool {
	state := wfa.startState
	transition := wfa.transitions[state][0]
	return transition.wfaState == state && transition.weight == 0
}

func verifySpecialSetsHaveClaimedProperty(wfa dwfa, specialSets specialSets) bool {
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

	endNonPositive := specialSets.nonPositive.contains(endState)
	startNonPositive := specialSets.nonPositive.contains(startState)
	if endNonPositive {
		if !startNonPositive {
			return false
		}
		if weight > 0 {
			return false
		}
	}

	endNonNegative := specialSets.nonNegative.contains(endState)
	startNonNegative := specialSets.nonNegative.contains(startState)
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

func verifyStartConfigAccept(leftWFA, rightWFA dwfa, acceptSet acceptSet) bool {
	bounds, ok := acceptSet[config{TMSTARTSTATE, TMSTARTSYMBOL, leftWFA.startState, rightWFA.startState}]
	if !ok {
		return false
	}
	if lowerbound, ok := bounds[LOWER]; ok && lowerbound > 0 {
		return false
	}
	if upperbound, ok := bounds[UPPER]; ok && upperbound < 0 {
		return false
	}
	return true
}

func verifyNoHaltingConfigAccepted(tm turingMachine, acceptSet acceptSet) bool {
	for condition := range acceptSet {
		if condition.tmState < 0 || int(condition.tmState) >= tm.states {
			return false
		}
		if haltsNextStep(tm, condition.tmState, condition.tmSymbol) {
			return false
		}
	}
	return true
}

func haltsNextStep(tm turingMachine, tmState tmState, symbol symbol) bool {
	transition, ok := tm.transitions[tmState][symbol]
	if !ok {
		return true
	}
	if transition.tmState < 0 || int(transition.tmState) >= tm.states {
		return true
	}
	return false
}

func verifyForwardClosed(tm turingMachine, leftWFA, rightWFA dwfa, leftSpecialSets, rightSpecialSets specialSets, acceptSet acceptSet) bool {
	for config, bounds := range acceptSet {
		for _, nextConfigWithWeightChange := range nextConfigsWithWeightChange(config, tm, leftWFA, rightWFA) {
			if !nextConfigWithWeightChangeIsAccepted(nextConfigWithWeightChange, bounds, leftSpecialSets, rightSpecialSets, acceptSet) {
				return false
			}
		}
	}
	return true
}

func nextConfigsWithWeightChange(oldConfig config, tm turingMachine, leftWFA, rightWFA dwfa) []configWithWeight {
	result := []configWithWeight{}
	tmTransition, ok := tm.transitions[oldConfig.tmState][oldConfig.tmSymbol]
	if !ok {
		return result
	}
	switch tmTransition.direction {
	case L:
		for nextLeftState, leftStateTransitions := range leftWFA.transitions {
			for nextSymbol, leftTransition := range leftStateTransitions {
				if leftTransition.wfaState == oldConfig.leftState {
					rightTransition := rightWFA.transitions[oldConfig.rightState][tmTransition.symbol]

					nextConfig := config{tmTransition.tmState, nextSymbol, nextLeftState, rightTransition.wfaState}
					weightChange := rightTransition.weight - leftTransition.weight
					check(weightChange)

					result = append(result, configWithWeight{nextConfig, weightChange})
				}
			}
		}
	case R:
		for nextRightState, rightStateTransitions := range rightWFA.transitions {
			for nextSymbol, rightTransition := range rightStateTransitions {
				if rightTransition.wfaState == oldConfig.rightState {
					leftTransition := leftWFA.transitions[oldConfig.leftState][tmTransition.symbol]

					nextConfig := config{tmTransition.tmState, nextSymbol, leftTransition.wfaState, nextRightState}
					weightChange := leftTransition.weight - rightTransition.weight
					check(weightChange)

					result = append(result, configWithWeight{nextConfig, weightChange})
				}
			}
		}
	}

	return result
}

func nextConfigWithWeightChangeIsAccepted(nextConfigWithWeightChange configWithWeight, bounds bounds, leftSpecialSets, rightSpecialSets specialSets, acceptSet acceptSet) bool {
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

	//adjust bounds according to the special sets
	_, leftStateNonNegative := leftSpecialSets.nonNegative[nextConfig.leftState]
	_, rightStateNonNegative := rightSpecialSets.nonNegative[nextConfig.rightState]
	if leftStateNonNegative && rightStateNonNegative {
		if !lowerExists || lowerbound < 0 {
			lowerExists = true
			lowerbound = 0
		}
	}
	_, leftStateNonPositive := leftSpecialSets.nonPositive[nextConfig.leftState]
	_, rightStateNonPositive := rightSpecialSets.nonPositive[nextConfig.rightState]
	if leftStateNonPositive && rightStateNonPositive {
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
		return true
	}
	return acceptSetCountainsConfigBounds(acceptSet, nextConfig, nextBounds)
}

func acceptSetCountainsConfigBounds(acceptSet acceptSet, nextConfig config, nextBounds map[boundType]weight) bool {
	acceptBounds, ok := acceptSet[nextConfig]
	if !ok {
		return false
	}
	acceptedLower, acceptedLowerExists := acceptBounds[LOWER]
	nextLower, nextLowerExists := nextBounds[LOWER]
	if acceptedLowerExists && (!nextLowerExists || acceptedLower > nextLower) {
		return false
	}

	acceptedUpper, acceptedUpperExists := acceptBounds[UPPER]
	nextUpper, nextUpperExists := nextBounds[UPPER]
	if acceptedUpperExists && (!nextUpperExists || acceptedUpper < nextUpper) {
		return false
	}
	return true
}
