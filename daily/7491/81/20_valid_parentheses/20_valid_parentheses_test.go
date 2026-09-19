package validparentheses

import (
	"math/rand"
	"strings"
	"testing"
	"time"
)

// brackets holds every character the task allows, in the order open, close.
const brackets = "()[]{}"

// closeOf gives the close bracket that belongs to an open bracket.
var closeOf = map[byte]byte{'(': ')', '[': ']', '{': '}'}

// brute gives the answer of the task in the slow and plain way: it takes the
// pairs "()", "[]" and "{}" out of the string again and again, until nothing
// more comes out. The string is valid when nothing is left. The tests compare
// the solution with it.
func brute(s string) bool {
	for {
		next := s
		for _, pair := range []string{"()", "[]", "{}"} {
			next = strings.ReplaceAll(next, pair, "")
		}

		if next == s {
			return s == ""
		}

		s = next
	}
}

// check calls isValid and compares the answer with want. It gives back false
// on the first error, so a loop over thousands of strings can stop.
func check(t *testing.T, s string, want bool) bool {
	t.Helper()

	if got := isValid(s); got != want {
		t.Errorf("isValid(%q) = %v, want %v", s, got, want)
		return false
	}

	return true
}

// everyString calls f once for every string of n characters over alphabet.
func everyString(alphabet string, n int, f func(string)) {
	word := make([]byte, n)

	var build func(i int)
	build = func(i int) {
		if i == n {
			f(string(word))
			return
		}
		for k := range len(alphabet) {
			word[i] = alphabet[k]
			build(i + 1)
		}
	}

	build(0)
}

// randomString gives a word of n characters out of alphabet.
func randomString(n int, alphabet string, rnd *rand.Rand) string {
	word := make([]byte, n)
	for i := range word {
		word[i] = alphabet[rnd.Intn(len(alphabet))]
	}

	return string(word)
}

// randomValid gives a valid string of n characters, or of n-1 characters when
// n is odd. The types and the depth are random, so the string mixes nesting
// and brackets side by side.
func randomValid(n int, rnd *rand.Rand) string {
	const opens = "([{"

	n -= n % 2

	out := make([]byte, 0, n)
	var open []byte // the close brackets that are still missing

	for len(out) < n {
		left := n - len(out)

		// A new open bracket needs a place for itself and a place for its
		// close bracket, next to the close brackets that are already open.
		if left >= len(open)+2 && (len(open) == 0 || rnd.Intn(2) == 0) {
			c := opens[rnd.Intn(len(opens))]
			out = append(out, c)
			open = append(open, closeOf[c])

			continue
		}

		out = append(out, open[len(open)-1])
		open = open[:len(open)-1]
	}

	return string(out)
}

// mutate changes one place of s: it drops a character, it writes another
// character over one, or it swaps two characters. A mutation of a valid string
// is a near miss, and sometimes it stays valid, so the tests take the answer
// from brute.
func mutate(s string, rnd *rand.Rand) string {
	if len(s) == 0 {
		return s
	}

	b := []byte(s)
	i := rnd.Intn(len(b))

	switch rnd.Intn(3) {
	case 0:
		return string(append(b[:i:i], b[i+1:]...))
	case 1:
		b[i] = brackets[rnd.Intn(len(brackets))]
	default:
		j := rnd.Intn(len(b))
		b[i], b[j] = b[j], b[i]
	}

	return string(b)
}

