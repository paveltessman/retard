package longestsubstring

import (
	"math/rand"
	"strings"
	"testing"
)

// brute gives the answer of the task in the slow and plain way: it starts a
// window at every position and extends it until a letter comes back.
// The tests compare the solution with it.
func brute(s string) int {
	best := 0

	for i := range len(s) {
		var seen [256]bool
		for j := i; j < len(s); j++ {
			if seen[s[j]] {
				break
			}
			seen[s[j]] = true

			if j-i+1 > best {
				best = j - i + 1
			}
		}
	}

	return best
}

// printable gives n different characters out of the printable ASCII range,
// and leaves out the characters of skip.
// The task allows letters, digits, symbols, and spaces.
func printable(n int, skip string) string {
	var out []byte

	for c := byte(' '); c <= '~' && len(out) < n; c++ {
		if strings.IndexByte(skip, c) < 0 {
			out = append(out, c)
		}
	}

	return string(out)
}

// randomString gives a word of n characters out of alphabet.
func randomString(n int, alphabet string, rnd *rand.Rand) string {
	word := make([]byte, n)
	for i := range word {
		word[i] = alphabet[rnd.Intn(len(alphabet))]
	}
	return string(word)
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

func TestLengthOfLongestSubstring(t *testing.T) {
	cases := []struct {
		name string
		s    string
		want int
	}{
		{"Example1", "abcabcbb", 3},
		{"Example2", "bbbbb", 1},
		{"Example3", "pwwkew", 3},
		{"Empty", "", 0},
		{"OneCharacter", "a", 1},
		{"OneSpace", " ", 1},
		{"TwoSame", "aa", 1},
		{"TwoDifferent", "ab", 2},
		{"AllDifferent", "abcdef", 6},
		{"AllSame", "cccccccc", 1},
		{"RepeatAtTheEnds", "abcdea", 5},
		{"RepeatInTheMiddle", "abcbdef", 5},
		{"PairInFront", "aab", 2},
		{"PairAtTheBack", "cdd", 2},
		{"Palindrome", "abba", 2},
		{"PalindromeLonger", "abccba", 3},
		{"WindowStartStaysPut", "tmmzuxt", 5},
		{"WindowStartMovesOnce", "dvdf", 3},
		{"BestWindowLast", "aabaabcdefg", 7},
		{"BestWindowFirst", "abcdefgaabaa", 7},
		{"BestWindowInTheMiddle", "aaabcdefaaa", 6},
		{"Digits", "1231234", 4},
		{"Spaces", "a b c a", 3},
		{"Symbols", "!@#!@#$", 4},
		{"UpperAndLower", "aAbBaA", 4},
		{"MixedCharacters", "a1! a1!", 4},
		{"LongTail", "aaaaabcdef", 6},
		{"EveryPrintableCharacter", printable(95, ""), 95},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := lengthOfLongestSubstring(c.s); got != c.want {
				t.Errorf("lengthOfLongestSubstring(%q) = %d, want %d", c.s, got, c.want)
			}
		})
	}
}

// TestLengthOfLongestSubstringEveryShortString walks every string of up to
// 8 characters over three letters, and every string of up to 12 characters
// over two letters. It catches a window that starts or stops one step too
// late, and a window start that moves back.
func TestLengthOfLongestSubstringEveryShortString(t *testing.T) {
	alphabets := []struct {
		letters string
		maxLen  int
	}{
		{"abc", 8},
		{"ab", 12},
	}

	for _, a := range alphabets {
		for n := 0; n <= a.maxLen; n++ {
			failed := false

			everyString(a.letters, n, func(s string) {
				if failed {
					return
				}
				if got, want := lengthOfLongestSubstring(s), brute(s); got != want {
					t.Errorf("lengthOfLongestSubstring(%q) = %d, want %d", s, got, want)
					failed = true
				}
			})
		}
	}
}

// TestLengthOfLongestSubstringRandom compares the solution with brute on
// random strings. A small alphabet gives short windows and many repeats.
// A large alphabet gives long windows.
func TestLengthOfLongestSubstringRandom(t *testing.T) {
	rnd := rand.New(rand.NewSource(1))
	alphabets := []string{"ab", "abc", "abcde", "abcdefghijklmnopqrstuvwxyz", printable(95, "")}

	for _, alphabet := range alphabets {
		for range 200 {
			s := randomString(rnd.Intn(200), alphabet, rnd)

			if got, want := lengthOfLongestSubstring(s), brute(s); got != want {
				t.Fatalf("lengthOfLongestSubstring(%q) = %d, want %d", s, got, want)
			}
		}
	}
}

