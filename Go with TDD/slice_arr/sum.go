package slice_arr

func Sum(numbers []int) int {
	sum := 0
	for _, number := range numbers {
		sum += number
	}
	return sum
}

// this is for old stuff basically the old way we used to write loops in go
// for i := 0; i < 5; i++ {
//		sum += numbers[i]
//	}
// The _ is the blank identifier.
// make sure that
// // Using _ tells Go: "I know there's a value here (the index), but I'm ignoring it intentionally."
// tells  just about its number and not index since we dont need it atm
//

func SumAll(numbersTosum ...[]int) []int {
	var sums []int

	for _, numbers := range numbersTosum {
		sums = append(sums, Sum(numbers))
	}

	return sums
}

func SumAllTails(numberToSum ...[]int) []int {

	var sums []int
	for _, numbers := range numberToSum {
		tail := numbers[1:]
		sums = append(sums, Sum(tail))
	}
	return sums
}
