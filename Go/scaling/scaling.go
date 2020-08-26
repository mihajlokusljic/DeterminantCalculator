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

var serialCodeShare float64 = 0
var parallelCodeShare float64 = 1

// ScalingResult holds information about determinant calculation for matrix of certain order to be stored in a results file.
//
// Attributes:
// 	MatrixOrder: order of matrix for which the determinant is calculated
// 	SerialExecTimeMs: execution time of serial implementation in milliseconds
// 	ParallelTasksNum: number of processes used in parallel calculation)
// 	ParallelExecTimeMs: execution time of parallel implementation in milliseconds
// 	AchievedSpeedup: actual speedup calculated using serial and parallel execution time
// 	MaxTheoreticalSpeedup: maximum possible speedup considering serial and parallel portions of code
type ScalingResult struct {
	MatrixOrder           int
	SerialExecTimeMs      int64
	ParallelTasksNum      int
	ParallelExecTimeMs    int64
	AchievedSpeedup       float64
	MaxTheoreticalSpeedup float64
}

// Populates a given CSV file with given results. If file does not exists it will be created.
// The header contains: matrix_order (order of matrix for which the determinant is calculated),
// serial_exec_time_ms (execution time of serial implementation in milliseconds),
// parallel_tasks_num (number of processes used in parallel calculation), parallel_exec_time_ms
// (execution time of parallel implementation in milliseconds), achieved_speedup (actual speedup calculated using
// serial and parallel execution time), max_theoretical_speedup (maximum possible speedup considering serial and
// parallel portions of code).
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

	resultsFile.WriteString("matrix_order,serial_exec_time_ms,parallel_tasks_num,parallel_exec_time_ms,achieved_speedup,max_theoretical_speedup\n")
	for _, res := range scalingResults {
		resultsFile.WriteString(fmt.Sprintf("%d,%d,%d,%d,%f,%f\n",
			res.MatrixOrder, res.SerialExecTimeMs, res.ParallelTasksNum, res.ParallelExecTimeMs,
			res.AchievedSpeedup, res.MaxTheoreticalSpeedup))
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

func maxSpeedupAmdahl(processesNum int) float64 {
	return 1.0 / (serialCodeShare + parallelCodeShare/float64(processesNum))
}

func maxSpeedupGustafson(processesNum int) float64 {
	return serialCodeShare + parallelCodeShare*float64(processesNum)
}

// Loads a predefined matrix for strong scaling and calculates the determinant
// using parallel implementation with number of tasks ranging from one to the
// order of the matrix. Statistics about calculations are printed onto the console
// and written into a predefined results file: strong_scaling_results_go.csv.
//
// Return:
// 	nil
func strongScalingExperiment() {
	fmt.Println("==========================")
	fmt.Println("Strong scaling experiment:")
	fmt.Println("==========================\n\n")

	var matrix [][]float64 = loadTestMatrix(strongScalingMatrixOrder)
	var results []ScalingResult

	start := time.Now()
	_, potentialParallelCodeExecTimeMs := determinantCalc.DetSerial(matrix)
	elapsed := time.Since(start)
	serialExecTimeMs := elapsed.Milliseconds()
	parallelCodeShare = float64(potentialParallelCodeExecTimeMs) / float64(serialExecTimeMs)
	serialCodeShare = 1 - parallelCodeShare

	fmt.Printf("Serial determinant calculation of matrix of order %d took %d ms.\n",
		strongScalingMatrixOrder, serialExecTimeMs)
	fmt.Printf("Execution of code that can be run in parallel took %d ms.\n", potentialParallelCodeExecTimeMs)
	fmt.Printf("Parallel code share is: %f\nSerial code share is: %f\n\n", parallelCodeShare, serialCodeShare)

	for tasksNum := 2; tasksNum <= strongScalingMatrixOrder; tasksNum++ {
		start := time.Now()
		determinantCalc.DetParallel(matrix, tasksNum)
		elapsed := time.Since(start)
		parallelExecTimeMs := elapsed.Milliseconds()
		achievedSpeedup := calculateSpeedup(serialExecTimeMs, parallelExecTimeMs)
		maxSpeedup := maxSpeedupAmdahl(tasksNum)

		scalingResult := ScalingResult{
			AchievedSpeedup:       achievedSpeedup,
			MatrixOrder:           strongScalingMatrixOrder,
			MaxTheoreticalSpeedup: maxSpeedup,
			ParallelExecTimeMs:    parallelExecTimeMs,
			ParallelTasksNum:      tasksNum,
			SerialExecTimeMs:      serialExecTimeMs,
		}
		results = append(results, scalingResult)

		fmt.Printf("Parallel determinant calculation of matrix of order %d with %d tasks took %d ms.\n",
			strongScalingMatrixOrder, tasksNum, parallelExecTimeMs)
		fmt.Printf("Achieved speedup is: %fX.\nMaximum speedup according to Amdahl’s law is: %fX.\n\n",
			achievedSpeedup, maxSpeedup)
	}

	writeScalingResults(fmt.Sprintf("%s/strong_scaling_results_go.csv", resultsBasePath), results)
	fmt.Println("Successfully finished strong scaling experiment.")
}

func calculateSpeedup(serialExecTimeMs int64, parallelExecTimeMs int64) float64 {
	if serialExecTimeMs == 0 && parallelExecTimeMs == 0 {
		return 1.0
	}
	return float64(serialExecTimeMs) / float64(parallelExecTimeMs)
}

// Loads matrices of different orders from predefined files and calculates
// their determinants using parallel implementation where the number of tasks is
// equal to the order of matrix. Statistics about calculations are printed onto the console
// and written into a predefined results file: weak_scaling_results_go.csv.
//
// Return:
// 	nil
func weakScalingExperiment() {
	fmt.Println("\n\n========================")
	fmt.Println("Weak scaling experiment:")
	fmt.Println("========================\n\n")
	var results []ScalingResult

	for _, n := range availableMatrixOrders {
		matrix := loadTestMatrix(n)

		start := time.Now()
		determinantCalc.DetSerial(matrix)
		elapsed := time.Since(start)
		serialExecTimeMs := elapsed.Milliseconds()
		fmt.Printf("Serial determinant calculation of matrix of order %d took %d ms.\n", n, serialExecTimeMs)

		start = time.Now()
		determinantCalc.DetParallel(matrix, n)
		elapsed = time.Since(start)
		parallelExecTimeMs := elapsed.Milliseconds()

		achievedSpeedup := calculateSpeedup(serialExecTimeMs, parallelExecTimeMs)
		maxSpeedup := maxSpeedupGustafson(n)

		scalingResult := ScalingResult{
			AchievedSpeedup:       achievedSpeedup,
			MatrixOrder:           n,
			MaxTheoreticalSpeedup: maxSpeedup,
			ParallelExecTimeMs:    parallelExecTimeMs,
			ParallelTasksNum:      n,
			SerialExecTimeMs:      serialExecTimeMs,
		}
		results = append(results, scalingResult)

		fmt.Printf("Parallel determinant calculation of matrix of order %d with %d tasks took %d ms.\n",
			n, n, parallelExecTimeMs)
		fmt.Printf("Achieved speedup is: %fX.\nMaximum speedup according to Gustafson’s law is: %fX.\n\n",
			achievedSpeedup, maxSpeedup)
	}

	writeScalingResults(fmt.Sprintf("%s/weak_scaling_results_go.csv", resultsBasePath), results)
	fmt.Println("Successfully finished weak scaling experiment.")
}

func main() {
	strongScalingExperiment()
	weakScalingExperiment()
}
