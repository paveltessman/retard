package productexceptself

import (
	"math/rand"
	"slices"
	"testing"
)

// brute gives the answer of the task in the slow and plain way: for every
// position it walks the whole array again and multiplies the other numbers.
// It never divides, so it also holds for an array with a zero. The tests
// compare the solution with it.
func brute(nums []int) []int {
	answer := make([]int, len(nums))

	for i := range nums {
		product := 1
		for j, v := range nums {
			if j != i {
				product *= v
			}
		}
		answer[i] = product
	}

	return answer
}

// check calls productExceptSelf on a copy of nums and compares the answer with
// want. It also reads the copy back, because productExceptSelf must leave the
// array as it was.
func check(t *testing.T, nums, want []int) {
	t.Helper()

	got := slices.Clone(nums)

	if answer := productExceptSelf(got); !slices.Equal(answer, want) {
		t.Errorf("productExceptSelf(%v) = %v, want %v", nums, answer, want)
	}
	if !slices.Equal(got, nums) {
		t.Errorf("productExceptSelf(%v) changed the array to %v", nums, got)
	}
}

// everyArray calls f once for every array of n numbers taken from vals.
func everyArray(n int, vals []int, f func([]int)) {
	nums := make([]int, n)

	var build func(i int)
	build = func(i int) {
		if i == n {
			f(nums)
			return
		}
		for _, v := range vals {
			nums[i] = v
			build(i + 1)
		}
	}

	build(0)
}

// randomNums gives n numbers between -max and max.
func randomNums(n, max int, rnd *rand.Rand) []int {
	nums := make([]int, n)
	for i := range nums {
		nums[i] = rnd.Intn(2*max+1) - max
	}

	return nums
}

// repeat gives an array of n copies of v.
func repeat(n, v int) []int {
	nums := make([]int, n)
	for i := range nums {
		nums[i] = v
	}

	return nums
}

func TestProductExceptSelf(t *testing.T) {
	cases := []struct {
		name string
		nums []int
		want []int
	}{
		{"Example1", []int{1, 2, 3, 4}, []int{24, 12, 8, 6}},
		{"Example2", []int{-1, 1, 0, -3, 3}, []int{0, 0, 9, 0, 0}},
		{"TwoNumbers", []int{2, 3}, []int{3, 2}},
		{"TwoEqualNumbers", []int{5, 5}, []int{5, 5}},
		{"TwoOnes", []int{1, 1}, []int{1, 1}},
		{"TwoZeros", []int{0, 0}, []int{0, 0}},
		{"OneZeroOfTwo", []int{0, 4}, []int{4, 0}},
		{"OneZeroAtTheStart", []int{0, 1, 2, 3}, []int{6, 0, 0, 0}},
		{"OneZeroAtTheEnd", []int{1, 2, 3, 0}, []int{0, 0, 0, 6}},
		{"OneZeroInTheMiddle", []int{1, 2, 0, 3, 4}, []int{0, 0, 24, 0, 0}},
		{"TwoZerosApart", []int{1, 0, 3, 0, 5}, []int{0, 0, 0, 0, 0}},
		{"TwoZerosTogether", []int{2, 0, 0, 3}, []int{0, 0, 0, 0}},
		{"ThreeZeros", []int{0, 0, 0}, []int{0, 0, 0}},
		{"AllZerosButOne", []int{0, 0, 7, 0}, []int{0, 0, 0, 0}},
		{"AllNegative", []int{-1, -2, -3, -4}, []int{-24, -12, -8, -6}},
		{"OneNegative", []int{-2, 3, 4}, []int{12, -8, -6}},
		{"TwoNegatives", []int{-2, -3, 4}, []int{-12, -8, 6}},
		{"ThreeNegatives", []int{-2, -3, -4}, []int{12, 8, 6}},
		{"MinusOnes", []int{-1, -1, -1}, []int{1, 1, 1}},
		{"MinusOnesEvenCount", []int{-1, -1, -1, -1}, []int{-1, -1, -1, -1}},
		{"Ones", []int{1, 1, 1, 1, 1}, []int{1, 1, 1, 1, 1}},
		{"AlternatingSigns", []int{1, -1, 1, -1}, []int{1, -1, 1, -1}},
		{"WithOnes", []int{1, 3, 1, 2}, []int{6, 2, 6, 3}},
		{"LargestNumbers", []int{30, 30, 30, 30}, []int{27000, 27000, 27000, 27000}},
		{"SmallestNumbers", []int{-30, -30, -30}, []int{900, 900, 900}},
		{"MixedLargest", []int{30, -30, 30, -30}, []int{27000, -27000, 27000, -27000}},
		{"LongArrayOfOnes", repeat(50, 1), repeat(50, 1)},
		{"NearTheLimit", []int{30, 30, 30, 30, 30, -1}, []int{-810000, -810000, -810000, -810000, -810000, 24300000}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			check(t, c.nums, c.want)
		})
	}
}

