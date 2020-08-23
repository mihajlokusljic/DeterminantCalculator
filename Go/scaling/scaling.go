package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"gitlab.com/mihajlokusljic/DeterminantCalculator/IO"
	"gitlab.com/mihajlokusljic/DeterminantCalculator/determinantCalc"
)

const matrixPathTemplate = "../../test_data/matrica%dx%d.txt"
const resultsBasePath = "../../results"
const strongScalingMatrixOrder = 10

var availableMatrixOrders = [...]int{3, 5, 8, 9, 10, 11}

// ScalingResult holds information about parallel determinant calculation to be stored in a results file.
//
// Attributes:
//     MatrixOrder: order of matrix for which the determinant is calculated
//     TasksNum: number of processes used during calculation
//     ExecTimeMs: duration of the calculation in milliseconds
type ScalingResult struct {
	MatrixOrder int
	TasksNum    int
	ExecTimeMs  int64
}

// Populates a given CSV file with given results. If file does not exists it will be created.
// The header contains: matrix_order (order of matrix for which the determinant is calculated),
// number_of_tasks (number of processes used in calculation), exec_time_ms (duration of the calculation
// in milliseconds).
//
// Args:
// 	resultsFilePath: path of the CSV file to hold results
// 	scalingResults: information about different determinant calculation attempts
//
// Return:
// 	nil
func writeScalingResults(resultsFilePath string, scalingResults []ScalingResult) {
	resultsFile, err := os.OpenFile(resultsFilePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Unable to access results file.\n", err)
	}
	defer resultsFile.Close()

	resultsFile.WriteString("matrix_order,number_of_tasks,exec_time_ms\n")
	for _, res := range scalingResults {
		resultsFile.WriteString(fmt.Sprintf("%d,%d,%d\n",
			res.MatrixOrder, res.TasksNum, res.ExecTimeMs))
	}
}

// Loads a test matrix of given order from a predefined directory.
//
// Args:
//     matrixOrder: order of desired matrix
//
// Return:
//     matrix loaded into memory
func loadTestMatrix(matrixOrder int) [][]float64 {
	var matrixFilePath string = fmt.Sprintf(matrixPathTemplate, matrixOrder, matrixOrder)
	var matrix [][]float64 = IO.ReadMatrix(matrixFilePath)
	fmt.Printf("Successfully loaded matrix of order %d. Calculating determinant(s)...\n", len(matrix))
	return matrix
}

// Loads a predefined matrix for strong scaling and calculates the determinant
// using parallel implementation with number of tasks ranging from one to the
// order of the matrix. Statistics about calculations are printed onto the console
// and written into a predefined results file: strong_scaling_results_go.csv.
//
// Return:
// 	nil
func strongScalingExperiment() {
	fmt.Println("Strong scaling experiment:\n")
	var matrix [][]float64 = loadTestMatrix(strongScalingMatrixOrder)
	var results []ScalingResult

	for tasksNum := 1; tasksNum <= strongScalingMatrixOrder; tasksNum++ {
		start := time.Now()
		determinantCalc.DetParallel(matrix, tasksNum)
		elapsed := time.Since(start)
		duration := elapsed.Milliseconds()
		scalingResult := ScalingResult{
			ExecTimeMs:  duration,
			MatrixOrder: strongScalingMatrixOrder,
			TasksNum:    tasksNum,
		}
		results = append(results, scalingResult)

		fmt.Printf("Parallel determinant calculation of matrix of order %d with %d tasks took %d ms.\n",
			strongScalingMatrixOrder, tasksNum, duration)
	}

	writeScalingResults(fmt.Sprintf("%s/strong_scaling_results_go.csv", resultsBasePath), results)
}

// Loads matrices of different orders from predefined files and calculates
// their determinants using parallel implementation where the number of tasks is
// equal to the order of matrix. Statistics about calculations are printed onto the console
// and written into a predefined results file: weak_scaling_results_go.csv.
//
// Return:
// 	nil
func weakScalingExperiment() {
	fmt.Println("\nWeak scaling experiment:\n")
	var results []ScalingResult

	for _, n := range availableMatrixOrders {
		matrix := loadTestMatrix(n)

		start := time.Now()
		determinantCalc.DetParallel(matrix, n)
		elapsed := time.Since(start)
		duration := elapsed.Milliseconds()
		scalingResult := ScalingResult{
			ExecTimeMs:  duration,
			MatrixOrder: n,
			TasksNum:    n,
		}
		results = append(results, scalingResult)

		fmt.Printf("Parallel determinant calculation of matrix of order %d with %d tasks took %d ms.\n", n, n, duration)
	}

	writeScalingResults(fmt.Sprintf("%s/weak_scaling_results_go.csv", resultsBasePath), results)
}

func main() {
	strongScalingExperiment()
	weakScalingExperiment()
}