// TestLengthOfLongestSubstringNoRepeat gives strings of different characters
// only. The answer is the length of the whole string.
func TestLengthOfLongestSubstringNoRepeat(t *testing.T) {
	for n := 0; n <= 95; n++ {
		s := printable(n, "")

		if got := lengthOfLongestSubstring(s); got != n {
			t.Errorf("lengthOfLongestSubstring(%q) = %d, want %d", s, got, n)
		}
	}
}

// TestLengthOfLongestSubstringOneCharacter gives one character many times.
// The answer is 1 for every length.
func TestLengthOfLongestSubstringOneCharacter(t *testing.T) {
	for n := 1; n <= 200; n++ {
		s := strings.Repeat("x", n)

		if got := lengthOfLongestSubstring(s); got != 1 {
			t.Errorf("lengthOfLongestSubstring of %d times %q = %d, want 1", n, "x", got)
		}
	}
}

// TestLengthOfLongestSubstringBestWindowEverywhere puts the one window of
// different characters at every position of the string. The characters around
// the window repeat, so the answer is the length of the window.
func TestLengthOfLongestSubstringBestWindowEverywhere(t *testing.T) {
	const window = "wxyz"

	for n := 0; n <= 12; n++ {
		for at := 0; at <= n; at++ {
			// The filler holds the same character many times, and that
			// character is not part of the window.
			s := strings.Repeat("a", at) + window + strings.Repeat("a", n-at)

			// The window takes one "a" as well, in front of it or behind
			// it. It cannot take both, because the two are the same
			// character.
			want := len(window)
			if n > 0 {
				want++
			}

			if got := lengthOfLongestSubstring(s); got != want {
				t.Errorf("lengthOfLongestSubstring(%q) = %d, want %d", s, got, want)
			}
		}
	}
}

// TestLengthOfLongestSubstringLarge uses the largest string of the
// constraints. The last case holds the best window at the end, so the
// solution must read the whole string.
func TestLengthOfLongestSubstringLarge(t *testing.T) {
	const n = 100000

	t.Run("SameCharacter", func(t *testing.T) {
		s := strings.Repeat("a", n)

		if got := lengthOfLongestSubstring(s); got != 1 {
			t.Errorf("lengthOfLongestSubstring of %d times %q = %d, want 1", n, "a", got)
		}
	})

	t.Run("RepeatedBlock", func(t *testing.T) {
		// The block holds every printable character once. A window of 96
		// characters always crosses the end of a block and takes one
		// character twice, so the answer is the length of the block.
		block := printable(95, "")
		s := strings.Repeat(block, n/len(block)+1)[:n]

		if got := lengthOfLongestSubstring(s); got != len(block) {
			t.Errorf("lengthOfLongestSubstring of %d repeats of the alphabet = %d, want %d", n, got, len(block))
		}
	})

	t.Run("BestWindowAtTheEnd", func(t *testing.T) {
		// The front holds one character many times, so every window there
		// holds 1 character. The tail holds 90 other characters, all
		// different, and it takes the last character of the front as well.
		tail := printable(90, "a")
		s := strings.Repeat("a", n-len(tail)) + tail
		want := len(tail) + 1

		if got := lengthOfLongestSubstring(s); got != want {
			t.Errorf("lengthOfLongestSubstring of %d characters with the best window at the end = %d, want %d", n, got, want)
		}
	})
}

func BenchmarkLengthOfLongestSubstringSmall(b *testing.B) {
	const s = "abcabcbb"

	for b.Loop() {
		lengthOfLongestSubstring(s)
	}
}

// BenchmarkLengthOfLongestSubstringSameCharacter gives one character many
// times. The window start moves as often as the window end.
func BenchmarkLengthOfLongestSubstringSameCharacter(b *testing.B) {
	const n = 100000

	s := strings.Repeat("a", n)

	b.ResetTimer()
	for b.Loop() {
		lengthOfLongestSubstring(s)
	}
}

// BenchmarkLengthOfLongestSubstringLongWindow repeats the whole printable
// alphabet. The window stays long, so the solution holds many characters
// at the same time.
func BenchmarkLengthOfLongestSubstringLongWindow(b *testing.B) {
	const n = 100000

	block := printable(95, "")
	s := strings.Repeat(block, n/len(block)+1)[:n]

	b.ResetTimer()
	for b.Loop() {
		lengthOfLongestSubstring(s)
	}
}

// BenchmarkLengthOfLongestSubstringRandom gives a random string over the
// 26 lowercase letters: the windows have different lengths.
func BenchmarkLengthOfLongestSubstringRandom(b *testing.B) {
	const n = 100000

	s := randomString(n, "abcdefghijklmnopqrstuvwxyz", rand.New(rand.NewSource(2)))

	b.ResetTimer()
	for b.Loop() {
		lengthOfLongestSubstring(s)
	}
}