// TestProductExceptSelfMatchesBrute walks the cases of the table again and
// compares the wanted answers with brute. It guards the table itself.
func TestProductExceptSelfMatchesBrute(t *testing.T) {
	arrays := [][]int{
		{1, 2, 3, 4},
		{-1, 1, 0, -3, 3},
		{0, 0},
		{30, 30, 30, 30},
		{1, 0, 3, 0, 5},
	}

	for _, nums := range arrays {
		want := brute(nums)

		got := productExceptSelf(slices.Clone(nums))
		if !slices.Equal(got, want) {
			t.Errorf("productExceptSelf(%v) = %v, want %v", nums, got, want)
		}
	}
}

// TestProductExceptSelfEveryShortArray walks every array of 2 to 6 numbers
// taken from -2, -1, 0, 1 and 3. It catches a solution that divides, because
// the arrays hold one zero and two zeros. It also catches a lost sign and a
// lost first or last position.
func TestProductExceptSelfEveryShortArray(t *testing.T) {
	const maxLen = 6

	vals := []int{-2, -1, 0, 1, 3}

	for n := 2; n <= maxLen; n++ {
		failed := false

		everyArray(n, vals, func(nums []int) {
			if failed {
				return
			}

			got := productExceptSelf(slices.Clone(nums))

			if want := brute(nums); !slices.Equal(got, want) {
				t.Errorf("productExceptSelf(%v) = %v, want %v", nums, got, want)
				failed = true
			}
		})
	}
}

// TestProductExceptSelfRandom compares the solution with brute on random
// arrays. The range of the numbers stays small, so the products stay inside a
// 32-bit integer.
func TestProductExceptSelfRandom(t *testing.T) {
	rnd := rand.New(rand.NewSource(1))

	for range 500 {
		nums := randomNums(2+rnd.Intn(14), 3, rnd)

		got := productExceptSelf(slices.Clone(nums))

		if want := brute(nums); !slices.Equal(got, want) {
			t.Fatalf("productExceptSelf(%v) = %v, want %v", nums, got, want)
		}
	}
}

// TestProductExceptSelfRandomWithoutZeros compares the solution with brute on
// arrays of -1 and 1 only. Here every answer is -1 or 1, so a wrong sign shows
// up at once.
func TestProductExceptSelfRandomWithoutZeros(t *testing.T) {
	rnd := rand.New(rand.NewSource(2))

	for range 500 {
		n := 2 + rnd.Intn(30)

		nums := make([]int, n)
		for i := range nums {
			nums[i] = 2*rnd.Intn(2) - 1
		}

		got := productExceptSelf(slices.Clone(nums))

		if want := brute(nums); !slices.Equal(got, want) {
			t.Fatalf("productExceptSelf(%v) = %v, want %v", nums, got, want)
		}
	}
}

