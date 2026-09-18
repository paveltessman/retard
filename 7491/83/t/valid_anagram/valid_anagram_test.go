package validanagram

import (
	"flag"
	"math/rand"
	"strings"
	"testing"
)

// followUp turns on the Unicode tests of the follow-up question.
// They stay off by default, because a solution that counts 26 letters
// panics on a letter outside a..z. Run them with:
//
//	go test ./1/t/valid_anagram -followup
var followUp = flag.Bool("followup", false, "run the Unicode follow-up tests")

// shuffled gives back the letters of s in a new order.
// The letters stay the same, so the answer of the task is true.
func shuffled(s string, rnd *rand.Rand) string {
	letters := []byte(s)
	rnd.Shuffle(len(letters), func(i, j int) {
		letters[i], letters[j] = letters[j], letters[i]
	})
	return string(letters)
}

// randomWord gives a word of n lowercase letters.
func randomWord(n int, rnd *rand.Rand) string {
	letters := make([]byte, n)
	for i := range letters {
		letters[i] = byte('a' + rnd.Intn(26))
	}
	return string(letters)
}

func TestIsAnagram(t *testing.T) {
	cases := []struct {
		name string
		s    string
		t    string
		want bool
	}{
		{"Example1", "anagram", "nagaram", true},
		{"Example2", "rat", "car", false},
		{"SameString", "listen", "listen", true},
		{"Reversed", "abcdef", "fedcba", true},
		{"OneLetterSame", "a", "a", true},
		{"OneLetterDifferent", "a", "b", false},
		{"TwoLetters", "ab", "ba", true},
		{"ShorterSecond", "anagram", "nagara", false},
		{"LongerSecond", "nagara", "anagram", false},
		{"ExtraLetterAtEnd", "abc", "abcd", false},
		{"MissingLetterAtEnd", "abcd", "abc", false},
		{"SameLettersDifferentCounts", "aabb", "abbb", false},
		{"RepeatedLetters", "aacc", "ccaa", true},
		{"RepeatedLettersWrong", "aacc", "ccac", false},
		{"AllOneLetter", "aaaa", "aaaa", true},
		{"AllOneLetterWrong", "aaaa", "aaab", false},
		{"Alphabet", "abcdefghijklmnopqrstuvwxyz", "zyxwvutsrqponmlkjihgfedcba", true},
		{"AlphabetWrong", "abcdefghijklmnopqrstuvwxyz", "zyxwvutsrqponmlkjihgfedcbb", false},
		{"FirstLetterOnly", "az", "za", true},
		{"LastLetterWrong", "abcz", "abcy", false},
		{"SameSumOfLetters", "ad", "bc", false},
		{"SameSumOfLettersLonger", "at", "bs", false},
		{"NoSharedLetter", "abc", "xyz", false},
		{"OneSharedLetter", "abc", "axy", false},
		{"BothEmpty", "", "", true},
		{"EmptyFirst", "", "a", false},
		{"EmptySecond", "a", "", false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := isAnagram(c.s, c.t); got != c.want {
				t.Errorf("isAnagram(%q, %q) = %v, want %v", c.s, c.t, got, c.want)
			}
		})
	}
}

// TestIsAnagramSymmetric checks that the order of the two words means nothing.
// A solution that reads only one of the two words fails here.
func TestIsAnagramSymmetric(t *testing.T) {
	pairs := [][2]string{
		{"anagram", "nagaram"},
		{"rat", "car"},
		{"aabb", "abbb"},
		{"abc", "abcd"},
		{"silent", "listen"},
		{"apple", "papel"},
		{"apple", "papal"},
	}

	for _, p := range pairs {
		forward := isAnagram(p[0], p[1])
		backward := isAnagram(p[1], p[0])

		if forward != backward {
			t.Errorf("isAnagram(%q, %q) = %v, but isAnagram(%q, %q) = %v",
				p[0], p[1], forward, p[1], p[0], backward)
		}
	}
}

// TestIsAnagramShuffles shuffles a word many times.
// Every shuffle holds the same letters, so the answer is always true.
func TestIsAnagramShuffles(t *testing.T) {
	rnd := rand.New(rand.NewSource(1))

	for n := 1; n <= 40; n++ {
		word := randomWord(n, rnd)
		for range 20 {
			other := shuffled(word, rnd)

			if !isAnagram(word, other) {
				t.Fatalf("isAnagram(%q, %q) = false, want true", word, other)
			}
		}
	}
}

