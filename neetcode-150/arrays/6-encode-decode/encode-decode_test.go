package encodedecode

import (
	"math/rand"
	"slices"
	"strings"
	"testing"
)

// check sends strs through encode and decode and compares the result with
// strs. An empty list and a nil list count as equal. It also reads strs back
// after encode, so it catches a solution that changes the input. It gives
// back false on an error.
func check(t *testing.T, strs []string) bool {
	t.Helper()

	in := slices.Clone(strs)
	enc := encode(in)
	ok := true

	if !slices.Equal(in, strs) {
		t.Errorf("encode(%q) changed the list to %q", strs, in)
		ok = false
	}
	if got := decode(enc); !slices.Equal(got, strs) {
		t.Errorf("decode(encode(%q)) = %q, want %q (encoded: %q)", strs, got, strs, enc)
		ok = false
	}

	return ok
}

// randomList gives a list of n strings. Each string has a random length
// from 0 to maxLen and holds random bytes from alphabet.
func randomList(rnd *rand.Rand, n, maxLen int, alphabet string) []string {
	strs := make([]string, n)
	for i := range strs {
		b := make([]byte, rnd.Intn(maxLen+1))
		for j := range b {
			b[j] = alphabet[rnd.Intn(len(alphabet))]
		}
		strs[i] = string(b)
	}
	return strs
}

// allBytes holds every one of the 256 byte values once, from 0 to 255.
var allBytes = func() string {
	b := make([]byte, 256)
	for i := range b {
		b[i] = byte(i)
	}
	return string(b)
}()

func TestEncodeDecode(t *testing.T) {
	cases := []struct {
		name string
		strs []string
	}{
		{"Example1", []string{"Hello", "World"}},
		{"Example2", []string{""}},
		{"EmptyList", []string{}},
		{"NilList", nil},
		{"TwoEmpty", []string{"", ""}},
		{"ThreeEmpty", []string{"", "", ""}},
		{"EmptyFirst", []string{"", "a"}},
		{"EmptyLast", []string{"a", ""}},
		{"EmptyInTheMiddle", []string{"a", "", "b"}},
		{"OneByte", []string{"x"}},
		{"OneString", []string{"Hello World"}},
		{"Space", []string{" ", "  ", ""}},
		{"Comma", []string{"a,b", ",", ",,"}},
		{"Hash", []string{"#", "##", "a#b", "#a#"}},
		{"Digits", []string{"0", "12", "345", "6789"}},
		{"DigitsAndHash", []string{"3#abc", "1#", "#1", "10#0123456789"}},
		{"LooksLikeALength", []string{"5#Hello", "5#World"}},
		{"Slash", []string{"/", "\\", "\\\\", "a\\#b"}},
		{"Quotes", []string{`"`, `'`, `""`, `","`}},
		{"Brackets", []string{"[", "]", `["a","b"]`}},
		{"Newlines", []string{"\n", "\r\n", "a\nb", "\t"}},
		{"ZeroByte", []string{"\x00", "a\x00b", "\x00\x00"}},
		{"HighBytes", []string{"\xff", "\x80\x81", "\xfe\xff\x00"}},
		{"EveryByte", []string{allBytes}},
		{"EveryByteSplit", []string{allBytes[:128], allBytes[128:]}},
		{"EveryByteAlone", strings.Split(allBytes[:99], "")},
		{"LengthNine", []string{"123456789", "a"}},
		{"LengthTen", []string{"1234567890", "a"}},
		{"LengthNinetyNine", []string{strings.Repeat("9", 99), "a"}},
		{"LengthHundred", []string{strings.Repeat("1", 100), "a"}},
		{"LongestString", []string{strings.Repeat("#", 199)}},
		{"LongAndShort", []string{strings.Repeat("z", 199), "", "q", strings.Repeat("7", 150)}},
		{"SameStringTwice", []string{"dup", "dup"}},
		{"Unicode", []string{"héllo", "世界", "🙂"}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			check(t, c.strs)
		})
	}
}

// TestEncodeDecodeEveryShort walks every list of up to 4 strings from a small
// set of tricky strings. It catches a list that loses or adds an empty
// string, and a string that breaks in two or joins its neighbour.
func TestEncodeDecodeEveryShort(t *testing.T) {
	set := []string{"", "a", "#", "1", "1#", "11", ",", "\x00", "\xff", "\\"}

	for n := 0; n <= 4; n++ {
		idx := make([]int, n)
		for {
			strs := make([]string, n)
			for i, j := range idx {
				strs[i] = set[j]
			}
			if !check(t, strs) {
				return
			}

			// Go to the next list, like a counter in base len(set).
			i := 0
			for i < n && idx[i] == len(set)-1 {
				idx[i] = 0
				i++
			}
			if i == n {
				break
			}
			idx[i]++
		}
	}
}

