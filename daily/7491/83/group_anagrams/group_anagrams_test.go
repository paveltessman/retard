package groupanagrams

import (
	"math/rand"
	"slices"
	"strings"
	"testing"
)

// signature gives the letters of s in alphabetical order.
// Two strings are anagrams when their signatures are equal.
// The tests work out the answer with it, without the solution.
func signature(s string) string {
	letters := []byte(s)
	slices.Sort(letters)
	return string(letters)
}

// canonical gives a form of the answer that does not depend on the order
// of the groups, or on the order of the words inside one group.
// The task accepts the answer in any order, so the tests compare this form.
func canonical(groups [][]string) [][]string {
	out := make([][]string, len(groups))
	for i, group := range groups {
		words := slices.Clone(group)
		slices.Sort(words)
		out[i] = words
	}
	slices.SortFunc(out, slices.Compare)
	return out
}

// wantGroups builds the answer of the task from strs.
// It puts the words with the same signature in one group.
func wantGroups(strs []string) [][]string {
	bySignature := make(map[string][]string)
	var order []string

	for _, s := range strs {
		sig := signature(s)
		if _, ok := bySignature[sig]; !ok {
			order = append(order, sig)
		}
		bySignature[sig] = append(bySignature[sig], s)
	}

	groups := make([][]string, 0, len(order))
	for _, sig := range order {
		groups = append(groups, bySignature[sig])
	}
	return groups
}

// checkGrouping checks that got is a real answer for strs.
// It reports an empty group, a group that mixes two different words,
// two groups that hold anagrams of the same word, a lost word,
// and a word that comes back more often than the input holds it.
func checkGrouping(t *testing.T, strs []string, got [][]string) {
	t.Helper()

	missing := make(map[string]int) // map[word]count
	for _, s := range strs {
		missing[s]++
	}

	groupOf := make(map[string]int) // map[signature]group index
	for i, group := range got {
		if len(group) == 0 {
			t.Errorf("group %d is empty", i)
			continue
		}

		sig := signature(group[0])
		for _, word := range group {
			if other := signature(word); other != sig {
				t.Errorf("group %d holds %q and %q, which are not anagrams", i, group[0], word)
			}

			missing[word]--
			if missing[word] < 0 {
				t.Errorf("the answer gives %q more often than the input holds it", word)
			}
		}

		if other, ok := groupOf[sig]; ok {
			t.Errorf("group %d and group %d both hold anagrams of %q", other, i, group[0])
		}
		groupOf[sig] = i
	}

	for word, count := range missing {
		if count > 0 {
			t.Errorf("the answer lost %q (missing copies: %d)", word, count)
		}
	}
}

// check runs both tests of the answer: the report of checkGrouping,
// and the comparison with the expected groups.
func check(t *testing.T, strs []string, got [][]string, want [][]string) {
	t.Helper()

	checkGrouping(t, strs, got)

	if !slices.EqualFunc(canonical(got), canonical(want), slices.Equal) {
		t.Errorf("groupAnagrams(%q) = %q, want %q", strs, canonical(got), canonical(want))
	}
}

// randomWord gives a word of n lowercase letters.
func randomWord(n int, rnd *rand.Rand) string {
	letters := make([]byte, n)
	for i := range letters {
		letters[i] = byte('a' + rnd.Intn(26))
	}
	return string(letters)
}

// shuffled gives back the letters of s in a new order.
// The letters stay the same, so the word belongs to the group of s.
func shuffled(s string, rnd *rand.Rand) string {
	letters := []byte(s)
	rnd.Shuffle(len(letters), func(i, j int) {
		letters[i], letters[j] = letters[j], letters[i]
	})
	return string(letters)
}

