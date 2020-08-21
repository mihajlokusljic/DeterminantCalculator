package determinantCalc

func removeColumn(cols []int, columnToRemove int) []int {
	n := len(cols)
	result := make([]int, n-1)
	currentIndex := 0

	for _, col := range cols {
		if col != columnToRemove {
			result[currentIndex] = col
			currentIndex++
		}
	}

	return result
}