// TestProductExceptSelfOneZeroEverywhere puts one zero at every position of an
// array of numbers above 1. Only the position of the zero holds a number above
// 0. A solution that divides fails here.
func TestProductExceptSelfOneZeroEverywhere(t *testing.T) {
	const n = 10

	base := make([]int, n)
	for i := range base {
		base[i] = 2 + i%3
	}

	for at := range n {
		nums := slices.Clone(base)
		nums[at] = 0

		want := make([]int, n)
		want[at] = 1
		for i, v := range base {
			if i != at {
				want[at] *= v
			}
		}

		check(t, nums, want)
	}
}

// TestProductExceptSelfTwoZerosEverywhere puts two zeros in an array of
// numbers above 1. Every answer is 0. A solution that counts the zeros wrong
// fails here.
func TestProductExceptSelfTwoZerosEverywhere(t *testing.T) {
	const n = 8

	for a := range n {
		for b := a + 1; b < n; b++ {
			nums := repeat(n, 2)
			nums[a] = 0
			nums[b] = 0

			check(t, nums, make([]int, n))
		}
	}
}

// TestProductExceptSelfKeepsNeighbours gives a part of a larger array to
// productExceptSelf. A solution that reads or writes outside of nums fails
// here.
func TestProductExceptSelfKeepsNeighbours(t *testing.T) {
	arr := []int{7, 7, 1, 2, 3, 7, 7}
	nums := arr[2:5]
	before := slices.Clone(arr)

	if got, want := productExceptSelf(nums), []int{6, 3, 2}; !slices.Equal(got, want) {
		t.Errorf("productExceptSelf(%v) = %v, want %v", nums, got, want)
	}
	if !slices.Equal(arr, before) {
		t.Errorf("productExceptSelf(%v) changed the array around it to %v", nums, arr)
	}
}

// TestProductExceptSelfKeepsSpareRoom gives a slice that has room after its
// last number. A solution that appends to nums writes into that room, and the
// test finds it.
func TestProductExceptSelfKeepsSpareRoom(t *testing.T) {
	arr := make([]int, 3, 10)
	copy(arr, []int{1, 2, 3})

	for i := 3; i < cap(arr); i++ {
		arr[:cap(arr)][i] = 99
	}

	if got, want := productExceptSelf(arr), []int{6, 3, 2}; !slices.Equal(got, want) {
		t.Errorf("productExceptSelf(%v) = %v, want %v", arr, got, want)
	}

	room := arr[:cap(arr)][3:]
	if want := repeat(len(room), 99); !slices.Equal(room, want) {
		t.Errorf("productExceptSelf wrote %v into the room after the array", room)
	}
}

// TestProductExceptSelfLength reads the length of the answer for arrays of
// many lengths. The answer must hold one number for every number of nums.
func TestProductExceptSelfLength(t *testing.T) {
	for n := 2; n <= 40; n++ {
		nums := repeat(n, 1)

		if got := productExceptSelf(nums); len(got) != n {
			t.Errorf("productExceptSelf(%d ones) gave %d numbers, want %d", n, len(got), n)
		}
	}
}

