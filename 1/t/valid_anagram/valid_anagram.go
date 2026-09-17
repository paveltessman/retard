package validanagram

func isAnagram(s string, t string) bool {
	if len(s) != len(t) {
		return false
	}

	sMap := make(map[rune]int) // map[char]count
	for _, char := range s {
		sMap[char]++
	}

	tMap := make(map[rune]int)
	for _, char := range t {
		tMap[char]++
	}

	for char, sCount := range sMap { // okay to iterate randomly
		tCount, ok := tMap[char]
		if !ok {
			return false
		}

		if sCount != tCount {
			return false
		}
	}
	return true
}
