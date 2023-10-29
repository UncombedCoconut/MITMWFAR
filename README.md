# MITM-WFAR

Using weighted automatons to describe non-regular languages that solve the halting problem for some turing machines

# Theory
 
We extend a Meet-in-the-Middle DFA verifier to work with some irregular languages. We do this by using weighted automatons with integer weights instead of DFAs. 

When we run the left and right half of the tape through the WA any word will result in a weight in addition to the automaton state. We now only accept configurations if the sum of the weights of the sides is correct. For each combination of (head configuration, left state, right state) we accept if the weight sum is in a given interval of integers.

It is fairly simple to check if such a language given accept set is forward-closed under TM transitions. Since we know the weight of the edges you follow when going left/right by one symbol we can keep track of the difference while doing the usual checks for MITM-DFAs.

Unless regular MITM-DFA would have solved the TM we'll also need to make use of the weights. To do so we consider special sets of states in the WA: States were the weight is always nonnegative or nonpositive. Note that in a WA with only positive weights all states will always have a nonnegative weight. States that cannot be reached after passing through a weighted transition will also be nonpositive. These special sets allow us to rule out some combinations of WA states and weight sums when checking if an accept set is forward-closed: If both left and right WA are in a nonpositive state then the weight sum can never be positive. So even if such a combination arises from the typical MITM-DFA check we can disregard it since these residues can never be the result of an actual tm configuration.

If the accept set is forward-closed, accepts the starting configuration and doesn't accept any head-configurations that halt, then we know that the tm can never halt.

There is another explanation of this [on discord](https://discord.com/channels/960643023006490684/960655108578881597/1085884481241616444). Though with slight differences, like using the difference of weights between the WA instead of the sum.

# Certificates

The certificates this decider uses should be usable to prove that the tm specified in the certificate doesn't halt with minimal additional computation required. They are given in text format over multiple lines. Certificates for many machines can be appended into the same file.

## Full Certificate

The full certificates has 6 lines:
1. the TM in standard text format
2. the left WA
3. the right WA
4. the left special sets
5. the right special sets
6. the accept set with all accepted 6-tuples of (tm state, tm symbol, left WA state, right WA state, lower bound of weight sum, upper bound of weight sum)

When checking the certificates the decider ensures that all given information is correct. It checks that the states in the special sets are indeed nonnegative/nonpositive and the accept set has the required properties. 

## Short Certificate

Short certificates only include the first 3 lines of the full certificate, reminiscent of MITM-DFA certificates, where the accept sets can be derived from the DFA. Here we can obtain the special sets of the WA easily enough. The accept set can be derived by starting with the tuple that accepts the start configuration and then expanding the accept set as necessary.

However, other than in the DFA case this is not a deterministic process as we have the potential to have infinite accept sets with ever growing weight intervals. In an attempt to find finite accept sets in those cases we switch to unbounded intervals after there length exceeds 1000. This heuristic makes the process non-deterministic and sensitive to implementation details. While this implementation is consistent and will always extend a given short certificate to the same full certificate that might not hold with reproductions. So this is not a true certificate. It is however easier to digest for humans and should contain enough information to find the full solution quickly.

# Search Strategy

To find a proof for a given TM we essentially want to try all possible WA pairs for the short certificate with the TM and see what sticks. But that search space is very big and grows exponentially as we consider bigger WA, so we have to put some limits in and try to go through it efficiently.

In my search I only consider WA that are based on a DFA that is closed under TM transitions. These base DFA have an explicit dead state. Using MITM-DFA checks I ensure that this dead state will never be reached by the TM configurations. (I do allow halting head configurations to occur here.) I can enumerate these DFA in a process similar to enumerating TMs in TNF by starting with a base DFA pair and changing the transitions to the dead state into all possible non-dead transitions when they come up. I can then put a bound on the number of non-dead transitions to limit the search space.

For a given base DFA I then try all possible pairs of transitions, one left and one right, to weigh with 1 and -1. There is an option to try additional weight pairs, but I have not had any success with that.

# Usage

The decider reads from stdin and outputs to stdout. With `-pm=0` (or by default) it will print all TM for which it found a proof. With `-pm=1` it will print short certificates for those TM. With `-pm=2` it will print full certicates.

With `-sc` it will read short certificates from the input and verify them. With `-fc` it will read and verify full certificates.

With `-n` it will read a list of TM and try to decide them. It will search through WA with up to n non-dead transitions. `-m` can be added to transform the WA just before trying to build the accept set in order to give them a m long memory of the last WA transitions used.

Examples:
```
MITMWFAR -n=9 -m=1 -pm=1 < holdouts.std.txt > solved.sc.txt
MITMWFAR -sc -pm=2 < solved.sc.txt > solved.fc.txt
MITMWFAR -fc < solved.fc.txt
```