// TestIsAnagramOneLetterChanged changes one letter of a shuffled word.
// The new letter is not in the word, so the answer is false.
// The test moves the change through every position of the word.
func TestIsAnagramOneLetterChanged(t *testing.T) {
	rnd := rand.New(rand.NewSource(2))

	for n := 1; n <= 30; n++ {
		// The word holds letters a..m only, so 'z' is never one of them.
		word := make([]byte, n)
		for i := range word {
			word[i] = byte('a' + rnd.Intn(13))
		}
		s := string(word)

		for i := range n {
			letters := []byte(shuffled(s, rnd))
			letters[i] = 'z'
			other := string(letters)

			if isAnagram(s, other) {
				t.Fatalf("isAnagram(%q, %q) = true, want false", s, other)
			}
		}
	}
}

// TestIsAnagramCountMatters moves one letter from one word to the other.
// The two words keep the same length and the same letter set,
// so only the counts tell them apart.
func TestIsAnagramCountMatters(t *testing.T) {
	cases := [][2]string{
		{"aab", "abb"},
		{"aaab", "abbb"},
		{"xxyz", "xyzz"},
		{"mississippi", "misisssippi"},
		{"mississippi", "mississipii"},
	}

	for _, c := range cases {
		s, other := c[0], c[1]
		want := sortedLetters(s) == sortedLetters(other)

		if got := isAnagram(s, other); got != want {
			t.Errorf("isAnagram(%q, %q) = %v, want %v", s, other, got, want)
		}
	}
}

// sortedLetters gives the letters of s in alphabetical order.
// The test uses it to work out the answer without the solution.
func sortedLetters(s string) string {
	var count [26]int
	for i := range len(s) {
		count[s[i]-'a']++
	}

	var b strings.Builder
	for i, n := range count {
		b.WriteString(strings.Repeat(string(rune('a'+i)), n))
	}
	return b.String()
}

// TestIsAnagramLarge uses the largest words of the constraints.
// The two words differ in the last letter only, so the solution must
// read both words to the end.
func TestIsAnagramLarge(t *testing.T) {
	const n = 50000

	rnd := rand.New(rand.NewSource(3))
	s := randomWord(n, rnd)

	same := shuffled(s, rnd)
	if !isAnagram(s, same) {
		t.Errorf("isAnagram of a large word and its shuffle = false, want true")
	}

	letters := []byte(same)
	// Turn the last letter into the next one of the alphabet.
	letters[n-1] = 'a' + (letters[n-1]-'a'+1)%26
	if isAnagram(s, string(letters)) {
		t.Errorf("isAnagram of a large word and a changed shuffle = true, want false")
	}
}

// TestIsAnagramUnicode answers the follow-up question.
// It runs only with the -followup flag. See the note on followUp.
func TestIsAnagramUnicode(t *testing.T) {
	if !*followUp {
		t.Skip("run with -followup to test the Unicode follow-up")
	}

	cases := []struct {
		name string
		s    string
		t    string
		want bool
	}{
		{"Accents", "café", "féac", true},
		{"AccentsWrong", "café", "cafe", false},
		{"Cyrillic", "привет", "тевирп", true},
		{"CyrillicWrong", "привет", "тевирб", false},
		{"Greek", "λόγος", "ςογόλ", true},
		{"GreekFinalSigma", "λόγος", "σογόλ", false},
		{"Emoji", "a🙂b", "b🙂a", true},
		{"EmojiWrong", "a🙂b", "b🙃a", false},
		{"SameBytesDifferentOrder", "日本語", "語本日", true},
		{"UpperAndLower", "Ab", "bA", true},
		{"CaseMatters", "Ab", "ab", false},
		{"DifferentLengthInRunes", "ää", "ä", false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := isAnagram(c.s, c.t); got != c.want {
				t.Errorf("isAnagram(%q, %q) = %v, want %v", c.s, c.t, got, c.want)
			}
		})
	}
}

func BenchmarkIsAnagramSmall(b *testing.B) {
	for b.Loop() {
		isAnagram("anagram", "nagaram")
	}
}

// BenchmarkIsAnagramLarge uses the largest words of the constraints.
// A counting solution stays near the word length. A sorting solution
// pays for the sort of both words.
func BenchmarkIsAnagramLarge(b *testing.B) {
	const n = 50000

	rnd := rand.New(rand.NewSource(4))
	s := randomWord(n, rnd)
	other := shuffled(s, rnd)

	b.ResetTimer()
	for b.Loop() {
		isAnagram(s, other)
	}
}

// BenchmarkIsAnagramEarlyExit gives two words that differ in length.
// A solution that compares the lengths first answers at once.
func BenchmarkIsAnagramEarlyExit(b *testing.B) {
	const n = 50000

	rnd := rand.New(rand.NewSource(5))
	s := randomWord(n, rnd)
	other := s[:n-1]

	b.ResetTimer()
	for b.Loop() {
		isAnagram(s, other)
	}
}
