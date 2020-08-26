package determinantCalc

import "time"

// minorCalcSerial calculates an arbitrary minor of given matrix using Laplace expansion.
// The original matrix is accessed for calculation (only the necessary rows and columns).
// Returns two values:
//     1. the value of the minor
//     2. time in milliseconds spent executing code that can be parallelized if argument
//     measureParallelCode is true, and zero otherwise (it will not be measured)
//
// Args:
// 	matrix: matrix containing the submatrix of the minor
// 	beginRowIndex: index of the first row of the submatrix for the minor, in the original matrix
// 	columnIndexes: indexes of columns of the submatrix for the minor, in the original matrix.
//  measureParallelCode: indicates whether the execution time for code that can be parallelized should be measured
//
// Return:
// 	value of the given minor, time spent executing code that can be parallelized
func minorCalcSerial(matrix [][]float64, beginRowIndex int, columnIndexes []int, measureParallelCode bool) (float64, int64) {
	n := len(columnIndexes) // order of submatrix

	if n == 1 {
		return matrix[beginRowIndex][columnIndexes[0]], 0 // return the only element in the submatrix
	}

	var result float64 = 0.0
	minors := make([]float64, n)
	sgn := 1.0
	var parallelCodeExecTime int64 = 0

	// calculate the minors
	for idx, col := range columnIndexes {
		minorCols := removeColumn(columnIndexes, col)

		if measureParallelCode {
			start := time.Now()
			minor, _ := minorCalcSerial(matrix, beginRowIndex+1, minorCols, false)
			minors[idx] = minor
			elapsed := time.Since(start)
			minorCalcTimeMs := elapsed.Milliseconds()
			parallelCodeExecTime += minorCalcTimeMs
		} else {
			minor, _ := minorCalcSerial(matrix, beginRowIndex+1, minorCols, false)
			minors[idx] = minor
		}
	}

	// calculate the determinant using values of minors
	for j := 0; j < n; j++ {
		result += sgn * matrix[beginRowIndex][columnIndexes[j]] * minors[j]
		sgn *= -1
	}

	return result, parallelCodeExecTime
}

// DetSerial calculates and returns the determinant of given matrix using serial implementation, as well
// as time in milliseconds, spent executing code that can be parallelized.
//
// Args:
//     matrix: matrix for which the determinant is calculated.
//
// Return:
//     value of the determinant, time spent executing code that can be parallelized
func DetSerial(matrix [][]float64) (float64, int64) {
	n := len(matrix)
	cols := make([]int, n)
	for i := 0; i < n; i++ {
		cols[i] = i
	}

	return minorCalcSerial(matrix, 0, cols, true)
}
