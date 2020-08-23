package determinantCalc

// removeColumn is udes to determine columns in a submatrix for a given minor.
// For a given slice of column indexes return a new slice in which the specified column
// index is removed.
//
// Args:
// 	cols: column indexes of a matrix
//	columnToRemove: column index which cshould be removed
//
// Return:
//  a slice of column indexes where the specified value has been removed
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