func TestIsValid(t *testing.T) {
	cases := []struct {
		name string
		s    string
		want bool
	}{
		{"Example1", "()", true},
		{"Example2", "()[]{}", true},
		{"Example3", "(]", false},
		{"Example4", "([])", true},
		{"Example5", "([)]", false},
		{"OneOpenRound", "(", false},
		{"OneCloseRound", ")", false},
		{"OneOpenSquare", "[", false},
		{"OneCloseCurly", "}", false},
		{"PairRound", "()", true},
		{"PairSquare", "[]", true},
		{"PairCurly", "{}", true},
		{"CloseBeforeOpen", ")(", false},
		{"CloseBeforeOpenSquare", "][", false},
		{"TwoOpens", "((", false},
		{"TwoCloses", "))", false},
		{"WrongTypeRoundSquare", "(]", false},
		{"WrongTypeSquareCurly", "[}", false},
		{"WrongTypeCurlyRound", "{)", false},
		{"NestedTwice", "(())", true},
		{"NestedThreeTypes", "([{}])", true},
		{"NestedWrongInnerType", "([{)}]", false},
		{"SideBySideThreeTypes", "(){}[]", true},
		{"MissingOneClose", "(()", false},
		{"MissingOneOpen", "())", false},
		{"ExtraCloseAtTheStart", ")()", false},
		{"ExtraOpenAtTheEnd", "()(", false},
		{"CrossedTypes", "([)]", false},
		{"CrossedTypesLonger", "{[(})]", false},
		{"NestedThenSideBySide", "(()){}", true},
		{"SideBySideThenNested", "{}(())", true},
		{"DeepNesting", "((((()))))", true},
		{"DeepNestingOneShort", "((((())))", false},
		{"DeepNestingOneTooMany", "((((())))))", false},
		{"DeepNestingWrongInnermost", "((((]))))", false},
		{"DeepNestingWrongOutermost", "[(((())))}", false},
		{"ManyPairs", "()()()()()", true},
		{"ManyPairsWithABreak", "()()(]()()", false},
		{"AllTypesNested", "{[()]}", true},
		{"AllTypesNestedSwapped", "{[(])}", false},
		{"SameCountsWrongOrder", "))((", false},
		{"SameCountsWrongOrderMixed", "}{][)(", false},
		{"RightPrefixWrongTail", "()()()[", false},
		{"WrongPrefixRightTail", "]()()()", false},
		{"LongValidMix", "({[]})[]{()}", true},
		{"LongInvalidMix", "({[]})[]{()", false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			check(t, c.s, c.want)
		})
	}
}

// TestIsValidEveryShortString walks every short string over the brackets and
// compares the solution with brute. It catches a solution that counts the
// brackets without their type, a solution that reads the order wrong, and a
// solution that gives true for a string that is not closed at the end.
func TestIsValidEveryShortString(t *testing.T) {
	alphabets := []struct {
		letters string
		maxLen  int
	}{
		{brackets, 7},
		{"()[]", 9},
		{"()", 14},
	}

	for _, a := range alphabets {
		for n := 1; n <= a.maxLen; n++ {
			failed := false

			everyString(a.letters, n, func(s string) {
				if failed {
					return
				}
				if !check(t, s, brute(s)) {
					failed = true
				}
			})

			if failed {
				return
			}
		}
	}
}

// TestIsValidRandom compares the solution with brute on random strings. Most
// random strings are not valid, so the test also builds valid strings and
// changes one place of them. Such a near miss is the hard case.
func TestIsValidRandom(t *testing.T) {
	rnd := rand.New(rand.NewSource(1))

	for range 2000 {
		n := 1 + rnd.Intn(60)
		s := randomString(n, brackets, rnd)

		if !check(t, s, brute(s)) {
			return
		}
	}

	for range 2000 {
		s := randomValid(2+rnd.Intn(60), rnd)

		if !check(t, s, true) {
			return
		}

		for range 3 {
			m := mutate(s, rnd)

			if !check(t, m, brute(m)) {
				return
			}
		}
	}
}

// TestIsValidKeepsNeighbours gives a part of a longer string to isValid. The
// characters around the part change the answer, so a solution that reads
// outside of s fails here.
func TestIsValidKeepsNeighbours(t *testing.T) {
	const whole = "]]([])[[" // the middle part is valid, the whole string is not

	check(t, whole[2:6], true)
	check(t, whole, false)
	check(t, whole[2:7], false)
	check(t, whole[1:6], false)

	const other = "(({[]}))"

	check(t, other[2:6], true)
	check(t, other[2:5], false)
	check(t, other[3:6], false)
}

