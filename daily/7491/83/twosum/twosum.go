package twosum

import (
	"slices"
)

func twosumNaive(nums []int, target int) []int {
	for i := range len(nums) {
		for j := i + 1; j < len(nums); j++ {
			if nums[i]+nums[j] == target {
				return []int{i, j}
			}
		}
	}
	panic("fuk")
}

func twosumBetterButStill(nums []int, target int) []int {
	members := slices.Clone(nums)
	slices.Sort(members)

	for i, first := range nums {
		second := target - first

		_, ok := slices.BinarySearch(members, second)
		if !ok {
			continue
		}

		for j, val := range nums {
			if i != j && val == second {
				return []int{i, j}

			}
		}

	}
	panic("fuk2")
}

func twosum(nums []int, target int) []int {
	seen := make(map[int]int) // map[value]index

	for i, first := range nums {
		second := target - first
		j, ok := seen[second]
		seen[first] = i
		if !ok {
			continue
		}
		return []int{i, j}
	}
	panic("fuk")
}