// TestEncodeDecodeRandom sends random lists through encode and decode. A
// narrow alphabet of digits and marks gives strings that look like parts of
// an encoded string. The full alphabet gives every byte value.
func TestEncodeDecodeRandom(t *testing.T) {
	rnd := rand.New(rand.NewSource(1))

	shapes := []struct {
		name     string
		alphabet string
		maxLen   int
	}{
		{"DigitsAndMarks", "0123456789#,;:|/\\ ", 12},
		{"Letters", "abcdefghijklmnopqrstuvwxyz", 30},
		{"EveryByteShort", allBytes, 5},
		{"EveryByteLong", allBytes, 199},
	}

	for _, s := range shapes {
		t.Run(s.name, func(t *testing.T) {
			for range 500 {
				strs := randomList(rnd, rnd.Intn(100), s.maxLen, s.alphabet)
				if !check(t, strs) {
					return
				}
			}
		})
	}
}

// TestEncodeDecodeNested uses encoded strings as the input strings. The
// strings then hold exactly the bytes that encode writes. It catches a
// solution that cannot tell its own marks from data.
func TestEncodeDecodeNested(t *testing.T) {
	rnd := rand.New(rand.NewSource(2))

	for range 300 {
		var strs []string
		for range 1 + rnd.Intn(4) {
			inner := randomList(rnd, rnd.Intn(5), 8, "0123456789#,\\ab")
			enc := encode(inner)
			if len(enc) >= 200 {
				continue
			}
			strs = append(strs, enc)
			if rnd.Intn(2) == 0 {
				strs = append(strs, "")
			}
		}
		if !check(t, strs) {
			return
		}
	}
}

// TestEncodeDecodeSeparate encodes many lists first, and then decodes the
// strings in a shuffled order. decode runs on another machine, so it can know
// only its argument. It catches a solution that keeps the list in a variable
// between encode and decode.
func TestEncodeDecodeSeparate(t *testing.T) {
	rnd := rand.New(rand.NewSource(3))

	lists := make([][]string, 50)
	encoded := make([]string, len(lists))
	for i := range lists {
		lists[i] = randomList(rnd, rnd.Intn(10), 20, allBytes)
		encoded[i] = encode(slices.Clone(lists[i]))
	}

	order := rnd.Perm(len(lists))
	for _, i := range order {
		if got := decode(encoded[i]); !slices.Equal(got, lists[i]) {
			t.Fatalf("decode(encode(%q)) = %q after other calls of encode, want %q", lists[i], got, lists[i])
		}
	}
}

// TestEncodeDecodeCarriesTheData gives random bytes, which no encoding can
// make shorter. The encoded string must hold at least as many bytes as the
// strings together. It catches a solution that sends a short key and keeps
// the data on the side of encode.
func TestEncodeDecodeCarriesTheData(t *testing.T) {
	rnd := rand.New(rand.NewSource(4))

	for range 20 {
		strs := randomList(rnd, 1+rnd.Intn(99), 199, allBytes)
		total := 0
		for _, s := range strs {
			total += len(s)
		}

		if enc := encode(strs); len(enc) < total {
			t.Fatalf("encode of %d strings with %d bytes in total gave only %d bytes, so the encoded string does not carry the data", len(strs), total, len(enc))
		}
	}
}

// TestEncodeDecodeLarge builds lists of the largest size that the
// constraints allow: 99 strings of 199 bytes. The errors do not print the
// list, because it is too long.
func TestEncodeDecodeLarge(t *testing.T) {
	rnd := rand.New(rand.NewSource(5))

	run := func(t *testing.T, strs []string) {
		t.Helper()

		got := decode(encode(strs))
		if len(got) != len(strs) {
			t.Fatalf("decode(encode(list of %d strings)) gave %d strings", len(strs), len(got))
		}
		for i := range strs {
			if got[i] != strs[i] {
				t.Fatalf("decode(encode(list of %d strings)): string %d = %q, want %q", len(strs), i, got[i], strs[i])
			}
		}
	}

	t.Run("RandomBytes", func(t *testing.T) {
		run(t, randomList(rnd, 99, 199, allBytes))
	})

	t.Run("AllLongest", func(t *testing.T) {
		strs := make([]string, 99)
		for i := range strs {
			b := make([]byte, 199)
			for j := range b {
				b[j] = allBytes[rnd.Intn(256)]
			}
			strs[i] = string(b)
		}
		run(t, strs)
	})

	for _, fill := range []string{"#", "9", "\x00", "\xff", ","} {
		t.Run("AllSame", func(t *testing.T) {
			strs := make([]string, 99)
			for i := range strs {
				strs[i] = strings.Repeat(fill, 199)
			}
			run(t, strs)
		})
	}

	t.Run("AllEmpty", func(t *testing.T) {
		run(t, make([]string, 99))
	})
}

func BenchmarkEncodeDecodeExample(b *testing.B) {
	strs := []string{"Hello", "World"}

	for b.Loop() {
		decode(encode(strs))
	}
}

// BenchmarkEncodeDecodeLarge uses the largest list: 99 strings of 199 random
// bytes.
func BenchmarkEncodeDecodeLarge(b *testing.B) {
	rnd := rand.New(rand.NewSource(6))
	strs := make([]string, 99)
	for i := range strs {
		buf := make([]byte, 199)
		rnd.Read(buf)
		strs[i] = string(buf)
	}

	for b.Loop() {
		decode(encode(strs))
	}
}