// TestIsValidLarge uses the largest string of the constraints. The answers
// come from the way the strings are built, and the wrong place sits at the
// end, so the solution must read the whole string.
func TestIsValidLarge(t *testing.T) {
	const n = 10000

	cases := []struct {
		name string
		s    string
		want bool
	}{
		{"DeepNesting", strings.Repeat("(", n/2) + strings.Repeat(")", n/2), true},
		{"DeepNestingWrongLastType", strings.Repeat("(", n/2) + strings.Repeat(")", n/2-1) + "]", false},
		{"DeepNestingOneOpenTooMany", strings.Repeat("(", n/2+1) + strings.Repeat(")", n/2-1), false},
		{"DeepNestingOneCloseTooMany", strings.Repeat("(", n/2-1) + strings.Repeat(")", n/2+1), false},
		{"AllOpen", strings.Repeat("(", n), false},
		{"AllClose", strings.Repeat(")", n), false},
		{"OpensThenCloses", strings.Repeat("(", n/2) + strings.Repeat("]", n/2), false},
		{"ManyPairs", strings.Repeat("()", n/2), true},
		{"ManyPairsBrokenAtTheEnd", strings.Repeat("()", n/2-1) + "(]", false},
		{"ManyPairsBrokenAtTheStart", "(]" + strings.Repeat("()", n/2-1), false},
		{"MixedTypesNested", strings.Repeat("([{", n/6) + strings.Repeat("}])", n/6), true},
		{"MixedTypesNestedWrongInnermost", strings.Repeat("([{", n/6) + ")" + strings.Repeat("}])", n/6)[1:], false},
		{"EveryTypeSideBySide", strings.Repeat("(){}[]", n/6), true},
		{"NestedThreeTypes", nestedThreeTypes(n), true},
		{"NestedThreeTypesCrossedInTheMiddle", crossAtTheMiddle(nestedThreeTypes(n)), false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := isValid(c.s); got != c.want {
				t.Errorf("isValid of the string %s of %d characters = %v, want %v", c.name, len(c.s), got, c.want)
			}
		})
	}
}

// nestedThreeTypes gives a valid string of n characters: the three types of
// brackets one inside the other, again and again, from the outside to the
// middle and back.
func nestedThreeTypes(n int) string {
	const opens = "([{"

	half := n / 2

	front := make([]byte, half)
	back := make([]byte, half)

	for i := range half {
		c := opens[i%len(opens)]
		front[i] = c
		back[half-1-i] = closeOf[c]
	}

	return string(front) + string(back)
}

// crossAtTheMiddle swaps the two characters in the middle of s. A valid string
// of nested brackets loses its order, but it keeps the count of every type.
func crossAtTheMiddle(s string) string {
	b := []byte(s)
	i := len(b) / 2

	b[i-1], b[i] = b[i], b[i-1]

	return string(b)
}

// callWithLimit runs isValid on s in another goroutine and waits for the
// limit. It gives back ok = false when no answer comes in time.
func callWithLimit(s string, limit time.Duration) (got, ok bool) {
	done := make(chan bool, 1)

	go func() {
		done <- isValid(s)
	}()

	select {
	case got := <-done:
		return got, true
	case <-time.After(limit):
		return false, false
	}
}

// TestIsValidScales gives strings that are much longer than the constraints
// allow. A solution that reads the string a few times needs milliseconds here.
// A solution that walks the whole string again for every pair of brackets
// needs hours, so it runs into the limit.
func TestIsValidScales(t *testing.T) {
	if testing.Short() {
		t.Skip("the test needs very long strings")
	}

	const n = 2000000
	const limit = 15 * time.Second

	cases := []struct {
		name string
		s    string
		want bool
	}{
		{"DeepNesting", strings.Repeat("(", n/2) + strings.Repeat(")", n/2), true},
		{"ManyPairs", strings.Repeat("()", n/2), true},
		{"NestedThreeTypes", nestedThreeTypes(n), true},
		{"BrokenAtTheEnd", nestedThreeTypes(n)[:n-1] + "(", false},
	}

	for _, c := range cases {
		got, ok := callWithLimit(c.s, limit)

		if !ok {
			t.Fatalf("isValid needed more than %v for the string %s of %d characters, so it reads the string again and again", limit, c.name, len(c.s))
		}
		if got != c.want {
			t.Fatalf("isValid of the string %s of %d characters = %v, want %v", c.name, len(c.s), got, c.want)
		}
	}
}

func BenchmarkIsValid(b *testing.B) {
	const n = 100000

	shapes := []struct {
		name string
		s    string
	}{
		{"Short", "()[]{}"},
		{"DeepNesting", strings.Repeat("(", n/2) + strings.Repeat(")", n/2)},
		{"ManyPairs", strings.Repeat("()", n/2)},
		{"NestedThreeTypes", nestedThreeTypes(n)},
		{"WrongAtTheStart", "]" + strings.Repeat("()", n/2)},
		{"WrongAtTheEnd", strings.Repeat("()", n/2-1) + "(]"},
		{"AllOpen", strings.Repeat("(", n)},
		{"Random", randomString(n, brackets, rand.New(rand.NewSource(2)))},
	}

	for _, s := range shapes {
		b.Run(s.name, func(b *testing.B) {
			for b.Loop() {
				isValid(s.s)
			}
		})
	}
}
