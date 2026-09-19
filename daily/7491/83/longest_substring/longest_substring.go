package longestsubstring

func lengthOfLongestSubstringNaive(s string) int {
	best := 0
	var count int
	for i, char1 := range s {
		seen := make(map[rune]bool)
		seen[char1] = true
		count = 1
		for _, char2 := range s[i+1:] {
			if _, ok := seen[char2]; ok {
				if count > best {
					best = count
				}
				break
			}
			seen[char2] = true
			count++
		}
		if count > best {
			best = count
		}
	}
	return best
}

func lengthOfLongestSubstring(s string) int {
	seen := make(map[rune]int)
	maxLen := 0
	left := 0

	for right, char := range s {

		if idx, found := seen[char]; found && idx >= left {
			left = idx + 1
		}

		seen[char] = right

		windowLen := right - left + 1
		if windowLen > maxLen {
			maxLen = windowLen
		}
	}
	return maxLen
}
