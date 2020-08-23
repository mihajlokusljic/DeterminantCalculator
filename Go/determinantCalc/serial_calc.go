package determinantCalc

// minorCalcSerial calculates an arbitrary minor of given matrix using Laplace expansion.
// The original matrix is accessed for calculation (only the necessary rows and columns).
//
// Args:
// 	matrix: matrix containing the submatrix of the minor
// 	begin_row_index: index of the first row of the submatrix for the minor, in the original matrix
// 	column_indexes: indexes of columns of the submatrix for the minor, in the original matrix.
//
// Return:
// 	value of the given minor
func minorCalcSerial(matrix [][]float64, beginRowIndex int, columnIndexes []int) float64 {
	n := len(columnIndexes) // order of submatrix

	if n == 1 {
		return matrix[beginRowIndex][columnIndexes[0]] // return the only element in the submatrix
	}

	var result float64 = 0.0
	minors := make([]float64, n)
	sgn := 1.0

	// calculate the minors
	for idx, col := range columnIndexes {
		minorCols := removeColumn(columnIndexes, col)
		minors[idx] = minorCalcSerial(matrix, beginRowIndex+1, minorCols)
	}

	// calculate the determinant using values of minors
	for j := 0; j < n; j++ {
		result += sgn * matrix[beginRowIndex][columnIndexes[j]] * minors[j]
		sgn *= -1
	}

	return result
}

// DetSerial calculates and returns the determinant of given matrix using serial implementation.
//
// Args:
//     matrix: matrix for which the determinant is calculated.
//
// Return:
//     value of the determinant
func DetSerial(matrix [][]float64) float64 {
	n := len(matrix)
	cols := make([]int, n)
	for i := 0; i < n; i++ {
		cols[i] = i
	}

	return minorCalcSerial(matrix, 0, cols)
}
