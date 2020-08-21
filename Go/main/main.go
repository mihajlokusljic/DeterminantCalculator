package main

import (
	"fmt"
	"os"
	"time"

	"gitlab.com/mihajlokusljic/DeterminantCalculator/IO"
	"gitlab.com/mihajlokusljic/DeterminantCalculator/determinantCalc"
)

func main() {
	matrixPaths := os.Args[1 : len(os.Args)-1]
	resultsFilePath := os.Args[len(os.Args)-1]
	var results []IO.ExecutionResult

	// process matrixes
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