func TestGroupAnagrams(t *testing.T) {
	cases := []struct {
		name string
		strs []string
		want [][]string
	}{
		{
			"Example1",
			[]string{"eat", "tea", "tan", "ate", "nat", "bat"},
			[][]string{{"bat"}, {"nat", "tan"}, {"ate", "eat", "tea"}},
		},
		{"Example2", []string{""}, [][]string{{""}}},
		{"Example3", []string{"a"}, [][]string{{"a"}}},
		{"TwoAnagrams", []string{"listen", "silent"}, [][]string{{"listen", "silent"}}},
		{"TwoDifferent", []string{"abc", "xyz"}, [][]string{{"abc"}, {"xyz"}}},
		{"SameWordTwice", []string{"ab", "ab"}, [][]string{{"ab", "ab"}}},
		{"SameWordThreeTimes", []string{"ab", "ab", "ab"}, [][]string{{"ab", "ab", "ab"}}},
		{"EmptyStrings", []string{"", "", ""}, [][]string{{"", "", ""}}},
		{"EmptyAndWord", []string{"", "a"}, [][]string{{""}, {"a"}}},
		{"OneLetterEach", []string{"a", "b", "a", "c"}, [][]string{{"a", "a"}, {"b"}, {"c"}}},
		{
			"AllDifferent",
			[]string{"a", "ab", "abc", "abcd"},
			[][]string{{"a"}, {"ab"}, {"abc"}, {"abcd"}},
		},
		{
			"OneGroup",
			[]string{"abc", "acb", "bac", "bca", "cab", "cba"},
			[][]string{{"abc", "acb", "bac", "bca", "cab", "cba"}},
		},
		{
			"SameLettersDifferentCounts",
			[]string{"aab", "abb", "aabb"},
			[][]string{{"aab"}, {"abb"}, {"aabb"}},
		},
		{
			"SameLetterSetDifferentLength",
			[]string{"ab", "aab", "aaab", "ba"},
			[][]string{{"ab", "ba"}, {"aab"}, {"aaab"}},
		},
		{
			"SameSumOfLetters",
			[]string{"ad", "bc", "da", "cb"},
			[][]string{{"ad", "da"}, {"bc", "cb"}},
		},
		{
			"SameSumOfLettersLonger",
			[]string{"abc", "aad", "cba", "daa"},
			[][]string{{"abc", "cba"}, {"aad", "daa"}},
		},
		{
			"GroupsOfDifferentSize",
			[]string{"tea", "eat", "tan", "nat", "bat", "ate", "tab"},
			[][]string{{"ate", "eat", "tea"}, {"nat", "tan"}, {"bat", "tab"}},
		},
		{
			"LongWords",
			[]string{strings.Repeat("ab", 50), strings.Repeat("ba", 50), strings.Repeat("a", 100)},
			[][]string{{strings.Repeat("ab", 50), strings.Repeat("ba", 50)}, {strings.Repeat("a", 100)}},
		},
		{"NoStrings", []string{}, [][]string{}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			check(t, c.strs, groupAnagrams(c.strs), c.want)
		})
	}
}

// TestGroupAnagramsRandom builds a list out of shuffles of a few base words.
// It checks the answer against the groups that the signatures give.
func TestGroupAnagramsRandom(t *testing.T) {
	rnd := rand.New(rand.NewSource(1))

	for range 200 {
		var strs []string
		for range 1 + rnd.Intn(6) {
			base := randomWord(1+rnd.Intn(8), rnd)
			for range 1 + rnd.Intn(4) {
				strs = append(strs, shuffled(base, rnd))
			}
		}
		rnd.Shuffle(len(strs), func(i, j int) {
			strs[i], strs[j] = strs[j], strs[i]
		})

		check(t, strs, groupAnagrams(slices.Clone(strs)), wantGroups(strs))
	}
}

// TestGroupAnagramsOneGroupPerWord gives words with different signatures only.
// The answer holds one group for every word.
func TestGroupAnagramsOneGroupPerWord(t *testing.T) {
	for n := 1; n <= 30; n++ {
		strs := make([]string, n)
		for i := range strs {
			// Every word holds i+1 letters "a" and fills the rest with "b".
			// The words keep the same length, but no two of them are anagrams.
			strs[i] = strings.Repeat("a", i+1) + strings.Repeat("b", 30-i)
		}

		got := groupAnagrams(slices.Clone(strs))

		checkGrouping(t, strs, got)
		if len(got) != n {
			t.Fatalf("groupAnagrams of %d different words gave %d groups, want %d", n, len(got), n)
		}
	}
}