// TestProductExceptSelfLarge uses the largest array of the constraints. A
// solution that walks the array again for every position needs much more time
// here. The numbers stay at -1, 0 and 1, so the products stay inside a 32-bit
// integer.
func TestProductExceptSelfLarge(t *testing.T) {
	const n = 100000

	t.Run("AllOnes", func(t *testing.T) {
		if got, want := productExceptSelf(repeat(n, 1)), repeat(n, 1); !slices.Equal(got, want) {
			t.Errorf("productExceptSelf(%d ones) gave a wrong answer", n)
		}
	})

	t.Run("EvenCountOfMinusOnes", func(t *testing.T) {
		// The product of all the numbers is 1, so every answer is -1.
		if got, want := productExceptSelf(repeat(n, -1)), repeat(n, -1); !slices.Equal(got, want) {
			t.Errorf("productExceptSelf(%d numbers of -1) gave a wrong answer", n)
		}
	})

	t.Run("OddCountOfMinusOnes", func(t *testing.T) {
		// The product of all the numbers is -1, so every answer is 1.
		if got, want := productExceptSelf(repeat(n-1, -1)), repeat(n-1, 1); !slices.Equal(got, want) {
			t.Errorf("productExceptSelf(%d numbers of -1) gave a wrong answer", n-1)
		}
	})

	t.Run("OneZero", func(t *testing.T) {
		// Only the position of the zero holds 1. The other numbers are 1, so
		// the product of the rest is 1.
		nums := repeat(n, 1)
		nums[n/3] = 0

		want := make([]int, n)
		want[n/3] = 1

		if got := productExceptSelf(nums); !slices.Equal(got, want) {
			t.Errorf("productExceptSelf(%d ones and one zero) gave a wrong answer", n)
		}
	})

	t.Run("TwoZeros", func(t *testing.T) {
		nums := repeat(n, 1)
		nums[0] = 0
		nums[n-1] = 0

		if got, want := productExceptSelf(nums), make([]int, n); !slices.Equal(got, want) {
			t.Errorf("productExceptSelf(%d ones and two zeros) gave a wrong answer", n)
		}
	})

	t.Run("AllZeros", func(t *testing.T) {
		if got, want := productExceptSelf(make([]int, n)), make([]int, n); !slices.Equal(got, want) {
			t.Errorf("productExceptSelf(%d zeros) gave a wrong answer", n)
		}
	})

	t.Run("Alternating", func(t *testing.T) {
		// The count of the numbers of -1 is even, so the product of all the
		// numbers is 1. Every answer holds the number of its own position.
		nums := make([]int, n)
		for i := range nums {
			if i%2 == 0 {
				nums[i] = 1
			} else {
				nums[i] = -1
			}
		}

		want := slices.Clone(nums)

		if got := productExceptSelf(nums); !slices.Equal(got, want) {
			t.Errorf("productExceptSelf(%d alternating numbers) gave a wrong answer", n)
		}
	})
}

// TestProductExceptSelfMedium compares the solution with brute on an array
// that brute still handles. It joins a run of ones, a run of equal numbers and
// a random part.
func TestProductExceptSelfMedium(t *testing.T) {
	const n = 2000

	rnd := rand.New(rand.NewSource(3))

	nums := make([]int, 0, n)
	nums = append(nums, repeat(n/4, 1)...)
	nums = append(nums, repeat(n/4, -1)...)
	for i := range n / 4 {
		if i%2 == 0 {
			nums = append(nums, 1)
		} else {
			nums = append(nums, -1)
		}
	}
	for range n / 4 {
		nums = append(nums, 2*rnd.Intn(2)-1)
	}

	got := productExceptSelf(slices.Clone(nums))

	if want := brute(nums); !slices.Equal(got, want) {
		t.Errorf("productExceptSelf(the mixed array) gave a wrong answer")
	}
}

func BenchmarkProductExceptSelfSmall(b *testing.B) {
	nums := []int{1, 2, 3, 4}

	for b.Loop() {
		productExceptSelf(nums)
	}
}

// BenchmarkProductExceptSelfLarge gives the largest array of the constraints.
func BenchmarkProductExceptSelfLarge(b *testing.B) {
	const n = 100000

	nums := make([]int, n)
	for i := range nums {
		nums[i] = 2*(i%2) - 1
	}

	b.ResetTimer()
	for b.Loop() {
		productExceptSelf(nums)
	}
}

// BenchmarkProductExceptSelfZeros gives an array of zeros. A solution that
// counts the zeros takes another path here.
func BenchmarkProductExceptSelfZeros(b *testing.B) {
	const n = 100000

	nums := make([]int, n)

	b.ResetTimer()
	for b.Loop() {
		productExceptSelf(nums)
	}
}
