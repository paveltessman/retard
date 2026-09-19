package groupanagrams

import (
	"slices"
)

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

func groupAnagramsNaive(strs []string) [][]string {
	taken := make([]bool, len(strs))
	groups := make(map[string][]string)

	for i := range slices.Backward(strs) {
		if taken[i] {
			continue
		}
		groups[strs[i]] = append(groups[strs[i]], strs[i])
		for j := range i {
			if taken[j] {
				continue
			}
			if isAnagram(strs[i], strs[j]) {
				taken[j] = true
				groups[strs[i]] = append(groups[strs[i]], strs[j])
			}
		}
	}

	result := make([][]string, 0, len(groups))
	for _, group := range groups { // okay to go random
		result = append(result, group)
	}
	return result
}

func groupAnagrams(strs []string) [][]string {
	groups := make(map[string][]string)

	for _, word := range strs {
		array := []byte(word)
		slices.Sort(array)
		key := string(array)
		groups[key] = append(groups[key], word)
	}

	result := make([][]string, 0, len(groups))
	for _, group := range groups {
		result = append(result, group)
	}
	return result
}