// TestGroupAnagramsSingleGroup gives shuffles of one word only.
// The answer holds one group with every word in it.
func TestGroupAnagramsSingleGroup(t *testing.T) {
	rnd := rand.New(rand.NewSource(2))

	for n := 1; n <= 30; n++ {
		base := randomWord(1+n%9, rnd)
		strs := make([]string, n)
		for i := range strs {
			strs[i] = shuffled(base, rnd)
		}

		got := groupAnagrams(slices.Clone(strs))

		checkGrouping(t, strs, got)
		if len(got) != 1 {
			t.Fatalf("groupAnagrams of %d shuffles of %q gave %d groups, want 1", n, base, len(got))
		}
	}
}

// TestGroupAnagramsKeepsInput checks that groupAnagrams leaves the input alone.
// A solution that sorts the words in place changes the array of the caller.
func TestGroupAnagramsKeepsInput(t *testing.T) {
	strs := []string{"eat", "tea", "tan", "ate", "nat", "bat"}
	want := slices.Clone(strs)

	groupAnagrams(strs)

	if !slices.Equal(strs, want) {
		t.Errorf("groupAnagrams changed the input to %q, want %q", strs, want)
	}
}

// TestGroupAnagramsLarge uses the largest input of the constraints:
// 10000 words of 100 letters each. It also puts one word of its own group
// at the end, so the solution must read the whole array.
func TestGroupAnagramsLarge(t *testing.T) {
	const (
		count  = 10000
		length = 100
		bases  = 50
	)

	rnd := rand.New(rand.NewSource(3))
	baseWords := make([]string, bases)
	for i := range baseWords {
		// Every base word holds i+1 letters "b" and fills the rest with "a",
		// so no two base words are anagrams.
		baseWords[i] = strings.Repeat("b", i+1) + strings.Repeat("a", length-i-1)
	}

	strs := make([]string, 0, count)
	for i := range count - 1 {
		strs = append(strs, shuffled(baseWords[i%bases], rnd))
	}
	strs = append(strs, strings.Repeat("z", length))

	got := groupAnagrams(slices.Clone(strs))

	checkGrouping(t, strs, got)
	if len(got) != bases+1 {
		t.Errorf("groupAnagrams of %d words gave %d groups, want %d", count, len(got), bases+1)
	}
}

func BenchmarkGroupAnagramsSmall(b *testing.B) {
	strs := []string{"eat", "tea", "tan", "ate", "nat", "bat"}

	for b.Loop() {
		groupAnagrams(strs)
	}
}

// BenchmarkGroupAnagramsLarge uses the largest input of the constraints.
// A solution that counts the 26 letters stays near the size of the input.
// A solution that sorts every word pays for every sort.
func BenchmarkGroupAnagramsLarge(b *testing.B) {
	const (
		count  = 10000
		length = 100
	)

	rnd := rand.New(rand.NewSource(4))
	strs := make([]string, count)
	for i := range strs {
		strs[i] = randomWord(length, rnd)
	}

	b.ResetTimer()
	for b.Loop() {
		groupAnagrams(strs)
	}
}

// BenchmarkGroupAnagramsOneGroup gives shuffles of one word only.
// Every word lands in the same group, so the solution pays for every
// word of the group.
func BenchmarkGroupAnagramsOneGroup(b *testing.B) {
	const (
		count  = 10000
		length = 100
	)

	rnd := rand.New(rand.NewSource(5))
	base := randomWord(length, rnd)
	strs := make([]string, count)
	for i := range strs {
		strs[i] = shuffled(base, rnd)
	}

	b.ResetTimer()
	for b.Loop() {
		groupAnagrams(strs)
	}
}
