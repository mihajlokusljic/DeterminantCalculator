package main

import (
	"fmt"
	"os"
	"time"

	"gitlab.com/mihajlokusljic/DeterminantCalculator/IO"
	"gitlab.com/mihajlokusljic/DeterminantCalculator/determinantCalc"
)

// main method loads matrices from files specified in arguments to the program,
// calculates the determinant of each matrix, both sequentially and parallel, and appends data about
// the execution times to a results file, specified as the last argument to the program. Results
// and execution times are also printed to the console.
// The first expected argument is the program itself. The next n arguments should be paths to
// valid text files containing a regular matrix: first row contains a number (n) which specifies the
// order of the matrix. The next n rows contain exactly n real numbers, separated by one or more
// space characters. The last argument to the program is a path to a csv file where information
// about the execution times for different matrices are stored.
//
// Return:
// 	None
func main() {
	matrixPaths := os.Args[1 : len(os.Args)-1]
	resultsFilePath := os.Args[len(os.Args)-1]
	var results []IO.ExecutionResult

	// process matrices
	for _, matrixFilePath := range matrixPaths {
		matrix := IO.ReadMatrix(matrixFilePath)
		fmt.Printf("\nSuccessfully loaded matrix: %s.\nCalculating determinants...\n\n", matrixFilePath)
		results = append(results, serialCalculation(matrix))
		results = append(results, parallelCalculation(matrix))

		if len(matrixPaths) > 1 {
			fmt.Println("====================")
		}
	}

	IO.WriteExecutionResults(results, resultsFilePath)
}

// serialCalculation calculates the determinant of given regular matrix using serial implementation
// and measures the execution time.
//
// Args:
//     matrix: regular matrix for which the determinant is calculated
//
// Return:
//     result object that can be stored in the results file
func serialCalculation(matrix [][]float64) IO.ExecutionResult {
	start := time.Now()
	var det float64 = determinantCalc.DetSerial(matrix)
	elapsed := time.Since(start)
	duration := elapsed.Milliseconds()
	fmt.Printf("Serial calculation result:\ndet(mat) = %f\n", det)
	fmt.Printf("Serial calculation took %d ms.\n\n", duration)

	return IO.ExecutionResult{
		MatrixOrder: len(matrix),
		ExecTimeMs:  duration,
		Serial:      true,
	}
}

// parallelCalculation calculates the determinant od given regular matrix using parallel implementation
// and measures the execution time. Number of tasks used to calculate the determinant is qeual to
// order of the matrix.
//
// Args:
//     matrix: regular matrix for which the determinant is calculated
//
// Return:
//     result object that can be stored in the results file
func parallelCalculation(matrix [][]float64) IO.ExecutionResult {
	start := time.Now()
	var det float64 = determinantCalc.DetParallel(matrix, len(matrix))
	elapsed := time.Since(start)
	duration := elapsed.Milliseconds()
	fmt.Printf("Parallel calculation result:\ndet(mat) = %f\n", det)
	fmt.Printf("Parallel calculation took %d ms.\n\n", duration)

	return IO.ExecutionResult{
		MatrixOrder: len(matrix),
		ExecTimeMs:  duration,
		Serial:      false,
	}
}
